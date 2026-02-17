package handlers

import (
	"testing"

	"github.com/madcok-co/unicorn-demo/internal/domain"
)

func TestCreateCustomer_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Test data
	req := domain.CreateCustomerDTO{
		Code:             "CUST001",
		Name:             "Test Customer Inc",
		Type:             "company",
		Email:            "contact@testcustomer.com",
		Phone:            "123-456-7890",
		Mobile:           "098-765-4321",
		TaxID:            "TAX123456789",
		BillingAddress:   "123 Main St",
		BillingCity:      "Jakarta",
		BillingState:     "DKI Jakarta",
		BillingZip:       "12345",
		BillingCountry:   "Indonesia",
		ShippingAddress:  "456 Ship St",
		ShippingCity:     "Jakarta",
		ShippingState:    "DKI Jakarta",
		ShippingZip:      "54321",
		ShippingCountry:  "Indonesia",
		CreditLimit:      50000000,
		PaymentTermDays:  30,
		Currency:         "IDR",
		ContactPerson:    "John Doe",
		ContactPersonJob: "Purchasing Manager",
		Notes:            "Important customer",
		Tags:             []string{"vip", "corporate"},
	}

	// Execute
	result, err := CreateCustomer(ctx, req)

	// Assert
	assertNoError(t, err, "CreateCustomer should not return error")
	assertNotNil(t, result, "CreateCustomer should return customer")
	assertEqual(t, result.Code, req.Code, "Customer code should match")
	assertEqual(t, result.Name, req.Name, "Customer name should match")
	assertEqual(t, result.Type, req.Type, "Customer type should match")
	assertEqual(t, result.Email, req.Email, "Customer email should match")
	assertEqual(t, result.TenantID, "acme", "Customer tenant ID should match")
	assertEqual(t, result.Active, true, "Customer should be active")
	assertEqual(t, result.CreatedBy, "user-123", "CreatedBy should match")
}

func TestCreateCustomer_DuplicateCode(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create first customer
	_ = createTestCustomer(db, "acme")

	// Try to create duplicate
	req := domain.CreateCustomerDTO{
		Code: "CUST-" + createTestCustomer(db, "acme").Code[5:],
		Name: "Duplicate Customer",
		Type: "individual",
	}

	// Execute
	_, err := CreateCustomer(ctx, req)

	// Assert - should fail due to duplicate code
	assertError(t, err, "CreateCustomer should fail with duplicate code")
}

func TestCreateCustomer_InvalidTenant(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Context without tenant_id
	ctx := setupTestContext("", "user-123")
	ctx.Set("db", db)
	ctx.Set("tenant_id", nil)

	req := domain.CreateCustomerDTO{
		Code: "CUST001",
		Name: "Test Customer",
		Type: "company",
	}

	// Execute
	_, err := CreateCustomer(ctx, req)

	// Assert
	assertError(t, err, "CreateCustomer should fail without tenant")
}

func TestGetCustomer_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test customer
	customer := createTestCustomer(db, "acme")

	// Execute
	result, err := GetCustomer(ctx, customer.ID)

	// Assert
	assertNoError(t, err, "GetCustomer should not return error")
	assertNotNil(t, result, "GetCustomer should return customer")
	assertEqual(t, result.ID, customer.ID, "Customer ID should match")
	assertEqual(t, result.Code, customer.Code, "Customer code should match")
	assertEqual(t, result.Name, customer.Name, "Customer name should match")
}

func TestGetCustomer_NotFound(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Execute with non-existent ID
	_, err := GetCustomer(ctx, "non-existent-id")

	// Assert
	assertError(t, err, "GetCustomer should fail with non-existent ID")
}

func TestGetCustomer_DifferentTenant(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Create customer for tenant "acme"
	customer := createTestCustomer(db, "acme")

	// Try to get customer with different tenant context
	ctx := setupTestContext("other-tenant", "user-123")
	ctx.Set("db", db)

	// Execute
	_, err := GetCustomer(ctx, customer.ID)

	// Assert - should fail due to tenant isolation
	assertError(t, err, "GetCustomer should fail with different tenant")
}

func TestListCustomers_WithFilters(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test customers
	customer1 := createTestCustomer(db, "acme")
	customer1.Type = "company"
	db.Save(customer1)
	_ = customer1 // Mark as used

	customer2 := createTestCustomer(db, "acme")
	customer2.Type = "individual"
	db.Save(customer2)

	customer3 := createTestCustomer(db, "acme")
	customer3.Type = "company"
	db.Save(customer3)

	// Test: List all customers
	req := domain.ListCustomersDTO{
		Page:  1,
		Limit: 10,
	}

	result, err := ListCustomers(ctx, req)
	assertNoError(t, err, "ListCustomers should not return error")
	assertNotNil(t, result, "ListCustomers should return result")

	// Test: Filter by type
	reqFiltered := domain.ListCustomersDTO{
		Page:  1,
		Limit: 10,
		Type:  "company",
	}

	resultFiltered, err := ListCustomers(ctx, reqFiltered)
	assertNoError(t, err, "ListCustomers with filter should not return error")
	assertNotNil(t, resultFiltered, "ListCustomers with filter should return result")

	// Test: Search by name
	reqSearch := domain.ListCustomersDTO{
		Page:   1,
		Limit:  10,
		Search: customer1.Name,
	}

	resultSearch, err := ListCustomers(ctx, reqSearch)
	assertNoError(t, err, "ListCustomers with search should not return error")
	assertNotNil(t, resultSearch, "ListCustomers with search should return result")
}

