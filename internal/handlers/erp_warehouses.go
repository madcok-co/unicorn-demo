package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListWarehouses returns paginated list of warehouses
func ListWarehouses(ctx *context.Context, req domain.ListWarehousesDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var warehouses []domain.Warehouse
	var total int64

	query := db.Model(&domain.Warehouse{}).Where("tenant_id = ? AND active = ?", tenantID, true)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR city LIKE ?", search, search, search)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Active != nil {
		query = query.Where("active = ?", *req.Active)
	}

	// Count total
	query.Count(&total)

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&warehouses).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch warehouses: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       warehouses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetWarehouse returns a single warehouse by ID
func GetWarehouse(ctx *context.Context, id string) (*domain.Warehouse, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var warehouse domain.Warehouse
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&warehouse).Error; err != nil {
		return nil, fmt.Errorf("warehouse not found: %w", err)
	}

	return &warehouse, nil
}

// CreateWarehouse creates a new warehouse
func CreateWarehouse(ctx *context.Context, req domain.CreateWarehouseDTO) (*domain.Warehouse, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.Warehouse
	if err := db.Where("code = ? AND tenant_id = ? AND active = ?", req.Code, tenantID, true).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("warehouse with code %s already exists", req.Code)
	}

	warehouse := &domain.Warehouse{
		ID:                 uuid.New().String(),
		Code:               req.Code,
		Name:               req.Name,
		Type:               req.Type,
		Address:            req.Address,
		City:               req.City,
		State:              req.State,
		Zip:                req.Zip,
		Country:            req.Country,
		Phone:              req.Phone,
		Email:              req.Email,
		ManagerID:          req.ManagerID,
		ManagerName:        req.ManagerName,
		TotalArea:          req.TotalArea,
		UsedArea:           0,
		TotalCapacity:      req.TotalCapacity,
		UsedCapacity:       0,
		AllowNegativeStock: req.AllowNegativeStock,
		IsDefault:          req.IsDefault,
		Description:        req.Description,
		Notes:              req.Notes,
		TenantID:           tenantID,
		Active:             true,
		CreatedBy:          userID,
		UpdatedBy:          userID,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := db.Create(warehouse).Error; err != nil {
		return nil, fmt.Errorf("failed to create warehouse: %w", err)
	}

	return warehouse, nil
}

// UpdateWarehouse updates an existing warehouse
func UpdateWarehouse(ctx *context.Context, id string, req domain.UpdateWarehouseDTO) (*domain.Warehouse, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var warehouse domain.Warehouse
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&warehouse).Error; err != nil {
		return nil, fmt.Errorf("warehouse not found: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		warehouse.Name = req.Name
	}
	if req.Type != "" {
		warehouse.Type = req.Type
	}
	if req.Address != "" {
		warehouse.Address = req.Address
	}
	if req.City != "" {
		warehouse.City = req.City
	}
	if req.State != "" {
		warehouse.State = req.State
	}
	if req.Zip != "" {
		warehouse.Zip = req.Zip
	}
	if req.Country != "" {
		warehouse.Country = req.Country
	}
	if req.Phone != "" {
		warehouse.Phone = req.Phone
	}
	if req.Email != "" {
		warehouse.Email = req.Email
	}
	if req.ManagerID != "" {
		warehouse.ManagerID = req.ManagerID
	}
	if req.ManagerName != "" {
		warehouse.ManagerName = req.ManagerName
	}
	if req.TotalArea != nil {
		warehouse.TotalArea = *req.TotalArea
	}
	if req.TotalCapacity != nil {
		warehouse.TotalCapacity = *req.TotalCapacity
	}
	if req.AllowNegativeStock != nil {
		warehouse.AllowNegativeStock = *req.AllowNegativeStock
	}
	if req.IsDefault != nil {
		warehouse.IsDefault = *req.IsDefault
	}
	if req.Active != nil {
		warehouse.Active = *req.Active
	}

	warehouse.UpdatedBy = userID
	warehouse.UpdatedAt = time.Now()

	if err := db.Save(&warehouse).Error; err != nil {
		return nil, fmt.Errorf("failed to update warehouse: %w", err)
	}

	return &warehouse, nil
}

// DeleteWarehouse soft deletes a warehouse
func DeleteWarehouse(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var warehouse domain.Warehouse
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&warehouse).Error; err != nil {
		return fmt.Errorf("warehouse not found: %w", err)
	}

	warehouse.Active = false
	warehouse.UpdatedBy = userID
	warehouse.UpdatedAt = time.Now()

	if err := db.Save(&warehouse).Error; err != nil {
		return fmt.Errorf("failed to delete warehouse: %w", err)
	}

	return nil
}

// GetWarehouseStats returns warehouse statistics
func GetWarehouseStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalWarehouses int64
	db.Model(&domain.Warehouse{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalWarehouses)

	var typeStats []struct {
		Type  string
		Count int64
	}
	db.Model(&domain.Warehouse{}).
		Select("type, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("type").
		Scan(&typeStats)

	var totalArea, totalCapacity float64
	db.Model(&domain.Warehouse{}).
		Select("COALESCE(SUM(total_area), 0) as total_area, COALESCE(SUM(total_capacity), 0) as total_capacity").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Row().Scan(&totalArea, &totalCapacity)

	return map[string]interface{}{
		"total_warehouses":   totalWarehouses,
		"warehouses_by_type": typeStats,
		"total_area":         totalArea,
		"total_capacity":     totalCapacity,
	}, nil
}
