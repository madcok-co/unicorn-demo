package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListCustomers returns paginated list of customers
func ListCustomers(ctx *context.Context, req domain.ListCustomersDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var customers []domain.Customer
	var total int64

	query := db.Model(&domain.Customer{}).Where("tenant_id = ? AND active = ?", tenantID, true)

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
	if err := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch customers: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       customers,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetCustomer returns a single customer by ID
func GetCustomer(ctx *context.Context, id string) (*domain.Customer, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var customer domain.Customer
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&customer).Error; err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	return &customer, nil
}

// CreateCustomer creates a new customer
func CreateCustomer(ctx *context.Context, req domain.CreateCustomerDTO) (*domain.Customer, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.Customer
	if err := db.Where("code = ? AND tenant_id = ? AND active = ?", req.Code, tenantID, true).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("customer with code %s already exists", req.Code)
	}

	customer := &domain.Customer{
		ID:               uuid.New().String(),
		Code:             req.Code,
		Name:             req.Name,
		Type:             req.Type,
		Email:            req.Email,
		Phone:            req.Phone,
		Mobile:           req.Mobile,
		TaxID:            req.TaxID,
		Website:          "",
		BillingAddress:   req.BillingAddress,
		BillingCity:      req.BillingCity,
		BillingState:     req.BillingState,
		BillingZip:       req.BillingZip,
		BillingCountry:   req.BillingCountry,
		ShippingAddress:  req.ShippingAddress,
		ShippingCity:     req.ShippingCity,
		ShippingState:    req.ShippingState,
		ShippingZip:      req.ShippingZip,
		ShippingCountry:  req.ShippingCountry,
		ContactPerson:    req.ContactPerson,
		ContactPersonJob: req.ContactPersonJob,
		CreditLimit:      req.CreditLimit,
		PaymentTermDays:  req.PaymentTermDays,
		Currency:         req.Currency,
		Notes:            req.Notes,
		Tags:             req.Tags,
		TenantID:         tenantID,
		Active:           true,
		CreatedBy:        userID,
		UpdatedBy:        userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.Create(customer).Error; err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

// UpdateCustomer updates an existing customer
func UpdateCustomer(ctx *context.Context, id string, req domain.UpdateCustomerDTO) (*domain.Customer, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var customer domain.Customer
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&customer).Error; err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.Phone != "" {
		customer.Phone = req.Phone
	}
	if req.Mobile != "" {
		customer.Mobile = req.Mobile
	}
	if req.TaxID != "" {
		customer.TaxID = req.TaxID
	}
	if req.Type != "" {
		customer.Type = req.Type
	}
	if req.BillingAddress != "" {
		customer.BillingAddress = req.BillingAddress
	}
	if req.BillingCity != "" {
		customer.BillingCity = req.BillingCity
	}
	if req.BillingState != "" {
		customer.BillingState = req.BillingState
	}
	if req.BillingZip != "" {
		customer.BillingZip = req.BillingZip
	}
	if req.BillingCountry != "" {
		customer.BillingCountry = req.BillingCountry
	}
	if req.ShippingAddress != "" {
		customer.ShippingAddress = req.ShippingAddress
	}
	if req.ShippingCity != "" {
		customer.ShippingCity = req.ShippingCity
	}
	if req.ShippingState != "" {
		customer.ShippingState = req.ShippingState
	}
	if req.ShippingZip != "" {
		customer.ShippingZip = req.ShippingZip
	}
	if req.ShippingCountry != "" {
		customer.ShippingCountry = req.ShippingCountry
	}
	if req.CreditLimit != nil {
		customer.CreditLimit = *req.CreditLimit
	}
	if req.PaymentTermDays != nil {
		customer.PaymentTermDays = *req.PaymentTermDays
	}
	if req.Currency != "" {
		customer.Currency = req.Currency
	}
	if req.ContactPerson != "" {
		customer.ContactPerson = req.ContactPerson
	}
	if req.ContactPersonJob != "" {
		customer.ContactPersonJob = req.ContactPersonJob
	}
	if req.Notes != "" {
		customer.Notes = req.Notes
	}
	if req.Tags != nil {
		customer.Tags = req.Tags
	}
	if req.Active != nil {
		customer.Active = *req.Active
	}

	customer.UpdatedBy = userID
	customer.UpdatedAt = time.Now()

	if err := db.Save(&customer).Error; err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return &customer, nil
}

// DeleteCustomer soft deletes a customer
func DeleteCustomer(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var customer domain.Customer
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&customer).Error; err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}

	customer.Active = false
	customer.UpdatedBy = userID
	customer.UpdatedAt = time.Now()

	if err := db.Save(&customer).Error; err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	return nil
}

// GetCustomerStats returns customer statistics
func GetCustomerStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalCustomers int64
	db.Model(&domain.Customer{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalCustomers)

	var typeStats []struct {
		Type  string
		Count int64
	}
	db.Model(&domain.Customer{}).
		Select("type, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("type").
		Scan(&typeStats)

	return map[string]interface{}{
		"total_customers":   totalCustomers,
		"customers_by_type": typeStats,
	}, nil
}
