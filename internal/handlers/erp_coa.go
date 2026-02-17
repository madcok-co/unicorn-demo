package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListChartOfAccounts returns paginated list of chart of accounts
func ListChartOfAccounts(ctx *context.Context, req domain.ListChartOfAccountsDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var accounts []domain.ChartOfAccount
	var total int64

	query := db.Model(&domain.ChartOfAccount{}).Where("tenant_id = ? AND active = ?", tenantID, true)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("name LIKE ? OR code LIKE ?", search, search)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.ParentID != nil {
		if *req.ParentID == "" {
			query = query.Where("(parent_id IS NULL OR parent_id = '')")
		} else {
			query = query.Where("parent_id = ?", *req.ParentID)
		}
	}
	if req.Active != nil {
		query = query.Where("active = ?", *req.Active)
	}
	if req.IsGroup != nil {
		query = query.Where("is_group = ?", *req.IsGroup)
	}

	// Count total
	query.Count(&total)

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("code ASC").Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch chart of accounts: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       accounts,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetChartOfAccount returns a single chart of account by ID
func GetChartOfAccount(ctx *context.Context, id string) (*domain.ChartOfAccount, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var account domain.ChartOfAccount
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&account).Error; err != nil {
		return nil, fmt.Errorf("chart of account not found: %w", err)
	}

	return &account, nil
}

// CreateChartOfAccount creates a new chart of account
func CreateChartOfAccount(ctx *context.Context, req domain.CreateChartOfAccountDTO) (*domain.ChartOfAccount, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.ChartOfAccount
	if err := db.Where("code = ? AND tenant_id = ? AND active = ?", req.Code, tenantID, true).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("chart of account with code %s already exists", req.Code)
	}

	// Calculate level based on parent
	level := 1
	if req.ParentID != "" {
		var parent domain.ChartOfAccount
		if err := db.Where("id = ? AND tenant_id = ? AND active = ?", req.ParentID, tenantID, true).First(&parent).Error; err != nil {
			return nil, fmt.Errorf("parent account not found: %w", err)
		}
		level = parent.Level + 1
	}

	account := &domain.ChartOfAccount{
		ID:                  uuid.New().String(),
		Code:                req.Code,
		Name:                req.Name,
		Type:                req.Type,
		Category:            req.Category,
		ParentID:            req.ParentID,
		Level:               level,
		IsGroup:             req.IsGroup,
		Currency:            req.Currency,
		CurrentBalance:      0,
		DebitBalance:        0,
		CreditBalance:       0,
		AllowReconciliation: req.AllowReconciliation,
		RequireTaxReporting: req.RequireTaxReporting,
		Description:         req.Description,
		Notes:               req.Notes,
		TenantID:            tenantID,
		Active:              true,
		CreatedBy:           userID,
		UpdatedBy:           userID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := db.Create(account).Error; err != nil {
		return nil, fmt.Errorf("failed to create chart of account: %w", err)
	}

	return account, nil
}

// UpdateChartOfAccount updates an existing chart of account
func UpdateChartOfAccount(ctx *context.Context, id string, req domain.UpdateChartOfAccountDTO) (*domain.ChartOfAccount, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var account domain.ChartOfAccount
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&account).Error; err != nil {
		return nil, fmt.Errorf("chart of account not found: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		account.Name = req.Name
	}
	if req.Category != "" {
		account.Category = req.Category
	}
	if req.Currency != "" {
		account.Currency = req.Currency
	}
	if req.AllowReconciliation != nil {
		account.AllowReconciliation = *req.AllowReconciliation
	}
	if req.RequireTaxReporting != nil {
		account.RequireTaxReporting = *req.RequireTaxReporting
	}
	if req.Description != "" {
		account.Description = req.Description
	}
	if req.Notes != "" {
		account.Notes = req.Notes
	}
	if req.Active != nil {
		account.Active = *req.Active
	}

	account.UpdatedBy = userID
	account.UpdatedAt = time.Now()

	if err := db.Save(&account).Error; err != nil {
		return nil, fmt.Errorf("failed to update chart of account: %w", err)
	}

	return &account, nil
}

// DeleteChartOfAccount soft deletes a chart of account
func DeleteChartOfAccount(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var account domain.ChartOfAccount
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&account).Error; err != nil {
		return fmt.Errorf("chart of account not found: %w", err)
	}

	// Check if account has children
	var childCount int64
	db.Model(&domain.ChartOfAccount{}).
		Where("parent_id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).
		Count(&childCount)
	if childCount > 0 {
		return fmt.Errorf("cannot delete account with active children")
	}

	account.Active = false
	account.UpdatedBy = userID
	account.UpdatedAt = time.Now()

	if err := db.Save(&account).Error; err != nil {
		return fmt.Errorf("failed to delete chart of account: %w", err)
	}

	return nil
}

// GetChartOfAccountStats returns chart of account statistics
func GetChartOfAccountStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalAccounts int64
	db.Model(&domain.ChartOfAccount{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalAccounts)

	var typeStats []struct {
		Type  string
		Count int64
	}
	db.Model(&domain.ChartOfAccount{}).
		Select("type, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("type").
		Scan(&typeStats)

	var groupAccounts int64
	db.Model(&domain.ChartOfAccount{}).
		Where("tenant_id = ? AND active = ? AND is_group = ?", tenantID, true, true).
		Count(&groupAccounts)

	var totalDebit, totalCredit float64
	db.Model(&domain.ChartOfAccount{}).
		Select("COALESCE(SUM(debit_balance), 0) as total_debit, COALESCE(SUM(credit_balance), 0) as total_credit").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Row().Scan(&totalDebit, &totalCredit)

	return map[string]interface{}{
		"total_accounts":   totalAccounts,
		"accounts_by_type": typeStats,
		"group_accounts":   groupAccounts,
		"total_debit":      totalDebit,
		"total_credit":     totalCredit,
	}, nil
}
