package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListTaxes returns paginated list of taxes
func ListTaxes(ctx *context.Context, req domain.ListTaxesDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var taxes []domain.Tax
	var total int64

	query := db.Model(&domain.Tax{}).Where("tenant_id = ? AND active = ?", tenantID, true)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("name LIKE ? OR code LIKE ?", search, search)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.IsDefault != nil {
		query = query.Where("is_default = ?", *req.IsDefault)
	}

	// Count total
	query.Count(&total)

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&taxes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch taxes: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       taxes,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetTax returns a single tax by ID
func GetTax(ctx *context.Context, id string) (*domain.Tax, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var tax domain.Tax
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&tax).Error; err != nil {
		return nil, fmt.Errorf("tax not found: %w", err)
	}

	return &tax, nil
}

// CreateTax creates a new tax
func CreateTax(ctx *context.Context, req domain.CreateTaxDTO) (*domain.Tax, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.Tax
	if err := db.Where("code = ? AND tenant_id = ? AND active = ?", req.Code, tenantID, true).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("tax with code %s already exists", req.Code)
	}

	tax := &domain.Tax{
		ID:              uuid.New().String(),
		Code:            req.Code,
		Name:            req.Name,
		Type:            req.Type,
		Rate:            req.Rate,
		TaxAccountID:    "",
		IsDefault:       req.IsDefault,
		IncludedInPrice: false,
		Description:     req.Description,
		TenantID:        tenantID,
		Active:          true,
		CreatedBy:       userID,
		UpdatedBy:       userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := db.Create(tax).Error; err != nil {
		return nil, fmt.Errorf("failed to create tax: %w", err)
	}

	return tax, nil
}

// UpdateTax updates an existing tax
func UpdateTax(ctx *context.Context, id string, req domain.UpdateTaxDTO) (*domain.Tax, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var tax domain.Tax
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&tax).Error; err != nil {
		return nil, fmt.Errorf("tax not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		tax.Name = *req.Name
	}
	if req.Type != nil {
		tax.Type = *req.Type
	}
	if req.Rate != nil {
		tax.Rate = *req.Rate
	}
	if req.IsDefault != nil {
		tax.IsDefault = *req.IsDefault
	}
	if req.Description != nil {
		tax.Description = *req.Description
	}

	tax.UpdatedBy = userID
	tax.UpdatedAt = time.Now()

	if err := db.Save(&tax).Error; err != nil {
		return nil, fmt.Errorf("failed to update tax: %w", err)
	}

	return &tax, nil
}

// DeleteTax soft deletes a tax
func DeleteTax(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var tax domain.Tax
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&tax).Error; err != nil {
		return fmt.Errorf("tax not found: %w", err)
	}

	tax.Active = false
	tax.UpdatedBy = userID
	tax.UpdatedAt = time.Now()

	if err := db.Save(&tax).Error; err != nil {
		return fmt.Errorf("failed to delete tax: %w", err)
	}

	return nil
}

// GetTaxStats returns tax statistics
func GetTaxStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalTaxes int64
	db.Model(&domain.Tax{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalTaxes)

	var typeStats []struct {
		Type  string
		Count int64
	}
	db.Model(&domain.Tax{}).
		Select("type, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("type").
		Scan(&typeStats)

	var avgRate float64
	db.Model(&domain.Tax{}).
		Select("COALESCE(AVG(rate), 0) as avg_rate").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Row().Scan(&avgRate)

	return map[string]interface{}{
		"total_taxes":   totalTaxes,
		"taxes_by_type": typeStats,
		"average_rate":  avgRate,
	}, nil
}
