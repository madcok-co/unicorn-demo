# ERP Module Tests

Comprehensive unit and integration tests for all ERP modules following Go testing best practices.

## Test Structure

```
internal/handlers/
├── test_helpers.go              # Test utilities and setup functions
├── erp_customers_test.go        # Customer management tests
├── sales_quotations_test.go     # Sales quotation tests  
├── purchase_requests_test.go    # Purchase request tests
└── integration_test.go          # End-to-end workflow tests
```

## Running Tests

### Run All Tests
```bash
cd internal/handlers
go test -v
```

### Run Specific Test File
```bash
go test -v -run TestCreateCustomer
go test -v -run TestSalesQuotation
go test -v -run TestPurchaseRequest
```

### Run Integration Tests Only
```bash
go test -v -run TestIntegration
go test -v -run TestWorkflow
```

### Run with Coverage
```bash
go test -v -cover
go test -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run with Race Detector
```bash
go test -v -race
```

## Test Helpers

The `test_helpers.go` file provides utilities for test setup:

### Database Setup
- `setupTestDB(t *testing.T) *gorm.DB` - Creates in-memory SQLite database with all tables
- `teardownTestDB(db *gorm.DB)` - Cleanup database connection

### Context Setup
- `setupTestContext(tenantID, userID string) *context.Context` - Create context with tenant and user

### Test Data Creation
- `createTestCustomer(db, tenantID)` - Create test customer
- `createTestVendor(db, tenantID)` - Create test vendor
- `createTestProduct(db, tenantID)` - Create test product
- `createTestWarehouse(db, tenantID)` - Create test warehouse
- `createTestTax(db, tenantID)` - Create test tax configuration
- `createTestCurrency(db, tenantID)` - Create test currency
- `createTestPaymentTerm(db, tenantID)` - Create test payment term
- `createTestChartOfAccount(db, tenantID)` - Create test chart of account
- `createTestUser(db, tenantID)` - Create test user
- `createTestSalesQuotation(db, tenantID, customer, product)` - Create test sales quotation
- `createTestSalesOrder(db, tenantID, customer, product)` - Create test sales order
- `createTestPurchaseRequest(db, tenantID, user, product)` - Create test purchase request
- `createTestPurchaseOrder(db, tenantID, vendor, product)` - Create test purchase order

### Assertion Helpers
- `assertNoError(t, err, msg)` - Check for no error
- `assertError(t, err, msg)` - Check for error
- `assertEqual(t, got, want, msg)` - Check equality
- `assertNotNil(t, val, msg)` - Check non-nil value

## Test Coverage

### ERP Master Data Tests (`erp_customers_test.go`)

**Customer Management:**
- ✅ Create customer - success case
- ✅ Create customer - duplicate code (should fail)
- ✅ Create customer - invalid tenant (should fail)
- ✅ Get customer - success
- ✅ Get customer - not found (should fail)
- ✅ Get customer - different tenant (should fail)
- ✅ List customers - with pagination and filters
- ✅ Update customer - success
- ✅ Update customer - not found (should fail)
- ✅ Delete customer - soft delete success
- ✅ Delete customer - not found (should fail)
- ✅ Get customer statistics
- ✅ Multi-tenant isolation
- ✅ Business rule validations

**Similar test coverage exists for:**
- Vendors
- Products
- Warehouses
- Chart of Accounts
- Taxes
- Payment Terms
- Currencies

### Sales Module Tests (`sales_quotations_test.go`)

**Sales Quotations:**
- ✅ Create quotation - success with items
- ✅ Create quotation - invalid customer (should fail)
- ✅ Create quotation - no items (validation)
- ✅ Update quotation - success
- ✅ Update quotation - non-draft status (should fail)
- ✅ Update quotation - status transitions
- ✅ Delete quotation - draft (success)
- ✅ Delete quotation - accepted (should fail)
- ✅ List quotations - with filters
- ✅ Get quotation statistics
- ✅ Multi-tenant isolation
- ✅ Calculate totals correctly
- ✅ Workflow state transitions

**Similar test coverage for:**
- Sales Orders
- Delivery Orders
- Sales Invoices

### Purchase Module Tests (`purchase_requests_test.go`)

**Purchase Requests:**
- ✅ Create request - success with items
- ✅ Create request - invalid requester (should fail)
- ✅ Update request - approval workflow
- ✅ Update request - rejection workflow
- ✅ Delete request - draft (success)
- ✅ Delete request - approved (should fail)
- ✅ List requests - with filters
- ✅ Get request statistics
- ✅ Multi-tenant isolation
- ✅ Priority handling
- ✅ Status restrictions
- ✅ Item quantity tracking

**Similar test coverage for:**
- Purchase Orders
- Goods Receipts
- Purchase Invoices

### Integration Tests (`integration_test.go`)

**Complete Workflows:**
- ✅ Sales workflow: Quotation → Order → Delivery → Invoice
- ✅ Purchase workflow: Request → Order → Receipt → Invoice
- ✅ Quantity tracking: Delivery order updates delivered_qty
- ✅ Quantity tracking: Goods receipt updates received_qty
- ✅ Payment status: Automatic calculation (unpaid → partial → paid)
- ✅ Cross-module integration: Sales order triggers purchase request

## Test Patterns

### Standard Test Structure
```go
func TestFeature_Scenario(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer teardownTestDB(db)
    
    ctx := setupTestContext("tenant-id", "user-id")
    ctx.Set("db", db)
    
    // Create test data
    customer := createTestCustomer(db, "tenant-id")
    
    // Test data
    req := domain.CreateCustomerDTO{
        Code: "CUST001",
        Name: "Test Customer",
        Type: "company",
    }
    
    // Execute
    result, err := CreateCustomer(ctx, req)
    
    // Assert
    assertNoError(t, err, "CreateCustomer should not return error")
    assertEqual(t, result.Code, req.Code, "Code should match")
}
```

### Table-Driven Tests
```go
func TestFeature_MultipleScenarios(t *testing.T) {
    tests := []struct {
        name        string
        input       domain.CreateDTO
        shouldError bool
        expected    string
    }{
        {
            name:        "Valid input",
            input:       validInput,
            shouldError: false,
            expected:    "success",
        },
        {
            name:        "Invalid input",
            input:       invalidInput,
            shouldError: true,
            expected:    "",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := CreateFeature(ctx, tt.input)
            
            if tt.shouldError {
                assertError(t, err, tt.name)
            } else {
                assertNoError(t, err, tt.name)
                assertEqual(t, result.Field, tt.expected, "Field should match")
            }
        })
    }
}
```

## Key Test Scenarios

### 1. Multi-Tenant Isolation
Tests verify that:
- Each tenant can only access their own data
- Queries are scoped by tenant_id
- Cross-tenant access is prevented

### 2. Business Rule Validation
Tests verify that:
- Required fields are validated
- Status transitions follow business rules
- Delete operations respect data states

### 3. Workflow State Management
Tests verify that:
- Status transitions are tracked with timestamps
- Invalid transitions are prevented
- Workflow states are persisted correctly

### 4. Quantity Tracking
Tests verify that:
- delivered_qty updates when delivery is completed
- received_qty updates when goods are received
- invoiced_qty tracks billed quantities
- Partial deliveries/receipts are supported

### 5. Financial Calculations
Tests verify that:
- Subtotals are calculated correctly
- Discounts (percentage and fixed) work properly
- Taxes are calculated accurately
- Grand totals include all components

### 6. Payment Status
Tests verify that:
- Status changes from unpaid → partial → paid
- amount_paid and amount_due are calculated
- Timestamps are set correctly

## Best Practices

### 1. Test Isolation
- Each test creates its own in-memory database
- No shared state between tests
- Cleanup is performed after each test

### 2. Meaningful Test Names
- Format: `TestFunction_Scenario`
- Examples: `TestCreateCustomer_Success`, `TestDeleteOrder_Approved`

### 3. Clear Assertions
- Use helper functions for consistent error messages
- Provide descriptive assertion messages
- Test both success and failure cases

### 4. Test Data Management
- Use helper functions to create test data
- Keep test data realistic and meaningful
- Include edge cases and boundary conditions

### 5. Documentation
- Add comments explaining complex test scenarios
- Document assumptions and prerequisites
- Note any known limitations

## Adding New Tests

When adding new handlers, create corresponding tests:

1. Create test file: `{module}_test.go`
2. Implement success case tests
3. Implement error case tests
4. Add multi-tenant isolation tests
5. Add business rule validation tests
6. Update this README with coverage

### Example Template
```go
package handlers

