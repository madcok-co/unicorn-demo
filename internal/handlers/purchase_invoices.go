package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListPurchaseInvoices returns paginated list of purchase invoices
func ListPurchaseInvoices(ctx *context.Context, req domain.ListPurchaseInvoicesDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var invoices []domain.PurchaseInvoice
	var total int64

	query := db.Model(&domain.PurchaseInvoice{}).
		Preload("Items").
		Where("tenant_id = ?", tenantID)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("code LIKE ? OR vendor_name LIKE ? OR reference_no LIKE ? OR vendor_invoice_no LIKE ? OR tax_invoice_no LIKE ?", search, search, search, search, search)
	}
	if req.VendorID != "" {
		query = query.Where("vendor_id = ?", req.VendorID)
	}
	if req.OrderID != "" {
		query = query.Where("order_id = ?", req.OrderID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.PaymentStatus != "" {
		query = query.Where("payment_status = ?", req.PaymentStatus)
	}
	if req.DateFrom != nil {
		query = query.Where("invoice_date >= ?", req.DateFrom)
	}
	if req.DateTo != nil {
		query = query.Where("invoice_date <= ?", req.DateTo)
	}
	if req.DueDateFrom != nil {
		query = query.Where("due_date >= ?", req.DueDateFrom)
	}
	if req.DueDateTo != nil {
		query = query.Where("due_date <= ?", req.DueDateTo)
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
	if err := query.Offset(offset).Limit(req.Limit).Order("invoice_date DESC, created_at DESC").Find(&invoices).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch purchase invoices: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       invoices,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetPurchaseInvoice returns a single purchase invoice by ID
func GetPurchaseInvoice(ctx *context.Context, id string) (*domain.PurchaseInvoice, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var invoice domain.PurchaseInvoice
	if err := db.Preload("Items").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&invoice).Error; err != nil {
		return nil, fmt.Errorf("purchase invoice not found: %w", err)
	}

	return &invoice, nil
}

// CreatePurchaseInvoice creates a new purchase invoice
func CreatePurchaseInvoice(ctx *context.Context, req domain.CreatePurchaseInvoiceDTO) (*domain.PurchaseInvoice, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.PurchaseInvoice
	if err := db.Where("code = ? AND tenant_id = ?", req.Code, tenantID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("purchase invoice with code %s already exists", req.Code)
	}

	// Get vendor info
	var vendor domain.Vendor
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.VendorID, tenantID, true).First(&vendor).Error; err != nil {
		return nil, fmt.Errorf("vendor not found: %w", err)
	}

	// Get order info if provided
	var orderCode *string
	var order *domain.PurchaseOrder
	if req.OrderID != nil {
		var o domain.PurchaseOrder
		if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", *req.OrderID, tenantID).First(&o).Error; err == nil {
			orderCode = &o.Code
			order = &o
		}
	}

	// Get receipt info if provided
	var receiptCode *string
	if req.ReceiptID != nil {
		var receipt domain.GoodsReceipt
		if err := db.Where("id = ? AND tenant_id = ?", *req.ReceiptID, tenantID).First(&receipt).Error; err == nil {
			receiptCode = &receipt.Code
		}
	}

	// Calculate totals
	var subtotal, totalDiscount, totalTax, grandTotal float64
	items := make([]domain.PurchaseInvoiceItem, len(req.Items))
	orderItemUpdates := make(map[string]float64) // Track invoiced qty updates per order item

	for i, itemReq := range req.Items {
		// Get product info
		var product domain.Product
		if err := db.Where("id = ? AND tenant_id = ? AND active = ?", itemReq.ProductID, tenantID, true).First(&product).Error; err != nil {
			return nil, fmt.Errorf("product not found for item %d: %w", i+1, err)
		}

		// If order item ID is provided, validate and track invoiced quantity
		var orderItem *domain.PurchaseOrderItem
		if itemReq.OrderItemID != nil && order != nil {
			for idx := range order.Items {
				if order.Items[idx].ID == *itemReq.OrderItemID {
					orderItem = &order.Items[idx]
					break
				}
			}

			if orderItem == nil {
				return nil, fmt.Errorf("order item %s not found in purchase order", *itemReq.OrderItemID)
			}

			// Validate product ID matches
			if orderItem.ProductID != itemReq.ProductID {
				return nil, fmt.Errorf("product ID mismatch for order item %s", *itemReq.OrderItemID)
			}

			// Check if invoiced quantity exceeds remaining quantity
			remainingQty := orderItem.Quantity - orderItem.InvoicedQty
			if itemReq.Quantity > remainingQty {
				return nil, fmt.Errorf("invoiced quantity %.2f exceeds remaining quantity %.2f for product %s",
					itemReq.Quantity, remainingQty, orderItem.ProductName)
			}

			// Track the invoiced quantity for this order item
			orderItemUpdates[*itemReq.OrderItemID] = orderItem.InvoicedQty + itemReq.Quantity
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

		items[i] = domain.PurchaseInvoiceItem{
			ID:           uuid.New().String(),
			TenantID:     tenantID,
			OrderItemID:  itemReq.OrderItemID,
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
			AccountID:    itemReq.AccountID,
			Notes:        itemReq.Notes,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		subtotal += itemSubtotal
		totalDiscount += itemDiscount
		totalTax += itemTax
	}

	grandTotal = subtotal - totalDiscount + totalTax + req.ShippingCost + req.OtherCost
	amountDue := grandTotal // Initially, amount due equals grand total

	invoice := &domain.PurchaseInvoice{
		ID:               uuid.New().String(),
		Code:             req.Code,
		TenantID:         tenantID,
		OrderID:          req.OrderID,
		OrderCode:        orderCode,
		ReceiptID:        req.ReceiptID,
		ReceiptCode:      receiptCode,
		VendorID:         vendor.ID,
		VendorName:       vendor.Name,
		VendorType:       vendor.Type,
		VendorTaxID:      vendor.TaxID,
		VendorInvoiceNo:  req.VendorInvoiceNo,
		InvoiceDate:      req.InvoiceDate,
		DueDate:          req.DueDate,
		Status:           "draft",
		ReferenceNo:      req.ReferenceNo,
		TaxInvoiceNo:     req.TaxInvoiceNo,
		Currency:         req.Currency,
		ExchangeRate:     req.ExchangeRate,
		Subtotal:         subtotal,
		TotalDiscount:    totalDiscount,
		TotalTax:         totalTax,
		ShippingCost:     req.ShippingCost,
		OtherCost:        req.OtherCost,
		GrandTotal:       grandTotal,
		AmountPaid:       0,
		AmountDue:        amountDue,
		PaymentTermID:    req.PaymentTermID,
		PaymentTermDays:  req.PaymentTermDays,
		PaymentStatus:    "unpaid",
		ExpenseAccountID: req.ExpenseAccountID,
		APAccountID:      req.APAccountID,
		Notes:            req.Notes,
		InternalNotes:    req.InternalNotes,
		Terms:            req.Terms,
		Tags:             req.Tags,
		Items:            items,
		Active:           true,
		CreatedBy:        userID,
		UpdatedBy:        userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// CRITICAL: If we have order item updates, create invoice in a transaction
	if len(orderItemUpdates) > 0 && order != nil {
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if err := tx.Create(invoice).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create purchase invoice: %w", err)
		}

		// CRITICAL: Update purchase order items with invoiced quantities
		for orderItemID, newInvoicedQty := range orderItemUpdates {
			if err := tx.Model(&domain.PurchaseOrderItem{}).
				Where("id = ? AND tenant_id = ?", orderItemID, tenantID).
				Update("invoiced_qty", newInvoicedQty).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to update order item invoiced quantity: %w", err)
			}
		}

		// Update purchase order invoiced quantity
		var totalOrderedQty, totalInvoicedQty float64
		for _, item := range order.Items {
			totalOrderedQty += item.Quantity
			if updatedQty, exists := orderItemUpdates[item.ID]; exists {
				totalInvoicedQty += updatedQty
			} else {
				totalInvoicedQty += item.InvoicedQty
			}
		}

		if err := tx.Model(&domain.PurchaseOrder{}).
			Where("id = ? AND tenant_id = ?", order.ID, tenantID).
			Updates(map[string]interface{}{
				"invoiced_qty": totalInvoicedQty,
				"updated_at":   time.Now(),
			}).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update purchase order invoiced quantity: %w", err)
		}

		if err := tx.Commit().Error; err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}
	} else {
		// No order item updates, just create the invoice
		if err := db.Create(invoice).Error; err != nil {
			return nil, fmt.Errorf("failed to create purchase invoice: %w", err)
		}
	}

	return invoice, nil
}

