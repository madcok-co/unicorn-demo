package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListVendors returns paginated list of vendors
func ListVendors(ctx *context.Context, req domain.ListVendorsDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var vendors []domain.Vendor
	var total int64

	query := db.Model(&domain.Vendor{}).Where("tenant_id = ? AND active = ?", tenantID, true)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR email LIKE ?", search, search, search)
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
	if err := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&vendors).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch vendors: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       vendors,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetVendor returns a single vendor by ID
func GetVendor(ctx *context.Context, id string) (*domain.Vendor, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var vendor domain.Vendor
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&vendor).Error; err != nil {
		return nil, fmt.Errorf("vendor not found: %w", err)
	}

	return &vendor, nil
}

// CreateVendor creates a new vendor
func CreateVendor(ctx *context.Context, req domain.CreateVendorDTO) (*domain.Vendor, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.Vendor
	if err := db.Where("code = ? AND tenant_id = ? AND active = ?", req.Code, tenantID, true).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("vendor with code %s already exists", req.Code)
	}

	vendor := &domain.Vendor{
		ID:               uuid.New().String(),
		Code:             req.Code,
		Name:             req.Name,
		Type:             req.Type,
		Email:            req.Email,
		Phone:            req.Phone,
		Mobile:           req.Mobile,
		Website:          "",
		TaxID:            req.TaxID,
		Address:          req.Address,
		City:             req.City,
		State:            req.State,
		Zip:              req.Zip,
		Country:          req.Country,
		ContactPerson:    req.ContactPerson,
		ContactPersonJob: req.ContactPersonJob,
		PaymentTermDays:  req.PaymentTermDays,
		Currency:         req.Currency,
		BankName:         req.BankName,
		BankAccount:      req.BankAccount,
		BankAccountName:  req.BankAccountName,
		Rating:           req.Rating,
		Category:         req.Category,
		Notes:            req.Notes,
		Tags:             req.Tags,
		TenantID:         tenantID,
		Active:           true,
		CreatedBy:        userID,
		UpdatedBy:        userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.Create(vendor).Error; err != nil {
		return nil, fmt.Errorf("failed to create vendor: %w", err)
	}

	return vendor, nil
}

// UpdateVendor updates an existing vendor
func UpdateVendor(ctx *context.Context, id string, req domain.UpdateVendorDTO) (*domain.Vendor, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var vendor domain.Vendor
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&vendor).Error; err != nil {
		return nil, fmt.Errorf("vendor not found: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		vendor.Name = req.Name
	}
	if req.Email != "" {
		vendor.Email = req.Email
	}
	if req.Phone != "" {
		vendor.Phone = req.Phone
	}
	if req.Mobile != "" {
		vendor.Mobile = req.Mobile
	}
	if req.TaxID != "" {
		vendor.TaxID = req.TaxID
	}
	if req.Type != "" {
		vendor.Type = req.Type
	}
	if req.Address != "" {
		vendor.Address = req.Address
	}
	if req.City != "" {
		vendor.City = req.City
	}
	if req.State != "" {
		vendor.State = req.State
	}
	if req.Zip != "" {
		vendor.Zip = req.Zip
	}
	if req.Country != "" {
		vendor.Country = req.Country
	}
	if req.PaymentTermDays != nil {
		vendor.PaymentTermDays = *req.PaymentTermDays
	}
	if req.Currency != "" {
		vendor.Currency = req.Currency
	}
	if req.BankName != "" {
		vendor.BankName = req.BankName
	}
	if req.BankAccount != "" {
		vendor.BankAccount = req.BankAccount
	}
	if req.BankAccountName != "" {
		vendor.BankAccountName = req.BankAccountName
	}
	if req.ContactPerson != "" {
		vendor.ContactPerson = req.ContactPerson
	}
	if req.ContactPersonJob != "" {
		vendor.ContactPersonJob = req.ContactPersonJob
	}
	if req.Rating != nil {
		vendor.Rating = *req.Rating
	}
	if req.Category != "" {
		vendor.Category = req.Category
	}
	if req.Notes != "" {
		vendor.Notes = req.Notes
	}
	if req.Tags != nil {
		vendor.Tags = req.Tags
	}
	if req.Active != nil {
		vendor.Active = *req.Active
	}

	vendor.UpdatedBy = userID
	vendor.UpdatedAt = time.Now()

	if err := db.Save(&vendor).Error; err != nil {
		return nil, fmt.Errorf("failed to update vendor: %w", err)
	}

	return &vendor, nil
}

// DeleteVendor soft deletes a vendor
func DeleteVendor(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var vendor domain.Vendor
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&vendor).Error; err != nil {
		return fmt.Errorf("vendor not found: %w", err)
	}

	vendor.Active = false
	vendor.UpdatedBy = userID
	vendor.UpdatedAt = time.Now()

	if err := db.Save(&vendor).Error; err != nil {
		return fmt.Errorf("failed to delete vendor: %w", err)
	}

	return nil
}

// GetVendorStats returns vendor statistics
func GetVendorStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalVendors int64
	db.Model(&domain.Vendor{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalVendors)

	var ratingStats []struct {
		Rating int
		Count  int64
	}
	db.Model(&domain.Vendor{}).
		Select("rating, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("rating").
		Scan(&ratingStats)

	return map[string]interface{}{
		"total_vendors":     totalVendors,
		"vendors_by_rating": ratingStats,
	}, nil
}
