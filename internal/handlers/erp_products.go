package handlers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListProducts returns paginated list of products
func ListProducts(ctx *context.Context, req domain.ListProductsDTO) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var products []domain.Product
	var total int64

	query := db.Model(&domain.Product{}).Where("tenant_id = ? AND active = ?", tenantID, true)

	// Apply filters
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR barcode LIKE ? OR sku LIKE ?", search, search, search, search)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Active != nil {
		query = query.Where("active = ?", *req.Active)
	}
	if req.CanBeSold != nil {
		query = query.Where("can_be_sold = ?", *req.CanBeSold)
	}
	if req.CanBePurchased != nil {
		query = query.Where("can_be_purchased = ?", *req.CanBePurchased)
	}
	if req.TrackInventory != nil {
		query = query.Where("track_inventory = ?", *req.TrackInventory)
	}

	// Count total
	query.Count(&total)

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	if err := query.Offset(offset).Limit(req.Limit).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	return &pagination.OffsetResult{
		Data:       products,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: int((total + int64(req.Limit) - 1) / int64(req.Limit)),
	}, nil
}

// GetProduct returns a single product by ID
func GetProduct(ctx *context.Context, id string) (*domain.Product, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var product domain.Product
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&product).Error; err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	return &product, nil
}

// CreateProduct creates a new product
func CreateProduct(ctx *context.Context, req domain.CreateProductDTO) (*domain.Product, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	// Check duplicate code
	var existing domain.Product
	if err := db.Where("code = ? AND tenant_id = ? AND active = ?", req.Code, tenantID, true).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("product with code %s already exists", req.Code)
	}

	product := &domain.Product{
		ID:               uuid.New().String(),
		Code:             req.Code,
		Name:             req.Name,
		Type:             req.Type,
		Category:         req.Category,
		UOM:              req.UOM,
		Barcode:          req.Barcode,
		SKU:              req.SKU,
		CanBeSold:        req.CanBeSold,
		CanBePurchased:   req.CanBePurchased,
		TrackInventory:   req.TrackInventory,
		MinStock:         req.MinStock,
		MaxStock:         req.MaxStock,
		ReorderLevel:     req.ReorderLevel,
		SalePrice:        req.SalePrice,
		PurchasePrice:    req.PurchasePrice,
		Cost:             req.Cost,
		Currency:         req.Currency,
		TaxCategory:      req.TaxCategory,
		TaxRate:          req.TaxRate,
		Weight:           req.Weight,
		Volume:           req.Volume,
		Length:           req.Length,
		Width:            req.Width,
		Height:           req.Height,
		WarrantyDays:     req.WarrantyDays,
		Description:      req.Description,
		InternalNotes:    req.InternalNotes,
		ImageURL:         req.ImageURL,
		AdditionalImages: req.AdditionalImages,
		IncomeAccountID:  req.IncomeAccountID,
		ExpenseAccountID: req.ExpenseAccountID,
		AssetAccountID:   req.AssetAccountID,
		Tags:             req.Tags,
		TenantID:         tenantID,
		Active:           true,
		CreatedBy:        userID,
		UpdatedBy:        userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := db.Create(product).Error; err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// UpdateProduct updates an existing product
func UpdateProduct(ctx *context.Context, id string, req domain.UpdateProductDTO) (*domain.Product, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, _ := getUserID(ctx)

	var product domain.Product
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&product).Error; err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Type != "" {
		product.Type = req.Type
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.UOM != "" {
		product.UOM = req.UOM
	}
	if req.Barcode != "" {
		product.Barcode = req.Barcode
	}
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.Active != nil {
		product.Active = *req.Active
	}
	if req.CanBeSold != nil {
		product.CanBeSold = *req.CanBeSold
	}
	if req.CanBePurchased != nil {
		product.CanBePurchased = *req.CanBePurchased
	}
	if req.TrackInventory != nil {
		product.TrackInventory = *req.TrackInventory
	}
	if req.MinStock != nil {
		product.MinStock = *req.MinStock
	}
	if req.MaxStock != nil {
		product.MaxStock = *req.MaxStock
	}
	if req.ReorderLevel != nil {
		product.ReorderLevel = *req.ReorderLevel
	}
	if req.SalePrice != nil {
		product.SalePrice = *req.SalePrice
	}
	if req.PurchasePrice != nil {
		product.PurchasePrice = *req.PurchasePrice
	}
	if req.Cost != nil {
		product.Cost = *req.Cost
	}

	product.UpdatedBy = userID
	product.UpdatedAt = time.Now()

	if err := db.Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &product, nil
}

// DeleteProduct soft deletes a product
func DeleteProduct(ctx *context.Context, id string) error {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return err
	}
	userID, _ := getUserID(ctx)

	var product domain.Product
	if err := db.Where("id = ? AND tenant_id = ? AND active = ?", id, tenantID, true).First(&product).Error; err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	product.Active = false
	product.UpdatedBy = userID
	product.UpdatedAt = time.Now()

	if err := db.Save(&product).Error; err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// GetProductStats returns product statistics
func GetProductStats(ctx *context.Context) (map[string]interface{}, error) {
	db := getDB(ctx)
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var totalProducts int64
	db.Model(&domain.Product{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&totalProducts)

	var typeStats []struct {
		Type  string
		Count int64
	}
	db.Model(&domain.Product{}).
		Select("type, COUNT(*) as count").
		Where("tenant_id = ? AND active = ?", tenantID, true).
		Group("type").
		Scan(&typeStats)

	var categoryStats []struct {
		Category string
		Count    int64
	}
	db.Model(&domain.Product{}).
		Select("category, COUNT(*) as count").
		Where("tenant_id = ? AND active = ? AND category != ''", tenantID, true).
		Group("category").
		Scan(&categoryStats)

	return map[string]interface{}{
		"total_products":       totalProducts,
		"products_by_type":     typeStats,
		"products_by_category": categoryStats,
	}, nil
}