// UpdatePurchaseInvoice updates an existing purchase invoice
func UpdatePurchaseInvoice(ctx *context.Context, id string, req domain.UpdatePurchaseInvoiceDTO) (*domain.PurchaseInvoice, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var invoice domain.PurchaseInvoice
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", id, tenantID).First(&invoice).Error; err != nil {
		return nil, fmt.Errorf("purchase invoice not found: %w", err)
	}

	// Only draft invoices can be updated (except for status and payment updates)
	if invoice.Status != "draft" && req.Status == "" && req.AmountPaid == nil {
		return nil, fmt.Errorf("only draft invoices can be updated")
	}

	// Update fields if provided
	if req.DueDate != nil {
		invoice.DueDate = *req.DueDate
	}
	if req.Status != "" {
		invoice.Status = req.Status
		// Update workflow timestamps
		now := time.Now()
		switch req.Status {
		case "received":
			if invoice.ReceivedAt == nil {
				invoice.ReceivedAt = &now
			}
		case "verified":
			if invoice.VerifiedAt == nil {
				invoice.VerifiedAt = &now
			}
		case "approved":
			if invoice.ApprovedAt == nil {
				invoice.ApprovedAt = &now
			}
		case "paid":
			if invoice.PaidAt == nil {
				invoice.PaidAt = &now
			}
			invoice.PaymentStatus = "paid"
			invoice.AmountPaid = invoice.GrandTotal
			invoice.AmountDue = 0
		}
	}
	if req.ReferenceNo != "" {
		invoice.ReferenceNo = req.ReferenceNo
	}
	if req.VendorInvoiceNo != "" {
		invoice.VendorInvoiceNo = req.VendorInvoiceNo
	}
	if req.TaxInvoiceNo != "" {
		invoice.TaxInvoiceNo = req.TaxInvoiceNo
	}
	if req.PaymentTermID != "" {
		invoice.PaymentTermID = req.PaymentTermID
	}
	if req.PaymentTermDays != nil {
		invoice.PaymentTermDays = *req.PaymentTermDays
	}
	if req.PaymentStatus != "" {
		invoice.PaymentStatus = req.PaymentStatus
	}
	if req.AmountPaid != nil {
		invoice.AmountPaid = *req.AmountPaid
		invoice.AmountDue = invoice.GrandTotal - invoice.AmountPaid

		// CRITICAL: Auto-update payment status based on amount paid (same logic as sales)
		if invoice.AmountPaid == 0 {
			invoice.PaymentStatus = "unpaid"
		} else if invoice.AmountPaid >= invoice.GrandTotal {
			invoice.PaymentStatus = "paid"
			now := time.Now()
			if invoice.PaidAt == nil {
				invoice.PaidAt = &now
			}
		} else {
			invoice.PaymentStatus = "partial"
		}

		// Check if overdue
		if invoice.AmountDue > 0 && time.Now().After(invoice.DueDate) {
			invoice.PaymentStatus = "overdue"
		}
	}
	if req.ExpenseAccountID != "" {
		invoice.ExpenseAccountID = req.ExpenseAccountID
	}
	if req.APAccountID != "" {
		invoice.APAccountID = req.APAccountID
	}
	if req.ShippingCost != nil {
		invoice.ShippingCost = *req.ShippingCost
		// Recalculate grand total and amount due
		invoice.GrandTotal = invoice.Subtotal - invoice.TotalDiscount + invoice.TotalTax + invoice.ShippingCost + invoice.OtherCost
		invoice.AmountDue = invoice.GrandTotal - invoice.AmountPaid
	}
	if req.OtherCost != nil {
		invoice.OtherCost = *req.OtherCost
		// Recalculate grand total and amount due
		invoice.GrandTotal = invoice.Subtotal - invoice.TotalDiscount + invoice.TotalTax + invoice.ShippingCost + invoice.OtherCost
		invoice.AmountDue = invoice.GrandTotal - invoice.AmountPaid
	}
	if req.Notes != "" {
		invoice.Notes = req.Notes
	}
	if req.InternalNotes != "" {
		invoice.InternalNotes = req.InternalNotes
	}
	if req.Terms != "" {
		invoice.Terms = req.Terms
	}
	if req.Tags != nil {
		invoice.Tags = req.Tags
	}

	invoice.UpdatedBy = userID
	invoice.UpdatedAt = time.Now()

	if err := db.Save(&invoice).Error; err != nil {
		return nil, fmt.Errorf("failed to update purchase invoice: %w", err)
	}

	return &invoice, nil
}

