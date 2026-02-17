package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListSalesQuotations returns paginated list of sales quotations
func ListSalesQuotations(ctx *context.Context, req domain.ListSalesQuotationsDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var quotations []domain.SalesQuotation
	var total int64

	query := db.Model(&domain.SalesQuotation{}).
		Preload("Items").
		Where("tenant_id = ?", tenantID)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("code LIKE ? OR customer_name LIKE ? OR reference_no LIKE ?", search, search, search)
	}
	if req.CustomerID != "" {
		query = query.Where("customer_id = ?", req.CustomerID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.SalespersonID != "" {
		query = query.Where("salesperson_id = ?", req.SalespersonID)
	}
	if req.DateFrom != nil {
		query = query.Where("quotation_date >= ?", req.DateFrom)
	}
	if req.DateTo != nil {
		query = query.Where("quotation_date <= ?", req.DateTo)
	}
	if req.Active != nil {
		query = query.Where("active = ?", *req.Active)
	}

	// Count total
	query.Count(&total)

	// Apply pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("quotation_date DESC, created_at DESC").Find(&quotations).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch sales quotations: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       quotations,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetSalesQuotation returns a single sales quotation by ID
func GetSalesQuotation(ctx *context.Context, id string) (*domain.SalesQuotation, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var quotation domain.SalesQuotation
	if err := db.Preload("Items").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&quotation).Error; err != nil {
		return nil, fmt.Errorf("sales quotation not found: %w", err)
	}

	return &quotation, nil
}

// CreateSalesQuotation creates a new sales quotation
func CreateSalesQuotation(ctx *context.Context, req domain.CreateSalesQuotationDTO) (*domain.SalesQuotation, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.SalesQuotation
	if err := db.Where("code = ? AND tenant_id = ?", req.Code, tenantID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("sales quotation with code %s already exists", req.Code)
	}

	// Get customer info
	var customer domain.Customer
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.CustomerID, tenantID, true).First(&customer).Error; err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Get salesperson info if provided
	var salespersonName string
	if req.SalespersonID != "" {
		var salesperson domain.User
		if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.SalespersonID, tenantID, true).First(&salesperson).Error; err == nil {
			salespersonName = salesperson.Name
		}
	}

	// Calculate totals
	var subtotal, totalDiscount, totalTax, grandTotal float64
	items := make([]domain.SalesQuotationItem, len(req.Items))

	for i, itemReq := range req.Items {
		// Get product info
		var product domain.Product
		if err := db.Where("id = ? AND tenant_id = ? AND active = ?", itemReq.ProductID, tenantID, true).First(&product).Error; err != nil {
			return nil, fmt.Errorf("product not found for item %d: %w", i+1, err)
		}

		// Calculate item totals
		itemSubtotal := itemReq.Quantity * itemReq.UnitPrice
		var itemDiscount float64
		if itemReq.DiscountType == "percentage" {
			itemDiscount = itemSubtotal * (itemReq.Discount / 100)
		} else {
			itemDiscount = itemReq.Discount
		}

		itemAfterDiscount := itemSubtotal - itemDiscount

		// Get tax rate
		var taxRate float64
		if itemReq.TaxID != "" {
			var tax domain.Tax
			if err := db.Where("id = ? AND tenant_id = ? AND active = ?", itemReq.TaxID, tenantID, true).First(&tax).Error; err == nil {
				taxRate = tax.Rate
			}
		}

		itemTax := itemAfterDiscount * (taxRate / 100)
		itemTotal := itemAfterDiscount + itemTax

		items[i] = domain.SalesQuotationItem{
			ID:           uuid.New().String(),
			TenantID:     tenantID,
			LineNumber:   i + 1,
			ProductID:    product.ID,
			ProductCode:  product.Code,
			ProductName:  product.Name,
			ProductType:  product.Type,
			Description:  itemReq.Description,
			Quantity:     itemReq.Quantity,
			UOM:          product.UOM,
			UnitPrice:    itemReq.UnitPrice,
			DiscountType: itemReq.DiscountType,
			Discount:     itemReq.Discount,
			TaxID:        itemReq.TaxID,
			TaxRate:      taxRate,
			TaxAmount:    itemTax,
			Subtotal:     itemSubtotal,
			Total:        itemTotal,
			Notes:        itemReq.Notes,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		subtotal += itemSubtotal
		totalDiscount += itemDiscount
		totalTax += itemTax
	}

	grandTotal = subtotal - totalDiscount + totalTax + req.ShippingCost + req.OtherCost

	quotation := &domain.SalesQuotation{
		ID:              uuid.New().String(),
		Code:            req.Code,
		TenantID:        tenantID,
		CustomerID:      customer.ID,
		CustomerName:    customer.Name,
		CustomerType:    customer.Type,
		QuotationDate:   req.QuotationDate,
		ValidUntil:      req.ValidUntil,
		ExpectedDate:    req.ExpectedDate,
		Status:          "draft",
		ReferenceNo:     req.ReferenceNo,
		SalespersonID:   req.SalespersonID,
		SalespersonName: salespersonName,
		Currency:        req.Currency,
		ExchangeRate:    req.ExchangeRate,
		Subtotal:        subtotal,
		TotalDiscount:   totalDiscount,
		TotalTax:        totalTax,
		ShippingCost:    req.ShippingCost,
		OtherCost:       req.OtherCost,
		GrandTotal:      grandTotal,
		PaymentTermID:   req.PaymentTermID,
		PaymentTermDays: req.PaymentTermDays,
		ShippingAddress: req.ShippingAddress,
		ShippingCity:    req.ShippingCity,
		ShippingState:   req.ShippingState,
		ShippingZip:     req.ShippingZip,
		ShippingCountry: req.ShippingCountry,
		Notes:           req.Notes,
		InternalNotes:   req.InternalNotes,
		Terms:           req.Terms,
		Tags:            req.Tags,
		Items:           items,
		Active:          true,
		CreatedBy:       userID,
		UpdatedBy:       userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := db.Create(quotation).Error; err != nil {
		return nil, fmt.Errorf("failed to create sales quotation: %w", err)
	}

	return quotation, nil
}

// UpdateSalesQuotation updates an existing sales quotation
func UpdateSalesQuotation(ctx *context.Context, id string, req domain.UpdateSalesQuotationDTO) (*domain.SalesQuotation, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var quotation domain.SalesQuotation
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", id, tenantID).First(&quotation).Error; err != nil {
		return nil, fmt.Errorf("sales quotation not found: %w", err)
	}

	// Only draft quotations can be updated
	if quotation.Status != "draft" && req.Status == "" {
		return nil, fmt.Errorf("only draft quotations can be updated")
	}

	// Update fields if provided
	if req.ValidUntil != nil {
		quotation.ValidUntil = *req.ValidUntil
	}
	if req.ExpectedDate != nil {
		quotation.ExpectedDate = req.ExpectedDate
	}
	if req.Status != "" {
		quotation.Status = req.Status
		// Update workflow timestamps
		now := time.Now()
		switch req.Status {
		case "sent":
			quotation.SentAt = &now
		case "accepted":
			quotation.AcceptedAt = &now
		case "rejected":
			quotation.RejectedAt = &now
		case "expired":
			quotation.ExpiredAt = &now
		}
	}
	if req.ReferenceNo != "" {
		quotation.ReferenceNo = req.ReferenceNo
	}
	if req.SalespersonID != "" {
		quotation.SalespersonID = req.SalespersonID
		// Update salesperson name
		var salesperson domain.User
		if err := db.Where("id = ? AND tenant_id = ?", req.SalespersonID, tenantID).First(&salesperson).Error; err == nil {
			quotation.SalespersonName = salesperson.Name
		}
	}
	if req.PaymentTermID != "" {
		quotation.PaymentTermID = req.PaymentTermID
	}
	if req.PaymentTermDays != nil {
		quotation.PaymentTermDays = *req.PaymentTermDays
	}
	if req.ShippingAddress != "" {
		quotation.ShippingAddress = req.ShippingAddress
	}
	if req.ShippingCity != "" {
		quotation.ShippingCity = req.ShippingCity
	}
	if req.ShippingState != "" {
		quotation.ShippingState = req.ShippingState
	}
	if req.ShippingZip != "" {
		quotation.ShippingZip = req.ShippingZip
	}
	if req.ShippingCountry != "" {
		quotation.ShippingCountry = req.ShippingCountry
	}
	if req.ShippingCost != nil {
		quotation.ShippingCost = *req.ShippingCost
		// Recalculate grand total
		quotation.GrandTotal = quotation.Subtotal - quotation.TotalDiscount + quotation.TotalTax + quotation.ShippingCost + quotation.OtherCost
	}
	if req.OtherCost != nil {
		quotation.OtherCost = *req.OtherCost
		// Recalculate grand total
		quotation.GrandTotal = quotation.Subtotal - quotation.TotalDiscount + quotation.TotalTax + quotation.ShippingCost + quotation.OtherCost
	}
	if req.Notes != "" {
		quotation.Notes = req.Notes
	}
	if req.InternalNotes != "" {
		quotation.InternalNotes = req.InternalNotes
	}
	if req.Terms != "" {
		quotation.Terms = req.Terms
	}
	if req.Tags != nil {
		quotation.Tags = req.Tags
	}

	quotation.UpdatedBy = userID
	quotation.UpdatedAt = time.Now()

	if err := db.Save(&quotation).Error; err != nil {
		return nil, fmt.Errorf("failed to update sales quotation: %w", err)
	}

	return &quotation, nil
}

