package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListGoodsReceipts returns paginated list of goods receipts
func ListGoodsReceipts(ctx *context.Context, req domain.ListGoodsReceiptsDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var receipts []domain.GoodsReceipt
	var total int64

	query := db.Model(&domain.GoodsReceipt{}).
		Preload("Items").
		Where("tenant_id = ?", tenantID)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("code LIKE ? OR order_code LIKE ? OR vendor_name LIKE ? OR reference_no LIKE ? OR delivery_note LIKE ?", search, search, search, search, search)
	}
	if req.OrderID != "" {
		query = query.Where("order_id = ?", req.OrderID)
	}
	if req.VendorID != "" {
		query = query.Where("vendor_id = ?", req.VendorID)
	}
	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.QualityStatus != "" {
		query = query.Where("quality_status = ?", req.QualityStatus)
	}
	if req.DateFrom != nil {
		query = query.Where("receipt_date >= ?", req.DateFrom)
	}
	if req.DateTo != nil {
		query = query.Where("receipt_date <= ?", req.DateTo)
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
	if err := query.Offset(offset).Limit(req.Limit).Order("receipt_date DESC, created_at DESC").Find(&receipts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch goods receipts: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       receipts,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetGoodsReceipt returns a single goods receipt by ID
func GetGoodsReceipt(ctx *context.Context, id string) (*domain.GoodsReceipt, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var receipt domain.GoodsReceipt
	if err := db.Preload("Items").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&receipt).Error; err != nil {
		return nil, fmt.Errorf("goods receipt not found: %w", err)
	}

	return &receipt, nil
}

// CreateGoodsReceipt creates a new goods receipt
func CreateGoodsReceipt(ctx *context.Context, req domain.CreateGoodsReceiptDTO) (*domain.GoodsReceipt, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.GoodsReceipt
	if err := db.Where("code = ? AND tenant_id = ?", req.Code, tenantID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("goods receipt with code %s already exists", req.Code)
	}

	// Get purchase order info
	var order domain.PurchaseOrder
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ? AND active = ?", req.OrderID, tenantID, true).First(&order).Error; err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	// Get warehouse info
	var warehouse domain.Warehouse
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.WarehouseID, tenantID, true).First(&warehouse).Error; err != nil {
		return nil, fmt.Errorf("warehouse not found: %w", err)
	}

	// Get receiver info if provided
	var receiverName string
	if req.ReceiverID != "" {
		var receiver domain.User
		if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.ReceiverID, tenantID, true).First(&receiver).Error; err == nil {
			receiverName = receiver.Name
		}
	}

	// Get inspector info if provided
	var inspectorName string
	if req.InspectorID != "" {
		var inspector domain.User
		if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.InspectorID, tenantID, true).First(&inspector).Error; err == nil {
			inspectorName = inspector.Name
		}
	}

	// Process receipt items and validate quantities
	items := make([]domain.GoodsReceiptItem, len(req.Items))
	orderItemUpdates := make(map[string]float64) // Track received qty updates per order item
	var totalRejectedQty float64

	for i, itemReq := range req.Items {
		// Find the order item
		var orderItem domain.PurchaseOrderItem
		found := false
		for _, oi := range order.Items {
			if oi.ID == itemReq.OrderItemID {
				orderItem = oi
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("order item %s not found in purchase order", itemReq.OrderItemID)
		}

		// Check if received quantity exceeds remaining quantity
		remainingQty := orderItem.Quantity - orderItem.ReceivedQty
		if itemReq.ReceivedQty > remainingQty {
			return nil, fmt.Errorf("received quantity %.2f exceeds remaining quantity %.2f for product %s",
				itemReq.ReceivedQty, remainingQty, orderItem.ProductName)
		}

		// Validate accepted + rejected = received
		if itemReq.AcceptedQty+itemReq.RejectedQty != itemReq.ReceivedQty {
			return nil, fmt.Errorf("accepted qty (%.2f) + rejected qty (%.2f) must equal received qty (%.2f) for product %s",
				itemReq.AcceptedQty, itemReq.RejectedQty, itemReq.ReceivedQty, orderItem.ProductName)
		}

		items[i] = domain.GoodsReceiptItem{
			ID:            uuid.New().String(),
			TenantID:      tenantID,
			OrderItemID:   itemReq.OrderItemID,
			LineNumber:    i + 1,
			ProductID:     orderItem.ProductID,
			ProductCode:   orderItem.ProductCode,
			ProductName:   orderItem.ProductName,
			Description:   orderItem.Description,
			OrderedQty:    orderItem.Quantity,
			ReceivedQty:   itemReq.ReceivedQty,
			AcceptedQty:   itemReq.AcceptedQty,
			RejectedQty:   itemReq.RejectedQty,
			UOM:           orderItem.UOM,
			SerialNumbers: itemReq.SerialNumbers,
			BatchNumbers:  itemReq.BatchNumbers,
			ExpiryDates:   itemReq.ExpiryDates,
			Notes:         itemReq.Notes,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Track the received quantity for this order item
		orderItemUpdates[itemReq.OrderItemID] = orderItem.ReceivedQty + itemReq.ReceivedQty
		totalRejectedQty += itemReq.RejectedQty
	}

	// Determine quality status based on rejected items
	var qualityStatus string
	if totalRejectedQty == 0 {
		qualityStatus = "pending"
	} else {
		var totalReceivedQty float64
		for _, item := range items {
			totalReceivedQty += item.ReceivedQty
		}
		if totalRejectedQty == totalReceivedQty {
			qualityStatus = "failed"
		} else {
			qualityStatus = "partial"
		}
	}

	receipt := &domain.GoodsReceipt{
		ID:            uuid.New().String(),
		Code:          req.Code,
		TenantID:      tenantID,
		OrderID:       order.ID,
		OrderCode:     order.Code,
		VendorID:      order.VendorID,
		VendorName:    order.VendorName,
		ReceiptDate:   req.ReceiptDate,
		ScheduledDate: req.ScheduledDate,
		Status:        "draft",
		ReferenceNo:   req.ReferenceNo,
		DeliveryNote:  req.DeliveryNote,
		PackingList:   req.PackingList,
		WarehouseID:   warehouse.ID,
		WarehouseName: warehouse.Name,
		ReceiverID:    req.ReceiverID,
		ReceiverName:  receiverName,
		InspectorID:   req.InspectorID,
		InspectorName: inspectorName,
		QualityStatus: qualityStatus,
		RejectedQty:   totalRejectedQty,
		Notes:         req.Notes,
		InternalNotes: req.InternalNotes,
		Tags:          req.Tags,
		Items:         items,
		Active:        true,
		CreatedBy:     userID,
		UpdatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// CRITICAL: Start transaction to create receipt and update order items
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(receipt).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create goods receipt: %w", err)
	}

	// CRITICAL: Update purchase order items with received quantities
	for orderItemID, newReceivedQty := range orderItemUpdates {
		if err := tx.Model(&domain.PurchaseOrderItem{}).
			Where("id = ? AND tenant_id = ?", orderItemID, tenantID).
			Update("received_qty", newReceivedQty).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update order item received quantity: %w", err)
		}
	}

	// CRITICAL: Update purchase order receipt status
	var totalOrderedQty, totalReceivedQty float64
	for _, item := range order.Items {
		totalOrderedQty += item.Quantity
		if updatedQty, exists := orderItemUpdates[item.ID]; exists {
			totalReceivedQty += updatedQty
		} else {
			totalReceivedQty += item.ReceivedQty
		}
	}

	var receiptStatus string
	if totalReceivedQty == 0 {
		receiptStatus = "pending"
	} else if totalReceivedQty < totalOrderedQty {
		receiptStatus = "partial"
	} else {
		receiptStatus = "completed"
	}

	if err := tx.Model(&domain.PurchaseOrder{}).
		Where("id = ? AND tenant_id = ?", order.ID, tenantID).
		Updates(map[string]interface{}{
			"received_qty":   totalReceivedQty,
			"receipt_status": receiptStatus,
			"updated_at":     time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update purchase order receipt status: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return receipt, nil
}

// UpdateGoodsReceipt updates an existing goods receipt
func UpdateGoodsReceipt(ctx *context.Context, id string, req domain.UpdateGoodsReceiptDTO) (*domain.GoodsReceipt, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var receipt domain.GoodsReceipt
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", id, tenantID).First(&receipt).Error; err != nil {
		return nil, fmt.Errorf("goods receipt not found: %w", err)
	}

	// Only draft or received receipts can be updated
	if receipt.Status != "draft" && receipt.Status != "received" && req.Status == "" {
		return nil, fmt.Errorf("only draft or received goods receipts can be updated")
	}

	// Update fields if provided
	if req.ScheduledDate != nil {
		receipt.ScheduledDate = *req.ScheduledDate
	}
	if req.ActualDate != nil {
		receipt.ActualDate = req.ActualDate
	}
	if req.Status != "" {
		receipt.Status = req.Status
	}
	if req.QualityStatus != "" {
		receipt.QualityStatus = req.QualityStatus
	}
	if req.ReferenceNo != "" {
		receipt.ReferenceNo = req.ReferenceNo
	}
	if req.DeliveryNote != "" {
		receipt.DeliveryNote = req.DeliveryNote
	}
	if req.PackingList != "" {
		receipt.PackingList = req.PackingList
	}
	if req.ReceiverID != "" {
		receipt.ReceiverID = req.ReceiverID
		// Update receiver name
		var receiver domain.User
		if err := db.Where("id = ? AND tenant_id = ?", req.ReceiverID, tenantID).First(&receiver).Error; err == nil {
			receipt.ReceiverName = receiver.Name
		}
	}
	if req.InspectorID != "" {
		receipt.InspectorID = req.InspectorID
		// Update inspector name
		var inspector domain.User
		if err := db.Where("id = ? AND tenant_id = ?", req.InspectorID, tenantID).First(&inspector).Error; err == nil {
			receipt.InspectorName = inspector.Name
		}
	}
	if req.RejectionReason != "" {
		receipt.RejectionReason = req.RejectionReason
	}
	if req.Notes != "" {
		receipt.Notes = req.Notes
	}
	if req.InternalNotes != "" {
		receipt.InternalNotes = req.InternalNotes
	}
	if req.Tags != nil {
		receipt.Tags = req.Tags
	}
	if req.DeliveryProof != nil {
		receipt.DeliveryProof = req.DeliveryProof
	}

	receipt.UpdatedBy = userID
	receipt.UpdatedAt = time.Now()

	if err := db.Save(&receipt).Error; err != nil {
		return nil, fmt.Errorf("failed to update goods receipt: %w", err)
	}

	return &receipt, nil
}

// DeleteGoodsReceipt soft deletes a goods receipt
func DeleteGoodsReceipt(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var receipt domain.GoodsReceipt
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", id, tenantID).First(&receipt).Error; err != nil {
		return fmt.Errorf("goods receipt not found: %w", err)
	}

	// Only draft receipts can be deleted
	if receipt.Status != "draft" {
		return fmt.Errorf("only draft goods receipts can be deleted")
	}

	// CRITICAL: Start transaction to delete receipt and revert order item quantities
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Revert the received quantities in purchase order items
	for _, item := range receipt.Items {
		if err := tx.Model(&domain.PurchaseOrderItem{}).
			Where("id = ? AND tenant_id = ?", item.OrderItemID, tenantID).
			Update("received_qty", db.Raw("received_qty - ?", item.ReceivedQty)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to revert order item received quantity: %w", err)
		}
	}

	// Recalculate purchase order receipt status
	var order domain.PurchaseOrder
	if err := tx.Preload("Items").Where("id = ? AND tenant_id = ?", receipt.OrderID, tenantID).First(&order).Error; err == nil {
		var totalOrderedQty, totalReceivedQty float64
		for _, item := range order.Items {
			totalOrderedQty += item.Quantity
			// Get updated received qty after revert
			var updatedItem domain.PurchaseOrderItem
			if err := tx.Where("id = ?", item.ID).First(&updatedItem).Error; err == nil {
				totalReceivedQty += updatedItem.ReceivedQty
			}
		}

		var receiptStatus string
		if totalReceivedQty == 0 {
			receiptStatus = "pending"
		} else if totalReceivedQty < totalOrderedQty {
			receiptStatus = "partial"
		} else {
			receiptStatus = "completed"
		}

		if err := tx.Model(&domain.PurchaseOrder{}).
			Where("id = ? AND tenant_id = ?", order.ID, tenantID).
			Updates(map[string]interface{}{
				"received_qty":   totalReceivedQty,
				"receipt_status": receiptStatus,
				"updated_at":     time.Now(),
			}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update purchase order receipt status: %w", err)
		}
	}

	// Soft delete the receipt
	receipt.Active = false
	receipt.UpdatedBy = userID
	receipt.UpdatedAt = time.Now()

	if err := tx.Save(&receipt).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete goods receipt: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetGoodsReceiptStats returns goods receipt statistics
func GetGoodsReceiptStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalReceipts int64
	db.Model(&domain.GoodsReceipt{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalReceipts)

	var statusStats []struct {
		Status string
		Count  int64
	}
	db.Model(&domain.GoodsReceipt{}).
		Select("status, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("status").
		Scan(&statusStats)

	var qualityStatusStats []struct {
		QualityStatus string
		Count         int64
	}
	db.Model(&domain.GoodsReceipt{}).
		Select("quality_status, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("quality_status").
		Scan(&qualityStatusStats)

	var todayReceipts int64
	today := time.Now().Format("2006-01-02")
	db.Model(&domain.GoodsReceipt{}).
		Where("tenant_id = ? AND active = ? AND DATE(receipt_date) = ?", tenantID, true, today).
		Count(&todayReceipts)

	var pendingInspection int64
	db.Model(&domain.GoodsReceipt{}).
		Where("tenant_id = ? AND active = ? AND status = ?", tenantID, true, "received").
		Count(&pendingInspection)

	var totalRejectedQty float64
	db.Model(&domain.GoodsReceipt{}).
		Select("SUM(rejected_qty)").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Scan(&totalRejectedQty)

	return map[string]interface{}{
		"total_receipts":      totalReceipts,
		"today_receipts":      todayReceipts,
		"pending_inspection":  pendingInspection,
		"total_rejected_qty":  totalRejectedQty,
		"receipts_by_status":  statusStats,
		"receipts_by_quality": qualityStatusStats,
	}, nil
}
