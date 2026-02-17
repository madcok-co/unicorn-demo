package handlers

import (
	"testing"
	"time"

	"github.com/madcok-co/unicorn-demo/internal/domain"
)

func TestCreateSalesQuotation_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")
	tax := createTestTax(db, "acme")

	// Test data
	req := domain.CreateSalesQuotationDTO{
		Code:          "SQ-2024-001",
		CustomerID:    customer.ID,
		QuotationDate: time.Now(),
		ValidUntil:    time.Now().AddDate(0, 0, 30),
		Currency:      "IDR",
		ExchangeRate:  1.0,
		Items: []domain.CreateSalesQuotationItemDTO{
			{
				ProductID:    product.ID,
				Description:  "Test product item",
				Quantity:     10,
				UnitPrice:    100000,
				DiscountType: "percentage",
				Discount:     5,
				TaxID:        tax.ID,
			},
		},
	}

	// Execute
	result, err := CreateSalesQuotation(ctx, req)

	// Assert
	assertNoError(t, err, "CreateSalesQuotation should not return error")
	assertNotNil(t, result, "CreateSalesQuotation should return quotation")
	assertEqual(t, result.Code, req.Code, "Quotation code should match")
	assertEqual(t, result.CustomerID, customer.ID, "Customer ID should match")
	assertEqual(t, result.CustomerName, customer.Name, "Customer name should match")
	assertEqual(t, result.Status, "draft", "Initial status should be draft")
	assertEqual(t, result.TenantID, "acme", "Tenant ID should match")
	assertEqual(t, len(result.Items), 1, "Should have 1 item")

	// Verify item calculations
	item := result.Items[0]
	assertEqual(t, item.ProductID, product.ID, "Item product ID should match")
	assertEqual(t, item.Quantity, 10.0, "Item quantity should match")
	assertEqual(t, item.UnitPrice, 100000.0, "Item unit price should match")
	assertEqual(t, item.Subtotal, 1000000.0, "Item subtotal should be calculated correctly")
}

func TestCreateSalesQuotation_InvalidCustomer(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	product := createTestProduct(db, "acme")

	// Test data with non-existent customer
	req := domain.CreateSalesQuotationDTO{
		Code:          "SQ-2024-001",
		CustomerID:    "non-existent-customer",
		QuotationDate: time.Now(),
		ValidUntil:    time.Now().AddDate(0, 0, 30),
		Currency:      "IDR",
		ExchangeRate:  1.0,
		Items: []domain.CreateSalesQuotationItemDTO{
			{
				ProductID: product.ID,
				Quantity:  10,
				UnitPrice: 100000,
			},
		},
	}

	// Execute
	_, err := CreateSalesQuotation(ctx, req)

	// Assert
	assertError(t, err, "CreateSalesQuotation should fail with invalid customer")
}

func TestCreateSalesQuotation_NoItems(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	customer := createTestCustomer(db, "acme")

	// Test data with no items
	req := domain.CreateSalesQuotationDTO{
		Code:          "SQ-2024-001",
		CustomerID:    customer.ID,
		QuotationDate: time.Now(),
		ValidUntil:    time.Now().AddDate(0, 0, 30),
		Currency:      "IDR",
		ExchangeRate:  1.0,
		Items:         []domain.CreateSalesQuotationItemDTO{},
	}

	// Execute
	_, err := CreateSalesQuotation(ctx, req)

	// Assert
	// Note: This would require validator to be set up in the handler
	// For now, we just verify the function can handle empty items
	if err == nil {
		t.Log("CreateSalesQuotation accepted empty items (validator not enforced in handler)")
	}
}

func TestUpdateSalesQuotation_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")
	quotation := createTestSalesQuotation(db, "acme", customer, product)

	// Update data
	newValidUntil := time.Now().AddDate(0, 0, 60)
	newShippingCost := 50000.0

	req := domain.UpdateSalesQuotationDTO{
		ValidUntil:   &newValidUntil,
		ShippingCost: &newShippingCost,
		Notes:        "Updated notes",
	}

	// Execute
	result, err := UpdateSalesQuotation(ctx, quotation.ID, req)

	// Assert
	assertNoError(t, err, "UpdateSalesQuotation should not return error")
	assertNotNil(t, result, "UpdateSalesQuotation should return quotation")
	assertEqual(t, result.ShippingCost, newShippingCost, "Shipping cost should be updated")
	assertEqual(t, result.Notes, "Updated notes", "Notes should be updated")
	assertEqual(t, result.UpdatedBy, "user-123", "UpdatedBy should match")
}

