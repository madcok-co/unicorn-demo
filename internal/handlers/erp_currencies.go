package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListCurrencies returns paginated list of currencies
func ListCurrencies(ctx *context.Context, req domain.ListCurrenciesDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var currencies []domain.Currency
	var total int64

	query := db.Model(&domain.Currency{}).Where("tenant_id = ? AND active = ?", tenantID, true)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR symbol LIKE ?", search, search, search)
	}
	if req.IsDefault != nil {
		query = query.Where("is_default = ?", *req.IsDefault)
	}

	// Count total
	query.Count(&total)

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&currencies).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch currencies: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       currencies,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetCurrency returns a single currency by ID
func GetCurrency(ctx *context.Context, id string) (*domain.Currency, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var currency domain.Currency
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&currency).Error; err != nil {
		return nil, fmt.Errorf("currency not found: %w", err)
	}

	return &currency, nil
}

// CreateCurrency creates a new currency
func CreateCurrency(ctx *context.Context, req domain.CreateCurrencyDTO) (*domain.Currency, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.Currency
	if err := db.Where("code = ? AND tenant_id = ? AND active = ?", req.Code, tenantID, true).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("currency with code %s already exists", req.Code)
	}

	currency := &domain.Currency{
		ID:            uuid.New().String(),
		Code:          req.Code,
		Name:          req.Name,
		Symbol:        req.Symbol,
		ExchangeRate:  req.ExchangeRate,
		DecimalPlaces: req.DecimalPlaces,
		IsDefault:     req.IsDefault,
		TenantID:      tenantID,
		Active:        true,
		CreatedBy:     userID,
		UpdatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := db.Create(currency).Error; err != nil {
		return nil, fmt.Errorf("failed to create currency: %w", err)
	}

	return currency, nil
}

// UpdateCurrency updates an existing currency
func UpdateCurrency(ctx *context.Context, id string, req domain.UpdateCurrencyDTO) (*domain.Currency, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var currency domain.Currency
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&currency).Error; err != nil {
		return nil, fmt.Errorf("currency not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		currency.Name = *req.Name
	}
	if req.Symbol != nil {
		currency.Symbol = *req.Symbol
	}
	if req.ExchangeRate != nil {
		currency.ExchangeRate = *req.ExchangeRate
	}
	if req.DecimalPlaces != nil {
		currency.DecimalPlaces = *req.DecimalPlaces
	}
	if req.IsDefault != nil {
		currency.IsDefault = *req.IsDefault
	}

	currency.UpdatedBy = userID
	currency.UpdatedAt = time.Now()

	if err := db.Save(&currency).Error; err != nil {
		return nil, fmt.Errorf("failed to update currency: %w", err)
	}

	return &currency, nil
}

// DeleteCurrency soft deletes a currency
func DeleteCurrency(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var currency domain.Currency
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&currency).Error; err != nil {
		return fmt.Errorf("currency not found: %w", err)
	}

	currency.Active = false
	currency.UpdatedBy = userID
	currency.UpdatedAt = time.Now()

	if err := db.Save(&currency).Error; err != nil {
		return fmt.Errorf("failed to delete currency: %w", err)
	}

	return nil
}

// GetCurrencyStats returns currency statistics
func GetCurrencyStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalCurrencies int64
	db.Model(&domain.Currency{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalCurrencies)

	var avgExchangeRate float64
	db.Model(&domain.Currency{}).
		Select("COALESCE(AVG(exchange_rate), 0) as avg_rate").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Row().Scan(&avgExchangeRate)

	var defaultCurrency string
	db.Model(&domain.Currency{}).
		Select("code").
		Where("tenant_id = ? AND active = ? AND is_default = ?", tenantID, true, true).
		Row().Scan(&defaultCurrency)

	return map[string]interface{}{
		"total_currencies":      totalCurrencies,
		"average_exchange_rate": avgExchangeRate,
		"default_currency":      defaultCurrency,
	}, nil
}