func TestUpdateCustomer_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test customer
	customer := createTestCustomer(db, "acme")

	// Update data
	newName := "Updated Customer Name"
	newEmail := "updated@test.com"
	newCreditLimit := 200000.0

	req := domain.UpdateCustomerDTO{
		Name:        newName,
		Email:       newEmail,
		CreditLimit: &newCreditLimit,
	}

	// Execute
	result, err := UpdateCustomer(ctx, customer.ID, req)

	// Assert
	assertNoError(t, err, "UpdateCustomer should not return error")
	assertNotNil(t, result, "UpdateCustomer should return customer")
	assertEqual(t, result.Name, newName, "Customer name should be updated")
	assertEqual(t, result.Email, newEmail, "Customer email should be updated")
	assertEqual(t, result.CreditLimit, newCreditLimit, "Customer credit limit should be updated")
	assertEqual(t, result.UpdatedBy, "user-123", "UpdatedBy should match")
}

func TestUpdateCustomer_NotFound(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	req := domain.UpdateCustomerDTO{
		Name: "Updated Name",
	}

	// Execute with non-existent ID
	_, err := UpdateCustomer(ctx, "non-existent-id", req)

	// Assert
	assertError(t, err, "UpdateCustomer should fail with non-existent ID")
}

func TestDeleteCustomer_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test customer
	customer := createTestCustomer(db, "acme")

	// Execute
	err := DeleteCustomer(ctx, customer.ID)

	// Assert
	assertNoError(t, err, "DeleteCustomer should not return error")

	// Verify customer is soft deleted (active = false)
	var deletedCustomer domain.Customer
	db.Where("id = ?", customer.ID).First(&deletedCustomer)
	assertEqual(t, deletedCustomer.Active, false, "Customer should be soft deleted")
}

func TestDeleteCustomer_NotFound(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Execute with non-existent ID
	err := DeleteCustomer(ctx, "non-existent-id")

	// Assert
	assertError(t, err, "DeleteCustomer should fail with non-existent ID")
}

func TestGetCustomerStats_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create test customers
	customer1 := createTestCustomer(db, "acme")
	customer1.Type = "company"
	db.Save(customer1)

	customer2 := createTestCustomer(db, "acme")
	customer2.Type = "individual"
	db.Save(customer2)

	customer3 := createTestCustomer(db, "acme")
	customer3.Type = "company"
	db.Save(customer3)

	// Execute
	stats, err := GetCustomerStats(ctx)

	// Assert
	assertNoError(t, err, "GetCustomerStats should not return error")
	assertNotNil(t, stats, "GetCustomerStats should return stats")
	assertNotNil(t, stats["total_customers"], "Stats should contain total_customers")
	assertNotNil(t, stats["customers_by_type"], "Stats should contain customers_by_type")
}

func TestCustomer_MultiTenantIsolation(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Create customers for different tenants
	_ = createTestCustomer(db, "acme")
	customer2 := createTestCustomer(db, "globex")

	// Test: Tenant "acme" should only see their customer
	ctxAcme := setupTestContext("acme", "user-123")
	ctxAcme.Set("db", db)

	req := domain.ListCustomersDTO{
		Page:  1,
		Limit: 100,
	}

	resultAcme, err := ListCustomers(ctxAcme, req)
	assertNoError(t, err, "ListCustomers should not return error")
	assertEqual(t, resultAcme.Total, int64(1), "Tenant acme should see 1 customer")

	// Test: Tenant "globex" should only see their customer
	ctxGlobex := setupTestContext("globex", "user-456")
	ctxGlobex.Set("db", db)

	resultGlobex, err := ListCustomers(ctxGlobex, req)
	assertNoError(t, err, "ListCustomers should not return error")
	assertEqual(t, resultGlobex.Total, int64(1), "Tenant globex should see 1 customer")

	// Test: Cannot get customer from different tenant
	_, err = GetCustomer(ctxAcme, customer2.ID)
	assertError(t, err, "Should not be able to get customer from different tenant")
}

func TestCustomer_BusinessRuleValidations(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	tests := []struct {
		name        string
		req         domain.CreateCustomerDTO
		shouldError bool
		description string
	}{
		{
			name: "Valid company customer",
			req: domain.CreateCustomerDTO{
				Code: "CUST001",
				Name: "Test Company",
				Type: "company",
			},
			shouldError: false,
			description: "Should create valid company customer",
		},
		{
			name: "Valid individual customer",
			req: domain.CreateCustomerDTO{
				Code: "CUST002",
				Name: "John Doe",
				Type: "individual",
			},
			shouldError: false,
			description: "Should create valid individual customer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CreateCustomer(ctx, tt.req)

			if tt.shouldError {
				assertError(t, err, tt.description)
			} else {
				assertNoError(t, err, tt.description)
				assertNotNil(t, result, "Should return customer")
			}
		})
	}
}