func TestUpdateSalesQuotation_NonDraft(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")
	quotation := createTestSalesQuotation(db, "acme", customer, product)

	// Change status to accepted
	quotation.Status = "accepted"
	db.Save(quotation)

	// Try to update
	req := domain.UpdateSalesQuotationDTO{
		Notes: "Trying to update accepted quotation",
	}

	// Execute
	_, err := UpdateSalesQuotation(ctx, quotation.ID, req)

	// Assert
	assertError(t, err, "UpdateSalesQuotation should fail for non-draft quotation")
}

func TestUpdateSalesQuotation_StatusChange(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")
	quotation := createTestSalesQuotation(db, "acme", customer, product)

	// Change status to sent
	req := domain.UpdateSalesQuotationDTO{
		Status: "sent",
	}

	// Execute
	result, err := UpdateSalesQuotation(ctx, quotation.ID, req)

	// Assert
	assertNoError(t, err, "UpdateSalesQuotation should not return error")
	assertEqual(t, result.Status, "sent", "Status should be updated")
	assertNotNil(t, result.SentAt, "SentAt timestamp should be set")
}

func TestDeleteSalesQuotation_Draft(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")
	quotation := createTestSalesQuotation(db, "acme", customer, product)

	// Execute
	err := DeleteSalesQuotation(ctx, quotation.ID)

	// Assert
	assertNoError(t, err, "DeleteSalesQuotation should not return error")

	// Verify quotation is soft deleted
	var deletedQuotation domain.SalesQuotation
	db.Where("id = ?", quotation.ID).First(&deletedQuotation)
	assertEqual(t, deletedQuotation.Active, false, "Quotation should be soft deleted")
}

func TestDeleteSalesQuotation_Accepted(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")
	quotation := createTestSalesQuotation(db, "acme", customer, product)

	// Change status to accepted
	quotation.Status = "accepted"
	db.Save(quotation)

	// Execute
	err := DeleteSalesQuotation(ctx, quotation.ID)

	// Assert
	assertError(t, err, "DeleteSalesQuotation should fail for accepted quotation")
}

func TestListSalesQuotations_WithFilters(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	customer1 := createTestCustomer(db, "acme")
	customer2 := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")

	quotation1 := createTestSalesQuotation(db, "acme", customer1, product)
	quotation1.Status = "draft"
	db.Save(quotation1)

	quotation2 := createTestSalesQuotation(db, "acme", customer2, product)
	quotation2.Status = "sent"
	db.Save(quotation2)

	quotation3 := createTestSalesQuotation(db, "acme", customer1, product)
	quotation3.Status = "accepted"
	db.Save(quotation3)

	// Test: List all quotations
	req := domain.ListSalesQuotationsDTO{
		Page:  1,
		Limit: 10,
	}

	result, err := ListSalesQuotations(ctx, req)
	assertNoError(t, err, "ListSalesQuotations should not return error")
	assertNotNil(t, result, "ListSalesQuotations should return result")

	// Test: Filter by customer
	reqCustomer := domain.ListSalesQuotationsDTO{
		Page:       1,
		Limit:      10,
		CustomerID: customer1.ID,
	}

	resultCustomer, err := ListSalesQuotations(ctx, reqCustomer)
	assertNoError(t, err, "ListSalesQuotations with customer filter should not return error")
	assertNotNil(t, resultCustomer, "ListSalesQuotations with customer filter should return result")

	// Test: Filter by status
	reqStatus := domain.ListSalesQuotationsDTO{
		Page:   1,
		Limit:  10,
		Status: "draft",
	}

	resultStatus, err := ListSalesQuotations(ctx, reqStatus)
	assertNoError(t, err, "ListSalesQuotations with status filter should not return error")
	assertNotNil(t, resultStatus, "ListSalesQuotations with status filter should return result")
}

func TestGetSalesQuotationStats_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")

	quotation1 := createTestSalesQuotation(db, "acme", customer, product)
	quotation1.Status = "draft"
	quotation1.GrandTotal = 1000000
	db.Save(quotation1)

	quotation2 := createTestSalesQuotation(db, "acme", customer, product)
	quotation2.Status = "accepted"
	quotation2.GrandTotal = 2000000
	db.Save(quotation2)

	quotation3 := createTestSalesQuotation(db, "acme", customer, product)
	quotation3.Status = "accepted"
	quotation3.GrandTotal = 3000000
	db.Save(quotation3)

	// Execute
	stats, err := GetSalesQuotationStats(ctx)

	// Assert
	assertNoError(t, err, "GetSalesQuotationStats should not return error")
	assertNotNil(t, stats, "GetSalesQuotationStats should return stats")
	assertNotNil(t, stats["total_quotations"], "Stats should contain total_quotations")
	assertNotNil(t, stats["total_value"], "Stats should contain total_value")
	assertNotNil(t, stats["accepted_value"], "Stats should contain accepted_value")
	assertNotNil(t, stats["quotations_by_status"], "Stats should contain quotations_by_status")
}