// DeletePurchaseInvoice soft deletes a purchase invoice
func DeletePurchaseInvoice(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var invoice domain.PurchaseInvoice
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", id, tenantID).First(&invoice).Error; err != nil {
		return fmt.Errorf("purchase invoice not found: %w", err)
	}

	// Only draft invoices can be deleted
	if invoice.Status != "draft" {
		return fmt.Errorf("only draft purchase invoices can be deleted")
	}

	// CRITICAL: If invoice is linked to an order, revert the invoiced quantities
	if invoice.OrderID != nil {
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Get the order
		var order domain.PurchaseOrder
		if err := tx.Preload("Items").Where("id = ? AND tenant_id = ?", *invoice.OrderID, tenantID).First(&order).Error; err == nil {
			// Revert invoiced quantities for each invoice item
			for _, invoiceItem := range invoice.Items {
				if invoiceItem.OrderItemID != nil {
					if err := tx.Model(&domain.PurchaseOrderItem{}).
						Where("id = ? AND tenant_id = ?", *invoiceItem.OrderItemID, tenantID).
						Update("invoiced_qty", db.Raw("invoiced_qty - ?", invoiceItem.Quantity)).Error; err != nil {
						tx.Rollback()
						return fmt.Errorf("failed to revert order item invoiced quantity: %w", err)
					}
				}
			}

			// Recalculate purchase order invoiced quantity
			var totalInvoicedQty float64
			for _, item := range order.Items {
				// Get updated invoiced qty after revert
				var updatedItem domain.PurchaseOrderItem
				if err := tx.Where("id = ?", item.ID).First(&updatedItem).Error; err == nil {
					totalInvoicedQty += updatedItem.InvoicedQty
				}
			}

			if err := tx.Model(&domain.PurchaseOrder{}).
				Where("id = ? AND tenant_id = ?", order.ID, tenantID).
				Updates(map[string]interface{}{
					"invoiced_qty": totalInvoicedQty,
					"updated_at":   time.Now(),
				}).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update purchase order invoiced quantity: %w", err)
			}
		}

		// Soft delete the invoice
		invoice.Active = false
		invoice.UpdatedBy = userID
		invoice.UpdatedAt = time.Now()

		if err := tx.Save(&invoice).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete purchase invoice: %w", err)
		}

		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	} else {
		// No order linked, just soft delete
		invoice.Active = false
		invoice.UpdatedBy = userID
		invoice.UpdatedAt = time.Now()

		if err := db.Save(&invoice).Error; err != nil {
			return fmt.Errorf("failed to delete purchase invoice: %w", err)
		}
	}

	return nil
}

