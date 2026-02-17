package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListDeliveryOrders returns paginated list of delivery orders
func ListDeliveryOrders(ctx *context.Context, req domain.ListDeliveryOrdersDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var deliveryOrders []domain.DeliveryOrder
	var total int64

	query := db.Model(&domain.DeliveryOrder{}).
		Preload("Items").
		Where("tenant_id = ?", tenantID)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("code LIKE ? OR order_code LIKE ? OR customer_name LIKE ? OR reference_no LIKE ? OR tracking_number LIKE ?", search, search, search, search, search)
	}
	if req.OrderID != "" {
		query = query.Where("order_id = ?", req.OrderID)
	}
	if req.CustomerID != "" {
		query = query.Where("customer_id = ?", req.CustomerID)
	}
	if req.WarehouseID != "" {
		query = query.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.DateFrom != nil {
		query = query.Where("delivery_date >= ?", req.DateFrom)
	}
	if req.DateTo != nil {
		query = query.Where("delivery_date <= ?", req.DateTo)
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
	if err := query.Offset(offset).Limit(req.Limit).Order("delivery_date DESC, created_at DESC").Find(&deliveryOrders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch delivery orders: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       deliveryOrders,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetDeliveryOrder returns a single delivery order by ID
func GetDeliveryOrder(ctx *context.Context, id string) (*domain.DeliveryOrder, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var deliveryOrder domain.DeliveryOrder
	if err := db.Preload("Items").
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&deliveryOrder).Error; err != nil {
		return nil, fmt.Errorf("delivery order not found: %w", err)
	}

	return &deliveryOrder, nil
}

// CreateDeliveryOrder creates a new delivery order
func CreateDeliveryOrder(ctx *context.Context, req domain.CreateDeliveryOrderDTO) (*domain.DeliveryOrder, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.DeliveryOrder
	if err := db.Where("code = ? AND tenant_id = ?", req.Code, tenantID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("delivery order with code %s already exists", req.Code)
	}

	// Get sales order info
	var order domain.SalesOrder
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ? AND active = ?", req.OrderID, tenantID, true).First(&order).Error; err != nil {
		return nil, fmt.Errorf("sales order not found: %w", err)
	}

	// Get warehouse info
	var warehouse domain.Warehouse
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.WarehouseID, tenantID, true).First(&warehouse).Error; err != nil {
		return nil, fmt.Errorf("warehouse not found: %w", err)
	}

	// Process delivery items
	items := make([]domain.DeliveryOrderItem, len(req.Items))
	orderItemUpdates := make(map[string]float64) // Track delivered qty updates per order item

	for i, itemReq := range req.Items {
		// Find the order item
		var orderItem domain.SalesOrderItem
		found := false
		for _, oi := range order.Items {
			if oi.ID == itemReq.OrderItemID {
				orderItem = oi
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("order item %s not found in sales order", itemReq.OrderItemID)
		}

		// Validate product ID matches
		if orderItem.ProductID != itemReq.ProductID {
			return nil, fmt.Errorf("product ID mismatch for order item %s", itemReq.OrderItemID)
		}

		// Check if delivered quantity exceeds remaining quantity
		remainingQty := orderItem.Quantity - orderItem.DeliveredQty
		if itemReq.DeliveredQty > remainingQty {
			return nil, fmt.Errorf("delivered quantity %.2f exceeds remaining quantity %.2f for product %s",
				itemReq.DeliveredQty, remainingQty, orderItem.ProductName)
		}

		items[i] = domain.DeliveryOrderItem{
			ID:            uuid.New().String(),
			TenantID:      tenantID,
			OrderItemID:   itemReq.OrderItemID,
			LineNumber:    i + 1,
			ProductID:     orderItem.ProductID,
			ProductCode:   orderItem.ProductCode,
			ProductName:   orderItem.ProductName,
			Description:   orderItem.Description,
			OrderedQty:    orderItem.Quantity,
			DeliveredQty:  itemReq.DeliveredQty,
			UOM:           orderItem.UOM,
			SerialNumbers: itemReq.SerialNumbers,
			BatchNumbers:  itemReq.BatchNumbers,
			Notes:         itemReq.Notes,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Track the delivered quantity for this order item
		orderItemUpdates[itemReq.OrderItemID] = orderItem.DeliveredQty + itemReq.DeliveredQty
	}

	deliveryOrder := &domain.DeliveryOrder{
		ID:              uuid.New().String(),
		Code:            req.Code,
		TenantID:        tenantID,
		OrderID:         order.ID,
		OrderCode:       order.Code,
		CustomerID:      order.CustomerID,
		CustomerName:    order.CustomerName,
		DeliveryDate:    req.DeliveryDate,
		ScheduledDate:   req.ScheduledDate,
		Status:          "draft",
		ReferenceNo:     req.ReferenceNo,
		TrackingNumber:  req.TrackingNumber,
		WarehouseID:     warehouse.ID,
		WarehouseName:   warehouse.Name,
		ShippingMethod:  req.ShippingMethod,
		ShippingCost:    req.ShippingCost,
		CourierName:     req.CourierName,
		DriverName:      req.DriverName,
		VehicleNumber:   req.VehicleNumber,
		ShippingAddress: req.ShippingAddress,
		ShippingCity:    req.ShippingCity,
		ShippingState:   req.ShippingState,
		ShippingZip:     req.ShippingZip,
		ShippingCountry: req.ShippingCountry,
		RecipientName:   req.RecipientName,
		RecipientPhone:  req.RecipientPhone,
		RecipientEmail:  req.RecipientEmail,
		Notes:           req.Notes,
		InternalNotes:   req.InternalNotes,
		Tags:            req.Tags,
		Items:           items,
		Active:          true,
		CreatedBy:       userID,
		UpdatedBy:       userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Start transaction to create delivery order and update order items
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(deliveryOrder).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create delivery order: %w", err)
	}

	// Update sales order items with delivered quantities
	for orderItemID, newDeliveredQty := range orderItemUpdates {
		if err := tx.Model(&domain.SalesOrderItem{}).
			Where("id = ? AND tenant_id = ?", orderItemID, tenantID).
			Update("delivered_qty", newDeliveredQty).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update order item delivered quantity: %w", err)
		}
	}

	// Update sales order fulfillment status
	var totalOrderedQty, totalDeliveredQty float64
	for _, item := range order.Items {
		totalOrderedQty += item.Quantity
		if updatedQty, exists := orderItemUpdates[item.ID]; exists {
			totalDeliveredQty += updatedQty
		} else {
			totalDeliveredQty += item.DeliveredQty
		}
	}

	var fulfillmentStatus string
	if totalDeliveredQty == 0 {
		fulfillmentStatus = "pending"
	} else if totalDeliveredQty < totalOrderedQty {
		fulfillmentStatus = "partial"
	} else {
		fulfillmentStatus = "completed"
	}

	if err := tx.Model(&domain.SalesOrder{}).
		Where("id = ? AND tenant_id = ?", order.ID, tenantID).
		Updates(map[string]interface{}{
			"delivered_qty":      totalDeliveredQty,
			"fulfillment_status": fulfillmentStatus,
			"updated_at":         time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update sales order fulfillment status: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return deliveryOrder, nil
}

// UpdateDeliveryOrder updates an existing delivery order
func UpdateDeliveryOrder(ctx *context.Context, id string, req domain.UpdateDeliveryOrderDTO) (*domain.DeliveryOrder, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var deliveryOrder domain.DeliveryOrder
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", id, tenantID).First(&deliveryOrder).Error; err != nil {
		return nil, fmt.Errorf("delivery order not found: %w", err)
	}

	// Only draft or ready delivery orders can be updated
	if deliveryOrder.Status != "draft" && deliveryOrder.Status != "ready" && req.Status == "" {
		return nil, fmt.Errorf("only draft or ready delivery orders can be updated")
	}

	// Update fields if provided
	if req.ScheduledDate != nil {
		deliveryOrder.ScheduledDate = *req.ScheduledDate
	}
	if req.ActualDate != nil {
		deliveryOrder.ActualDate = req.ActualDate
	}
	if req.Status != "" {
		deliveryOrder.Status = req.Status
		// Update workflow timestamps
		now := time.Now()
		switch req.Status {
		case "delivered":
			if deliveryOrder.ActualDate == nil {
				deliveryOrder.ActualDate = &now
			}
			if req.ReceivedBy != "" {
				deliveryOrder.ReceivedBy = req.ReceivedBy
				deliveryOrder.ReceivedAt = &now
			}
		}
	}
	if req.ReferenceNo != "" {
		deliveryOrder.ReferenceNo = req.ReferenceNo
	}
	if req.TrackingNumber != "" {
		deliveryOrder.TrackingNumber = req.TrackingNumber
	}
	if req.ShippingMethod != "" {
		deliveryOrder.ShippingMethod = req.ShippingMethod
	}
	if req.ShippingCost != nil {
		deliveryOrder.ShippingCost = *req.ShippingCost
	}
	if req.CourierName != "" {
		deliveryOrder.CourierName = req.CourierName
	}
	if req.DriverName != "" {
		deliveryOrder.DriverName = req.DriverName
	}
	if req.VehicleNumber != "" {
		deliveryOrder.VehicleNumber = req.VehicleNumber
	}
	if req.RecipientName != "" {
		deliveryOrder.RecipientName = req.RecipientName
	}
	if req.RecipientPhone != "" {
		deliveryOrder.RecipientPhone = req.RecipientPhone
	}
	if req.RecipientEmail != "" {
		deliveryOrder.RecipientEmail = req.RecipientEmail
	}
	if req.Notes != "" {
		deliveryOrder.Notes = req.Notes
	}
	if req.InternalNotes != "" {
		deliveryOrder.InternalNotes = req.InternalNotes
	}
	if req.Tags != nil {
		deliveryOrder.Tags = req.Tags
	}
	if req.ReceivedBy != "" {
		deliveryOrder.ReceivedBy = req.ReceivedBy
	}
	if req.ReceiverSignature != "" {
		deliveryOrder.ReceiverSignature = req.ReceiverSignature
	}
	if req.DeliveryProof != nil {
		deliveryOrder.DeliveryProof = req.DeliveryProof
	}

	deliveryOrder.UpdatedBy = userID
	deliveryOrder.UpdatedAt = time.Now()

	if err := db.Save(&deliveryOrder).Error; err != nil {
		return nil, fmt.Errorf("failed to update delivery order: %w", err)
	}

	return &deliveryOrder, nil
}

// DeleteDeliveryOrder soft deletes a delivery order
func DeleteDeliveryOrder(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var deliveryOrder domain.DeliveryOrder
	if err := db.Preload("Items").Where("id = ? AND tenant_id = ?", id, tenantID).First(&deliveryOrder).Error; err != nil {
		return fmt.Errorf("delivery order not found: %w", err)
	}

	// Only draft delivery orders can be deleted
	if deliveryOrder.Status != "draft" {
		return fmt.Errorf("only draft delivery orders can be deleted")
	}

	// Start transaction to delete delivery order and revert order item quantities
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Revert the delivered quantities in sales order items
	for _, item := range deliveryOrder.Items {
		if err := tx.Model(&domain.SalesOrderItem{}).
			Where("id = ? AND tenant_id = ?", item.OrderItemID, tenantID).
			Update("delivered_qty", db.Raw("delivered_qty - ?", item.DeliveredQty)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to revert order item delivered quantity: %w", err)
		}
	}

	// Recalculate sales order fulfillment status
	var order domain.SalesOrder
	if err := tx.Preload("Items").Where("id = ? AND tenant_id = ?", deliveryOrder.OrderID, tenantID).First(&order).Error; err == nil {
		var totalOrderedQty, totalDeliveredQty float64
		for _, item := range order.Items {
			totalOrderedQty += item.Quantity
			// Get updated delivered qty after revert
			var updatedItem domain.SalesOrderItem
			if err := tx.Where("id = ?", item.ID).First(&updatedItem).Error; err == nil {
				totalDeliveredQty += updatedItem.DeliveredQty
			}
		}

		var fulfillmentStatus string
		if totalDeliveredQty == 0 {
			fulfillmentStatus = "pending"
		} else if totalDeliveredQty < totalOrderedQty {
			fulfillmentStatus = "partial"
		} else {
			fulfillmentStatus = "completed"
		}

		if err := tx.Model(&domain.SalesOrder{}).
			Where("id = ? AND tenant_id = ?", order.ID, tenantID).
			Updates(map[string]interface{}{
				"delivered_qty":      totalDeliveredQty,
				"fulfillment_status": fulfillmentStatus,
				"updated_at":         time.Now(),
			}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update sales order fulfillment status: %w", err)
		}
	}

	// Soft delete the delivery order
	deliveryOrder.Active = false
	deliveryOrder.UpdatedBy = userID
	deliveryOrder.UpdatedAt = time.Now()

	if err := tx.Save(&deliveryOrder).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete delivery order: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetDeliveryOrderStats returns delivery order statistics
func GetDeliveryOrderStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalDeliveryOrders int64
	db.Model(&domain.DeliveryOrder{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalDeliveryOrders)

	var statusStats []struct {
		Status string
		Count  int64
	}
	db.Model(&domain.DeliveryOrder{}).
		Select("status, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("status").
		Scan(&statusStats)

	var todayDeliveries int64
	today := time.Now().Format("2006-01-02")
	db.Model(&domain.DeliveryOrder{}).
		Where("tenant_id = ? AND active = ? AND DATE(delivery_date) = ?", tenantID, true, today).
		Count(&todayDeliveries)

	var pendingDeliveries int64
	db.Model(&domain.DeliveryOrder{}).
		Where("tenant_id = ? AND active = ? AND status IN ?", tenantID, true, []string{"draft", "ready"}).
		Count(&pendingDeliveries)

	var completedDeliveries int64
	db.Model(&domain.DeliveryOrder{}).
		Where("tenant_id = ? AND active = ? AND status = ?", tenantID, true, "delivered").
		Count(&completedDeliveries)

	return map[string]interface{}{
		"total_delivery_orders": totalDeliveryOrders,
		"deliveries_by_status":  statusStats,
		"today_deliveries":      todayDeliveries,
		"pending_deliveries":    pendingDeliveries,
		"completed_deliveries":  completedDeliveries,
	}, nil
}
