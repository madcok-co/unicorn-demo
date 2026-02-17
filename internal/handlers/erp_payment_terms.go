package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListPaymentTerms returns paginated list of payment terms
func ListPaymentTerms(ctx *context.Context, req domain.ListPaymentTermsDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var paymentTerms []domain.PaymentTerm
	var total int64

	query := db.Model(&domain.PaymentTerm{}).Where("tenant_id = ? AND active = ?", tenantID, true)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("name LIKE ? OR code LIKE ?", search, search)
	}
	if req.IsDefault != nil {
		query = query.Where("is_default = ?", *req.IsDefault)
	}

	// Count total
	query.Count(&total)

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&paymentTerms).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment terms: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       paymentTerms,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetPaymentTerm returns a single payment term by ID
func GetPaymentTerm(ctx *context.Context, id string) (*domain.PaymentTerm, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var paymentTerm domain.PaymentTerm
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&paymentTerm).Error; err != nil {
		return nil, fmt.Errorf("payment term not found: %w", err)
	}

	return &paymentTerm, nil
}

// CreatePaymentTerm creates a new payment term
func CreatePaymentTerm(ctx *context.Context, req domain.CreatePaymentTermDTO) (*domain.PaymentTerm, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.PaymentTerm
	if err := db.Where("code = ? AND tenant_id = ? AND active = ?", req.Code, tenantID, true).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("payment term with code %s already exists", req.Code)
	}

	paymentTerm := &domain.PaymentTerm{
		ID:              uuid.New().String(),
		Code:            req.Code,
		Name:            req.Name,
		Days:            req.Days,
		Type:            "fixed",
		DiscountDays:    req.DiscountDays,
		DiscountPercent: req.DiscountPercent,
		IsDefault:       req.IsDefault,
		Description:     req.Description,
		TenantID:        tenantID,
		Active:          true,
		CreatedBy:       userID,
		UpdatedBy:       userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := db.Create(paymentTerm).Error; err != nil {
		return nil, fmt.Errorf("failed to create payment term: %w", err)
	}

	return paymentTerm, nil
}

// UpdatePaymentTerm updates an existing payment term
func UpdatePaymentTerm(ctx *context.Context, id string, req domain.UpdatePaymentTermDTO) (*domain.PaymentTerm, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var paymentTerm domain.PaymentTerm
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&paymentTerm).Error; err != nil {
		return nil, fmt.Errorf("payment term not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		paymentTerm.Name = *req.Name
	}
	if req.Days != nil {
		paymentTerm.Days = *req.Days
	}
	if req.DiscountDays != nil {
		paymentTerm.DiscountDays = *req.DiscountDays
	}
	if req.DiscountPercent != nil {
		paymentTerm.DiscountPercent = *req.DiscountPercent
	}
	if req.IsDefault != nil {
		paymentTerm.IsDefault = *req.IsDefault
	}
	if req.Description != nil {
		paymentTerm.Description = *req.Description
	}

	paymentTerm.UpdatedBy = userID
	paymentTerm.UpdatedAt = time.Now()

	if err := db.Save(&paymentTerm).Error; err != nil {
		return nil, fmt.Errorf("failed to update payment term: %w", err)
	}

	return &paymentTerm, nil
}

// DeletePaymentTerm soft deletes a payment term
func DeletePaymentTerm(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var paymentTerm domain.PaymentTerm
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&paymentTerm).Error; err != nil {
		return fmt.Errorf("payment term not found: %w", err)
	}

	paymentTerm.Active = false
	paymentTerm.UpdatedBy = userID
	paymentTerm.UpdatedAt = time.Now()

	if err := db.Save(&paymentTerm).Error; err != nil {
		return fmt.Errorf("failed to delete payment term: %w", err)
	}

	return nil
}

// GetPaymentTermStats returns payment term statistics
func GetPaymentTermStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalPaymentTerms int64
	db.Model(&domain.PaymentTerm{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalPaymentTerms)

	var avgDays float64
	db.Model(&domain.PaymentTerm{}).
		Select("COALESCE(AVG(days), 0) as avg_days").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Row().Scan(&avgDays)

	var withDiscount int64
	db.Model(&domain.PaymentTerm{}).
		Where("tenant_id = ? AND active = ? AND discount_percent > 0", tenantID, true).
		Count(&withDiscount)

	return map[string]interface{}{
		"total_payment_terms": totalPaymentTerms,
		"average_days":        avgDays,
		"with_discount":       withDiscount,
	}, nil
}