// GetPurchaseInvoiceStats returns purchase invoice statistics
func GetPurchaseInvoiceStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalInvoices int64
	db.Model(&domain.PurchaseInvoice{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalInvoices)

	var statusStats []struct {
		Status string
		Count  int64
		Total  float64
	}
	db.Model(&domain.PurchaseInvoice{}).
		Select("status, COUNT(*) as count, SUM(grandtotal) as total").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("status").
		Scan(&statusStats)

	var paymentStatusStats []struct {
		PaymentStatus string
		Count         int64
		TotalDue      float64
	}
	db.Model(&domain.PurchaseInvoice{}).
		Select("payment_status, COUNT(*) as count, SUM(amount_due) as total_due").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("payment_status").
		Scan(&paymentStatusStats)

	var totalValue float64
	db.Model(&domain.PurchaseInvoice{}).
		Select("SUM(grandtotal)").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Scan(&totalValue)

	var totalPaid float64
	db.Model(&domain.PurchaseInvoice{}).
		Select("SUM(amount_paid)").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Scan(&totalPaid)

	var totalDue float64
	db.Model(&domain.PurchaseInvoice{}).
		Select("SUM(amount_due)").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Scan(&totalDue)

	var overdueCount int64
	var overdueAmount float64
	db.Model(&domain.PurchaseInvoice{}).
		Where("tenant_id = ? AND active = ? AND payment_status = ?", tenantID, true, "overdue").
		Count(&overdueCount)
	db.Model(&domain.PurchaseInvoice{}).
		Select("SUM(amount_due)").
		Where("tenant_id = ? AND active = ? AND payment_status = ?", tenantID, true, "overdue").
		Scan(&overdueAmount)

	return map[string]interface{}{
		"total_invoices":             totalInvoices,
		"total_value":                totalValue,
		"total_paid":                 totalPaid,
		"total_due":                  totalDue,
		"overdue_count":              overdueCount,
		"overdue_amount":             overdueAmount,
		"invoices_by_status":         statusStats,
		"invoices_by_payment_status": paymentStatusStats,
	}, nil
}