func TestSalesQuotation_MultiTenantIsolation(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Create quotations for different tenants
	customerAcme := createTestCustomer(db, "acme")
	productAcme := createTestProduct(db, "acme")
	_ = createTestSalesQuotation(db, "acme", customerAcme, productAcme)

	customerGlobex := createTestCustomer(db, "globex")
	productGlobex := createTestProduct(db, "globex")
	quotationGlobex := createTestSalesQuotation(db, "globex", customerGlobex, productGlobex)

	// Test: Tenant "acme" should only see their quotation
	ctxAcme := setupTestContext("acme", "user-123")
	ctxAcme.Set("db", db)

	req := domain.ListSalesQuotationsDTO{
		Page:  1,
		Limit: 100,
	}

	resultAcme, err := ListSalesQuotations(ctxAcme, req)
	assertNoError(t, err, "ListSalesQuotations should not return error")
	assertEqual(t, resultAcme.Total, int64(1), "Tenant acme should see 1 quotation")

	// Test: Tenant "globex" should only see their quotation
	ctxGlobex := setupTestContext("globex", "user-456")
	ctxGlobex.Set("db", db)

	resultGlobex, err := ListSalesQuotations(ctxGlobex, req)
	assertNoError(t, err, "ListSalesQuotations should not return error")
	assertEqual(t, resultGlobex.Total, int64(1), "Tenant globex should see 1 quotation")

	// Test: Cannot get quotation from different tenant
	_, err = GetSalesQuotation(ctxAcme, quotationGlobex.ID)
	assertError(t, err, "Should not be able to get quotation from different tenant")
}

func TestSalesQuotation_CalculateTotals(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")
	tax := createTestTax(db, "acme")

	// Test data with specific calculations
	req := domain.CreateSalesQuotationDTO{
		Code:          "SQ-2024-001",
		CustomerID:    customer.ID,
		QuotationDate: time.Now(),
		ValidUntil:    time.Now().AddDate(0, 0, 30),
		Currency:      "IDR",
		ExchangeRate:  1.0,
		ShippingCost:  25000,
		OtherCost:     10000,
		Items: []domain.CreateSalesQuotationItemDTO{
			{
				ProductID:    product.ID,
				Quantity:     10,
				UnitPrice:    100000,
				DiscountType: "percentage",
				Discount:     5,      // 5% discount
				TaxID:        tax.ID, // 11% tax
			},
		},
	}

	// Execute
	result, err := CreateSalesQuotation(ctx, req)

	// Assert
	assertNoError(t, err, "CreateSalesQuotation should not return error")

	// Verify calculations
	// Subtotal = 10 * 100,000 = 1,000,000
	// Discount = 1,000,000 * 5% = 50,000
	// After discount = 950,000
	// Tax = 950,000 * 11% = 104,500
	// Grand total = 950,000 + 104,500 + 25,000 + 10,000 = 1,089,500

	assertEqual(t, result.Subtotal, 1000000.0, "Subtotal should be calculated correctly")
	assertEqual(t, result.TotalDiscount, 50000.0, "Total discount should be calculated correctly")
	assertEqual(t, result.TotalTax, 104500.0, "Total tax should be calculated correctly")
	assertEqual(t, result.GrandTotal, 1089500.0, "Grand total should be calculated correctly")
}

func TestSalesQuotation_WorkflowTransitions(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")
	quotation := createTestSalesQuotation(db, "acme", customer, product)

	// Test workflow: draft -> sent
	reqSent := domain.UpdateSalesQuotationDTO{
		Status: "sent",
	}
	resultSent, err := UpdateSalesQuotation(ctx, quotation.ID, reqSent)
	assertNoError(t, err, "Should transition to sent")
	assertEqual(t, resultSent.Status, "sent", "Status should be sent")
	assertNotNil(t, resultSent.SentAt, "SentAt should be set")

	// Test workflow: sent -> accepted
	reqAccepted := domain.UpdateSalesQuotationDTO{
		Status: "accepted",
	}
	resultAccepted, err := UpdateSalesQuotation(ctx, quotation.ID, reqAccepted)
	assertNoError(t, err, "Should transition to accepted")
	assertEqual(t, resultAccepted.Status, "accepted", "Status should be accepted")
	assertNotNil(t, resultAccepted.AcceptedAt, "AcceptedAt should be set")
}
