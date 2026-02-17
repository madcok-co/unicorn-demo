package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListPurchaseOrders returns paginated list of purchase orders
func ListPurchaseOrders(ctx *context.Context, req domain.ListPurchaseOrdersDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var orders []domain.PurchaseOrder
	var total int64

	query := db.Model(&domain.PurchaseOrder{}).
		Preload("Items").
		Where("tenant_id = ?", tenantID)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("code LIKE ? OR vendor_name LIKE ? OR reference_no LIKE ? OR vendor_quote_no LIKE ?", search, search, search, search)
	}
	if req.VendorID != "" {
		query = query.Where("vendor_id = ?", req.VendorID)
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
	if req.ReceiptStatus != "" {
		query = query.Where("receipt_status = ?", req.ReceiptStatus)
	}
	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.BuyerID != "" {
		query = query.Where("buyer_id = ?", req.BuyerID)
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
		return nil, fmt.Errorf("failed to fetch purchase orders: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       orders,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetPurchaseOrder returns a single purchase order by ID
func GetPurchaseOrder(ctx *context.Context, id string) (*domain.PurchaseOrder, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var order domain.PurchaseOrder
	if err := db.Preload("Items").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&order).Error; err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	return &order, nil
}

// CreatePurchaseOrder creates a new purchase order
func CreatePurchaseOrder(ctx *context.Context, req domain.CreatePurchaseOrderDTO) (*domain.PurchaseOrder, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.PurchaseOrder
	if err := db.Where("code = ? AND tenant_id = ?", req.Code, tenantID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("purchase order with code %s already exists", req.Code)
	}

	// Get vendor info
	var vendor domain.Vendor
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.VendorID, tenantID, true).First(&vendor).Error; err != nil {
		return nil, fmt.Errorf("vendor not found: %w", err)
	}

	// Get buyer info if provided
	var buyerName string
	if req.BuyerID != "" {
		var buyer domain.User
		if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.BuyerID, tenantID, true).First(&buyer).Error; err == nil {
			buyerName = buyer.Name
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

	// Get purchase request info if provided
	var requestCode *string
	if req.RequestID != nil {
		var request domain.PurchaseRequest
		if err := db.Where("id = ? AND tenant_id = ?", *req.RequestID, tenantID).First(&request).Error; err == nil {
			requestCode = &request.Code
		}
	}

	// Calculate totals
	var subtotal, totalDiscount, totalTax, grandTotal float64
	items := make([]domain.PurchaseOrderItem, len(req.Items))

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

		items[i] = domain.PurchaseOrderItem{
			ID:           uuid.New().String(),
			TenantID:     tenantID,
			LineNumber:   i + 1,
			ProductID:    product.ID,
			ProductCode:  product.Code,
			ProductName:  product.Name,
			ProductType:  product.Type,
			Description:  itemReq.Description,
			Quantity:     itemReq.Quantity,
			ReceivedQty:  0,
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

	order := &domain.PurchaseOrder{
		ID:              uuid.New().String(),
		Code:            req.Code,
		TenantID:        tenantID,
		RequestID:       req.RequestID,
		RequestCode:     requestCode,
		VendorID:        vendor.ID,
		VendorName:      vendor.Name,
		VendorType:      vendor.Type,
		OrderDate:       req.OrderDate,
		ExpectedDate:    req.ExpectedDate,
		Status:          "draft",
		Priority:        priority,
		ReferenceNo:     req.ReferenceNo,
		VendorQuoteNo:   req.VendorQuoteNo,
		BuyerID:         req.BuyerID,
		BuyerName:       buyerName,
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
		PaymentStatus:   "unpaid",
		WarehouseID:     req.WarehouseID,
		WarehouseName:   warehouseName,
		ReceiptStatus:   "pending",
		ReceivedQty:     0,
		InvoicedQty:     0,
		ShippingMethod:  req.ShippingMethod,
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

	if err := db.Create(order).Error; err != nil {
		return nil, fmt.Errorf("failed to create purchase order: %w", err)
	}

	// If linked to a purchase request, update the request items ordered_qty
	if req.RequestID != nil {
		var request domain.PurchaseRequest
		if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", *req.RequestID, tenantID).First(&request).Error; err == nil {
			tx := db.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			// Update ordered_qty for each request item (simplified - assumes same products)
			for _, orderItem := range items {
				for _, reqItem := range request.Items {
					if reqItem.ProductID == orderItem.ProductID {
						newOrderedQty := reqItem.OrderedQty + orderItem.Quantity
						if err := tx.Model(&domain.PurchaseRequestItem{}).
							Where("id = ? AND tenant_id = ?", reqItem.ID, tenantID).
							Update("ordered_qty", newOrderedQty).Error; err != nil {
							tx.Rollback()
							break
						}
					}
				}
			}

			tx.Commit()
		}
	}

	return order, nil
}

// UpdatePurchaseOrder updates an existing purchase order
func UpdatePurchaseOrder(ctx *context.Context, id string, req domain.UpdatePurchaseOrderDTO) (*domain.PurchaseOrder, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var order domain.PurchaseOrder
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", id, tenantID).First(&order).Error; err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	// Only draft or sent orders can be updated (not confirmed, receiving, completed, or cancelled)
	if order.Status != "draft" && order.Status != "sent" && req.Status == "" {
		return nil, fmt.Errorf("only draft or sent purchase orders can be updated")
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
	if req.VendorQuoteNo != "" {
		order.VendorQuoteNo = req.VendorQuoteNo
	}
	if req.BuyerID != "" {
		order.BuyerID = req.BuyerID
		// Update buyer name
		var buyer domain.User
		if err := db.Where("id = ? AND tenant_id = ?", req.BuyerID, tenantID).First(&buyer).Error; err == nil {
			order.BuyerName = buyer.Name
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
	if req.ReceiptStatus != "" {
		order.ReceiptStatus = req.ReceiptStatus
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
		return nil, fmt.Errorf("failed to update purchase order: %w", err)
	}

	return &order, nil
}

// DeletePurchaseOrder soft deletes a purchase order
func DeletePurchaseOrder(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var order domain.PurchaseOrder
	if err := db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&order).Error; err != nil {
		return fmt.Errorf("purchase order not found: %w", err)
	}

	// Only draft orders can be deleted
	if order.Status != "draft" {
		return fmt.Errorf("only draft purchase orders can be deleted")
	}

	order.Active = false
	order.UpdatedBy = userID
	order.UpdatedAt = time.Now()

	if err := db.Save(&order).Error; err != nil {
		return fmt.Errorf("failed to delete purchase order: %w", err)
	}

	return nil
}

// GetPurchaseOrderStats returns purchase order statistics
func GetPurchaseOrderStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalOrders int64
	db.Model(&domain.PurchaseOrder{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalOrders)

	var statusStats []struct {
		Status string
		Count  int64
		Total  float64
	}
	db.Model(&domain.PurchaseOrder{}).
		Select("status, COUNT(*) as count, SUM(grandtotal) as total").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("status").
		Scan(&statusStats)

	var paymentStatusStats []struct {
		PaymentStatus string
		Count         int64
		Total         float64
	}
	db.Model(&domain.PurchaseOrder{}).
		Select("payment_status, COUNT(*) as count, SUM(grandtotal) as total").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("payment_status").
		Scan(&paymentStatusStats)

	var receiptStatusStats []struct {
		ReceiptStatus string
		Count         int64
	}
	db.Model(&domain.PurchaseOrder{}).
		Select("receipt_status, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("receipt_status").
		Scan(&receiptStatusStats)

	var totalValue float64
	db.Model(&domain.PurchaseOrder{}).
		Select("SUM(grandtotal)").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Scan(&totalValue)

	var confirmedValue float64
	db.Model(&domain.PurchaseOrder{}).
		Select("SUM(grandtotal)").
		Where("tenant_id = ? AND active = ? AND status IN ?", tenantID, true, []string{"confirmed", "receiving", "completed"}).
		Scan(&confirmedValue)

	return map[string]interface{}{
		"total_orders":             totalOrders,
		"total_value":              totalValue,
		"confirmed_value":          confirmedValue,
		"orders_by_status":         statusStats,
		"orders_by_payment_status": paymentStatusStats,
		"orders_by_receipt_status": receiptStatusStats,
	}, nil
}