import (
    "testing"
    "github.com/madcok-co/unicorn-demo/internal/domain"
)

func TestCreateFeature_Success(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer teardownTestDB(db)
    ctx := setupTestContext("acme", "user-123")
    ctx.Set("db", db)
    
    // Create dependencies
    dependency := createTestDependency(db, "acme")
    
    // Test data
    req := domain.CreateFeatureDTO{
        // ... fields
    }
    
    // Execute
    result, err := CreateFeature(ctx, req)
    
    // Assert
    assertNoError(t, err, "CreateFeature should succeed")
    assertNotNil(t, result, "Result should not be nil")
    assertEqual(t, result.Field, req.Field, "Field should match")
}

func TestCreateFeature_ValidationError(t *testing.T) {
    // Test validation failures
}

func TestCreateFeature_MultiTenant(t *testing.T) {
    // Test tenant isolation
}
```

## Continuous Integration

### GitHub Actions Example
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: |
          cd internal/handlers
          go test -v -race -coverprofile=coverage.out
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

## Troubleshooting

### Test Database Issues
If you encounter database errors:
- Ensure SQLite driver is installed: `go get gorm.io/driver/sqlite`
- Check that AutoMigrate includes all models
- Verify foreign key constraints

### Context Issues
If handlers can't access DB:
- Ensure `ctx.Set("db", db)` is called
- Check that getDB helper returns non-nil
- Verify tenant_id and user_id are set

### Assertion Failures
If assertions fail unexpectedly:
- Use `-v` flag to see detailed output
- Check actual vs expected values
- Verify test data setup

## Contributing

When contributing tests:
1. Follow existing test patterns
2. Ensure tests are isolated and repeatable
3. Add both positive and negative test cases
4. Update this README with new coverage
5. Maintain high code coverage (>80%)

## License

Same as main project license.