// DeleteSalesQuotation soft deletes a sales quotation
func DeleteSalesQuotation(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var quotation domain.SalesQuotation
	if err := db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&quotation).Error; err != nil {
		return fmt.Errorf("sales quotation not found: %w", err)
	}

	// Only draft or rejected quotations can be deleted
	if quotation.Status != "draft" && quotation.Status != "rejected" {
		return fmt.Errorf("only draft or rejected quotations can be deleted")
	}

	quotation.Active = false
	quotation.UpdatedBy = userID
	quotation.UpdatedAt = time.Now()

	if err := db.Save(&quotation).Error; err != nil {
		return fmt.Errorf("failed to delete sales quotation: %w", err)
	}

	return nil
}

// GetSalesQuotationStats returns sales quotation statistics
func GetSalesQuotationStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalQuotations int64
	db.Model(&domain.SalesQuotation{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalQuotations)

	var statusStats []struct {
		Status string
		Count  int64
		Total  float64
	}
	db.Model(&domain.SalesQuotation{}).
		Select("status, COUNT(*) as count, SUM(grandtotal) as total").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("status").
		Scan(&statusStats)

	var totalValue float64
	db.Model(&domain.SalesQuotation{}).
		Select("SUM(grandtotal)").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Scan(&totalValue)

	var acceptedValue float64
	db.Model(&domain.SalesQuotation{}).
		Select("SUM(grandtotal)").
		Where("tenant_id = ? AND active = ? AND status = ?", tenantID, true, "accepted").
		Scan(&acceptedValue)

	return map[string]interface{}{
		"total_quotations":     totalQuotations,
		"total_value":          totalValue,
		"accepted_value":       acceptedValue,
		"quotations_by_status": statusStats,
	}, nil
}
