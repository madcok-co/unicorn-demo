package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListSalesOrders returns paginated list of sales orders
func ListSalesOrders(ctx *context.Context, req domain.ListSalesOrdersDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var orders []domain.SalesOrder
	var total int64

	query := db.Model(&domain.SalesOrder{}).
		Preload("Items").
		Where("tenant_id = ?", tenantID)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("code LIKE ? OR customer_name LIKE ? OR reference_no LIKE ? OR customer_po LIKE ?", search, search, search, search)
	}
	if req.CustomerID != "" {
		query = query.Where("customer_id = ?", req.CustomerID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Priority != "" {
		query = query.Where("priority = ?", req.Priority)
	}
	if req.PaymentStatus != "" {
		query = query.Where("payment_status = ?", req.PaymentStatus)
	}
	if req.FulfillmentStatus != "" {
		query = query.Where("fulfillment_status = ?", req.FulfillmentStatus)
	}
	if req.SalespersonID != "" {
		query = query.Where("salesperson_id = ?", req.SalespersonID)
	}
	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.DateFrom != nil {
		query = query.Where("order_date >= ?", req.DateFrom)
	}
	if req.DateTo != nil {
		query = query.Where("order_date <= ?", req.DateTo)
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
	if err := query.Offset(offset).Limit(req.Limit).Order("order_date DESC, created_at DESC").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch sales orders: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       orders,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetSalesOrder returns a single sales order by ID
func GetSalesOrder(ctx *context.Context, id string) (*domain.SalesOrder, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var order domain.SalesOrder
	if err := db.Preload("Items").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&order).Error; err != nil {
		return nil, fmt.Errorf("sales order not found: %w", err)
	}

	return &order, nil
}

// CreateSalesOrder creates a new sales order
func CreateSalesOrder(ctx *context.Context, req domain.CreateSalesOrderDTO) (*domain.SalesOrder, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.SalesOrder
	if err := db.Where("code = ? AND tenant_id = ?", req.Code, tenantID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("sales order with code %s already exists", req.Code)
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

	// Get warehouse info if provided
	var warehouseName string
	if req.WarehouseID != "" {
		var warehouse domain.Warehouse
		if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.WarehouseID, tenantID, true).First(&warehouse).Error; err == nil {
			warehouseName = warehouse.Name
		}
	}

	// Get quotation info if provided
	var quotationCode *string
	if req.QuotationID != nil {
		var quotation domain.SalesQuotation
		if err := db.Where("id = ? AND tenant_id = ?", *req.QuotationID, tenantID).First(&quotation).Error; err == nil {
			quotationCode = &quotation.Code
		}
	}

	// Calculate totals
	var subtotal, totalDiscount, totalTax, grandTotal float64
	items := make([]domain.SalesOrderItem, len(req.Items))

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

		items[i] = domain.SalesOrderItem{
			ID:           uuid.New().String(),
			TenantID:     tenantID,
			LineNumber:   i + 1,
			ProductID:    product.ID,
			ProductCode:  product.Code,
			ProductName:  product.Name,
			ProductType:  product.Type,
			Description:  itemReq.Description,
			Quantity:     itemReq.Quantity,
			DeliveredQty: 0,
			InvoicedQty:  0,
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

	// Set default priority if not provided
	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	order := &domain.SalesOrder{
		ID:                uuid.New().String(),
		Code:              req.Code,
		TenantID:          tenantID,
		QuotationID:       req.QuotationID,
		QuotationCode:     quotationCode,
		CustomerID:        customer.ID,
		CustomerName:      customer.Name,
		CustomerType:      customer.Type,
		OrderDate:         req.OrderDate,
		ExpectedDate:      req.ExpectedDate,
		Status:            "draft",
		Priority:          priority,
		ReferenceNo:       req.ReferenceNo,
		CustomerPO:        req.CustomerPO,
		SalespersonID:     req.SalespersonID,
		SalespersonName:   salespersonName,
		Currency:          req.Currency,
		ExchangeRate:      req.ExchangeRate,
		Subtotal:          subtotal,
		TotalDiscount:     totalDiscount,
		TotalTax:          totalTax,
		ShippingCost:      req.ShippingCost,
		OtherCost:         req.OtherCost,
		GrandTotal:        grandTotal,
		PaymentTermID:     req.PaymentTermID,
		PaymentTermDays:   req.PaymentTermDays,
		PaymentStatus:     "unpaid",
		WarehouseID:       req.WarehouseID,
		WarehouseName:     warehouseName,
		FulfillmentStatus: "pending",
		DeliveredQty:      0,
		InvoicedQty:       0,
		ShippingMethod:    req.ShippingMethod,
		ShippingAddress:   req.ShippingAddress,
		ShippingCity:      req.ShippingCity,
		ShippingState:     req.ShippingState,
		ShippingZip:       req.ShippingZip,
		ShippingCountry:   req.ShippingCountry,
		Notes:             req.Notes,
		InternalNotes:     req.InternalNotes,
		Terms:             req.Terms,
		Tags:              req.Tags,
		Items:             items,
		Active:            true,
		CreatedBy:         userID,
		UpdatedBy:         userID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := db.Create(order).Error; err != nil {
		return nil, fmt.Errorf("failed to create sales order: %w", err)
	}

	return order, nil
}

// UpdateSalesOrder updates an existing sales order
func UpdateSalesOrder(ctx *context.Context, id string, req domain.UpdateSalesOrderDTO) (*domain.SalesOrder, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var order domain.SalesOrder
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", id, tenantID).First(&order).Error; err != nil {
		return nil, fmt.Errorf("sales order not found: %w", err)
	}

	// Only draft or confirmed orders can be updated (not processing, completed, or cancelled)
	if order.Status != "draft" && order.Status != "confirmed" && req.Status == "" {
		return nil, fmt.Errorf("only draft or confirmed orders can be updated")
	}

	// Update fields if provided
	if req.ExpectedDate != nil {
		order.ExpectedDate = *req.ExpectedDate
	}
	if req.Status != "" {
		order.Status = req.Status
		// Update workflow timestamps
		now := time.Now()
		switch req.Status {
		case "confirmed":
			if order.ConfirmedDate == nil {
				order.ConfirmedDate = &now
			}
		}
	}
	if req.Priority != "" {
		order.Priority = req.Priority
	}
	if req.ReferenceNo != "" {
		order.ReferenceNo = req.ReferenceNo
	}
	if req.CustomerPO != "" {
		order.CustomerPO = req.CustomerPO
	}
	if req.SalespersonID != "" {
		order.SalespersonID = req.SalespersonID
		// Update salesperson name
		var salesperson domain.User
		if err := db.Where("id = ? AND tenant_id = ?", req.SalespersonID, tenantID).First(&salesperson).Error; err == nil {
			order.SalespersonName = salesperson.Name
		}
	}
	if req.PaymentTermID != "" {
		order.PaymentTermID = req.PaymentTermID
	}
	if req.PaymentTermDays != nil {
		order.PaymentTermDays = *req.PaymentTermDays
	}
	if req.PaymentStatus != "" {
		order.PaymentStatus = req.PaymentStatus
	}
	if req.WarehouseID != "" {
		order.WarehouseID = req.WarehouseID
		// Update warehouse name
		var warehouse domain.Warehouse
		if err := db.Where("id = ? AND tenant_id = ?", req.WarehouseID, tenantID).First(&warehouse).Error; err == nil {
			order.WarehouseName = warehouse.Name
		}
	}
	if req.FulfillmentStatus != "" {
		order.FulfillmentStatus = req.FulfillmentStatus
	}
	if req.ShippingMethod != "" {
		order.ShippingMethod = req.ShippingMethod
	}
	if req.ShippingAddress != "" {
		order.ShippingAddress = req.ShippingAddress
	}
	if req.ShippingCity != "" {
		order.ShippingCity = req.ShippingCity
	}
	if req.ShippingState != "" {
		order.ShippingState = req.ShippingState
	}
	if req.ShippingZip != "" {
		order.ShippingZip = req.ShippingZip
	}
	if req.ShippingCountry != "" {
		order.ShippingCountry = req.ShippingCountry
	}
	if req.ShippingCost != nil {
		order.ShippingCost = *req.ShippingCost
		// Recalculate grand total
		order.GrandTotal = order.Subtotal - order.TotalDiscount + order.TotalTax + order.ShippingCost + order.OtherCost
	}
	if req.OtherCost != nil {
		order.OtherCost = *req.OtherCost
		// Recalculate grand total
		order.GrandTotal = order.Subtotal - order.TotalDiscount + order.TotalTax + order.ShippingCost + order.OtherCost
	}
	if req.Notes != "" {
		order.Notes = req.Notes
	}
	if req.InternalNotes != "" {
		order.InternalNotes = req.InternalNotes
	}
	if req.Terms != "" {
		order.Terms = req.Terms
	}
	if req.Tags != nil {
		order.Tags = req.Tags
	}

	order.UpdatedBy = userID
	order.UpdatedAt = time.Now()

	if err := db.Save(&order).Error; err != nil {
		return nil, fmt.Errorf("failed to update sales order: %w", err)
	}

	return &order, nil
}

// DeleteSalesOrder soft deletes a sales order
func DeleteSalesOrder(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var order domain.SalesOrder
	if err := db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&order).Error; err != nil {
		return fmt.Errorf("sales order not found: %w", err)
	}

	// Only draft orders can be deleted
	if order.Status != "draft" {
		return fmt.Errorf("only draft orders can be deleted")
	}

	order.Active = false
	order.UpdatedBy = userID
	order.UpdatedAt = time.Now()

	if err := db.Save(&order).Error; err != nil {
		return fmt.Errorf("failed to delete sales order: %w", err)
	}

	return nil
}

// GetSalesOrderStats returns sales order statistics
func GetSalesOrderStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalOrders int64
	db.Model(&domain.SalesOrder{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalOrders)

	var statusStats []struct {
		Status string
		Count  int64
		Total  float64
	}
	db.Model(&domain.SalesOrder{}).
		Select("status, COUNT(*) as count, SUM(grandtotal) as total").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("status").
		Scan(&statusStats)

	var paymentStatusStats []struct {
		PaymentStatus string
		Count         int64
		Total         float64
	}
	db.Model(&domain.SalesOrder{}).
		Select("payment_status, COUNT(*) as count, SUM(grandtotal) as total").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("payment_status").
		Scan(&paymentStatusStats)

	var fulfillmentStatusStats []struct {
		FulfillmentStatus string
		Count             int64
	}
	db.Model(&domain.SalesOrder{}).
		Select("fulfillment_status, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("fulfillment_status").
		Scan(&fulfillmentStatusStats)

	var totalValue float64
	db.Model(&domain.SalesOrder{}).
		Select("SUM(grandtotal)").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Scan(&totalValue)

	var confirmedValue float64
	db.Model(&domain.SalesOrder{}).
		Select("SUM(grandtotal)").
		Where("tenant_id = ? AND active = ? AND status IN ?", tenantID, true, []string{"confirmed", "processing", "completed"}).
		Scan(&confirmedValue)

	return map[string]interface{}{
		"total_orders":                 totalOrders,
		"total_value":                  totalValue,
		"confirmed_value":              confirmedValue,
		"orders_by_status":             statusStats,
		"orders_by_payment_status":     paymentStatusStats,
		"orders_by_fulfillment_status": fulfillmentStatusStats,
	}, nil
}
