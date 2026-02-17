package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	ucontext "github.com/madcok-co/unicorn/core/pkg/context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB creates an in-memory SQLite database with all tables migrated
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		// Master Data
		&domain.Customer{},
		&domain.Vendor{},
		&domain.Product{},
		&domain.Warehouse{},
		&domain.ChartOfAccount{},
		&domain.Tax{},
		&domain.PaymentTerm{},
		&domain.Currency{},
		// Sales Module
		&domain.SalesQuotation{},
		&domain.SalesQuotationItem{},
		&domain.SalesOrder{},
		&domain.SalesOrderItem{},
		&domain.DeliveryOrder{},
		&domain.DeliveryOrderItem{},
		&domain.SalesInvoice{},
		&domain.SalesInvoiceItem{},
		// Purchase Module
		&domain.PurchaseRequest{},
		&domain.PurchaseRequestItem{},
		&domain.PurchaseOrder{},
		&domain.PurchaseOrderItem{},
		&domain.GoodsReceipt{},
		&domain.GoodsReceiptItem{},
		&domain.PurchaseInvoice{},
		&domain.PurchaseInvoiceItem{},
		// User & Auth
		&domain.User{},
		&domain.AuditLog{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

// setupTestContext creates a context with tenant and user IDs
func setupTestContext(tenantID, userID string) *ucontext.Context {
	ctx := ucontext.New(context.Background())
	ctx.Set("tenant_id", tenantID)
	ctx.Set("user_id", userID)
	return ctx
}

// createTestCustomer creates a test customer in the database
func createTestCustomer(db *gorm.DB, tenantID string) *domain.Customer {
	customer := &domain.Customer{
		ID:               uuid.New().String(),
		Code:             "CUST-" + uuid.New().String()[:8],
		Name:             "Test Customer",
		TenantID:         tenantID,
		Type:             "company",
		Email:            "customer@test.com",
		Phone:            "123-456-7890",
		Mobile:           "098-765-4321",
		TaxID:            "TAX123456",
		BillingAddress:   "123 Test St",
		BillingCity:      "Test City",
		BillingState:     "Test State",
		BillingZip:       "12345",
		BillingCountry:   "Indonesia",
		ShippingAddress:  "456 Ship St",
		ShippingCity:     "Ship City",
		ShippingState:    "Ship State",
		ShippingZip:      "54321",
		ShippingCountry:  "Indonesia",
		CreditLimit:      100000,
		PaymentTermDays:  30,
		Currency:         "IDR",
		ContactPerson:    "John Doe",
		ContactPersonJob: "Manager",
		Active:           true,
		CreatedBy:        "system",
		UpdatedBy:        "system",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	db.Create(customer)
	return customer
}

// createTestVendor creates a test vendor in the database
func createTestVendor(db *gorm.DB, tenantID string) *domain.Vendor {
	vendor := &domain.Vendor{
		ID:               uuid.New().String(),
		Code:             "VEND-" + uuid.New().String()[:8],
		Name:             "Test Vendor",
		TenantID:         tenantID,
		Type:             "company",
		Email:            "vendor@test.com",
		Phone:            "111-222-3333",
		Mobile:           "444-555-6666",
		TaxID:            "VTAX123456",
		Address:          "789 Vendor St",
		City:             "Vendor City",
		State:            "Vendor State",
		Zip:              "99999",
		Country:          "Indonesia",
		PaymentTermDays:  30,
		Currency:         "IDR",
		BankName:         "Test Bank",
		BankAccount:      "1234567890",
		BankAccountName:  "Test Vendor Inc",
		ContactPerson:    "Jane Smith",
		ContactPersonJob: "Sales Manager",
		Rating:           4,
		Category:         "raw-material",
		Active:           true,
		CreatedBy:        "system",
		UpdatedBy:        "system",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	db.Create(vendor)
	return vendor
}

// createTestProduct creates a test product in the database
func createTestProduct(db *gorm.DB, tenantID string) *domain.Product {
	product := &domain.Product{
		ID:             uuid.New().String(),
		Code:           "PROD-" + uuid.New().String()[:8],
		Name:           "Test Product",
		TenantID:       tenantID,
		Type:           "product",
		Category:       "electronics",
		UOM:            "pcs",
		Barcode:        "123456789012",
		SKU:            "SKU-TEST-001",
		Active:         true,
		CanBeSold:      true,
		CanBePurchased: true,
		TrackInventory: true,
		MinStock:       10,
		MaxStock:       100,
		ReorderLevel:   20,
		SalePrice:      100000,
		PurchasePrice:  80000,
		Cost:           75000,
		Currency:       "IDR",
		TaxCategory:    "taxable",
		TaxRate:        11,
		Weight:         1.5,
		Volume:         0.5,
		Length:         10,
		Width:          10,
		Height:         10,
		WarrantyDays:   365,
		Description:    "Test product description",
		InternalNotes:  "Internal notes for testing",
		CreatedBy:      "system",
		UpdatedBy:      "system",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	db.Create(product)
	return product
}

// createTestWarehouse creates a test warehouse in the database
func createTestWarehouse(db *gorm.DB, tenantID string) *domain.Warehouse {
	warehouse := &domain.Warehouse{
		ID:                 uuid.New().String(),
		Code:               "WH-" + uuid.New().String()[:8],
		Name:               "Test Warehouse",
		TenantID:           tenantID,
		Type:               "physical",
		Active:             true,
		Address:            "100 Warehouse Ave",
		City:               "Warehouse City",
		State:              "Warehouse State",
		Zip:                "11111",
		Country:            "Indonesia",
		Phone:              "777-888-9999",
		Email:              "warehouse@test.com",
		ManagerID:          "manager-123",
		ManagerName:        "Warehouse Manager",
		TotalArea:          1000,
		UsedArea:           500,
		TotalCapacity:      10000,
		UsedCapacity:       5000,
		AllowNegativeStock: false,
		IsDefault:          true,
		Description:        "Main warehouse for testing",
		CreatedBy:          "system",
		UpdatedBy:          "system",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	db.Create(warehouse)
	return warehouse
}

// createTestTax creates a test tax in the database
func createTestTax(db *gorm.DB, tenantID string) *domain.Tax {
	tax := &domain.Tax{
		ID:              uuid.New().String(),
		Code:            "PPN",
		Name:            "PPN 11%",
		TenantID:        tenantID,
		Type:            "both",
		Rate:            11.0,
		Active:          true,
		TaxAccountID:    "tax-account-123",
		IsDefault:       true,
		IncludedInPrice: false,
		Description:     "Value Added Tax 11%",
		CreatedBy:       "system",
		UpdatedBy:       "system",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	db.Create(tax)
	return tax
}

// createTestCurrency creates a test currency in the database
func createTestCurrency(db *gorm.DB, tenantID string) *domain.Currency {
	currency := &domain.Currency{
		ID:            uuid.New().String(),
		Code:          "IDR",
		Name:          "Indonesian Rupiah",
		Symbol:        "Rp",
		TenantID:      tenantID,
		Active:        true,
		ExchangeRate:  1.0,
		IsDefault:     true,
		DecimalPlaces: 0,
		CreatedBy:     "system",
		UpdatedBy:     "system",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	db.Create(currency)
	return currency
}

// createTestPaymentTerm creates a test payment term in the database
func createTestPaymentTerm(db *gorm.DB, tenantID string) *domain.PaymentTerm {
	paymentTerm := &domain.PaymentTerm{
		ID:              uuid.New().String(),
		Code:            "NET30",
		Name:            "Net 30 Days",
		TenantID:        tenantID,
		Days:            30,
		Type:            "fixed",
		Active:          true,
		DiscountDays:    10,
		DiscountPercent: 2.0,
		IsDefault:       true,
		Description:     "Payment due in 30 days, 2% discount if paid within 10 days",
		CreatedBy:       "system",
		UpdatedBy:       "system",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	db.Create(paymentTerm)
	return paymentTerm
}

// createTestChartOfAccount creates a test chart of account in the database
func createTestChartOfAccount(db *gorm.DB, tenantID string) *domain.ChartOfAccount {
	coa := &domain.ChartOfAccount{
		ID:                  uuid.New().String(),
		Code:                "1-1001",
		Name:                "Cash in Hand",
		TenantID:            tenantID,
		Type:                "asset",
		Category:            "current_asset",
		Level:               2,
		IsGroup:             false,
		Active:              true,
		Currency:            "IDR",
		CurrentBalance:      0,
		DebitBalance:        0,
		CreditBalance:       0,
		AllowReconciliation: true,
		RequireTaxReporting: false,
		Description:         "Cash and cash equivalents",
		CreatedBy:           "system",
		UpdatedBy:           "system",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	db.Create(coa)
	return coa
}

// createTestUser creates a test user in the database
func createTestUser(db *gorm.DB, tenantID string) *domain.User {
	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     "test@test.com",
		Name:      "Test User",
		TenantID:  tenantID,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(user)
	return user
}

// createTestSalesQuotation creates a test sales quotation with items
func createTestSalesQuotation(db *gorm.DB, tenantID string, customer *domain.Customer, product *domain.Product) *domain.SalesQuotation {
	quotation := &domain.SalesQuotation{
		ID:              uuid.New().String(),
		Code:            "SQ-" + uuid.New().String()[:8],
		TenantID:        tenantID,
		CustomerID:      customer.ID,
		CustomerName:    customer.Name,
		CustomerType:    customer.Type,
		QuotationDate:   time.Now(),
		ValidUntil:      time.Now().AddDate(0, 0, 30),
		Status:          "draft",
		Currency:        "IDR",
		ExchangeRate:    1.0,
		Subtotal:        1000000,
		TotalDiscount:   50000,
		TotalTax:        105000,
		ShippingCost:    25000,
		OtherCost:       10000,
		GrandTotal:      1090000,
		PaymentTermDays: 30,
		Active:          true,
		CreatedBy:       "system",
		UpdatedBy:       "system",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	quotationItem := &domain.SalesQuotationItem{
		ID:           uuid.New().String(),
		QuotationID:  quotation.ID,
		TenantID:     tenantID,
		LineNumber:   1,
		ProductID:    product.ID,
		ProductCode:  product.Code,
		ProductName:  product.Name,
		ProductType:  product.Type,
		Quantity:     10,
		UOM:          product.UOM,
		UnitPrice:    100000,
		DiscountType: "percentage",
		Discount:     5,
		TaxRate:      11,
		TaxAmount:    105000,
		Subtotal:     1000000,
		Total:        1055000,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	quotation.Items = []domain.SalesQuotationItem{*quotationItem}
	db.Create(quotation)
	db.Create(quotationItem)

	return quotation
}

// createTestSalesOrder creates a test sales order with items
func createTestSalesOrder(db *gorm.DB, tenantID string, customer *domain.Customer, product *domain.Product) *domain.SalesOrder {
	order := &domain.SalesOrder{
		ID:                uuid.New().String(),
		Code:              "SO-" + uuid.New().String()[:8],
		TenantID:          tenantID,
		CustomerID:        customer.ID,
		CustomerName:      customer.Name,
		CustomerType:      customer.Type,
		OrderDate:         time.Now(),
		ExpectedDate:      time.Now().AddDate(0, 0, 7),
		Status:            "draft",
		Priority:          "normal",
		Currency:          "IDR",
		ExchangeRate:      1.0,
		Subtotal:          1000000,
		TotalDiscount:     50000,
		TotalTax:          105000,
		ShippingCost:      25000,
		OtherCost:         10000,
		GrandTotal:        1090000,
		PaymentTermDays:   30,
		PaymentStatus:     "unpaid",
		FulfillmentStatus: "pending",
		DeliveredQty:      0,
		InvoicedQty:       0,
		Active:            true,
		CreatedBy:         "system",
		UpdatedBy:         "system",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	orderItem := &domain.SalesOrderItem{
		ID:           uuid.New().String(),
		OrderID:      order.ID,
		TenantID:     tenantID,
		LineNumber:   1,
		ProductID:    product.ID,
		ProductCode:  product.Code,
		ProductName:  product.Name,
		ProductType:  product.Type,
		Quantity:     10,
		DeliveredQty: 0,
		InvoicedQty:  0,
		UOM:          product.UOM,
		UnitPrice:    100000,
		DiscountType: "percentage",
		Discount:     5,
		TaxRate:      11,
		TaxAmount:    105000,
		Subtotal:     1000000,
		Total:        1055000,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	order.Items = []domain.SalesOrderItem{*orderItem}
	db.Create(order)
	db.Create(orderItem)

	return order
}

// createTestPurchaseRequest creates a test purchase request with items
func createTestPurchaseRequest(db *gorm.DB, tenantID string, user *domain.User, product *domain.Product) *domain.PurchaseRequest {
	request := &domain.PurchaseRequest{
		ID:            uuid.New().String(),
		Code:          "PR-" + uuid.New().String()[:8],
		TenantID:      tenantID,
		RequesterID:   user.ID,
		RequesterName: user.Name,
		RequestDate:   time.Now(),
		RequiredDate:  time.Now().AddDate(0, 0, 14),
		Status:        "draft",
		Priority:      "normal",
		Purpose:       "Testing purchase request",
		Active:        true,
		CreatedBy:     "system",
		UpdatedBy:     "system",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	requestItem := &domain.PurchaseRequestItem{
		ID:             uuid.New().String(),
		RequestID:      request.ID,
		TenantID:       tenantID,
		LineNumber:     1,
		ProductID:      product.ID,
		ProductCode:    product.Code,
		ProductName:    product.Name,
		ProductType:    product.Type,
		Quantity:       50,
		OrderedQty:     0,
		UOM:            product.UOM,
		EstimatedPrice: 75000,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	request.Items = []domain.PurchaseRequestItem{*requestItem}
	db.Create(request)
	db.Create(requestItem)

	return request
}

// createTestPurchaseOrder creates a test purchase order with items
func createTestPurchaseOrder(db *gorm.DB, tenantID string, vendor *domain.Vendor, product *domain.Product) *domain.PurchaseOrder {
	order := &domain.PurchaseOrder{
		ID:            uuid.New().String(),
		Code:          "PO-" + uuid.New().String()[:8],
		TenantID:      tenantID,
		VendorID:      vendor.ID,
		VendorName:    vendor.Name,
		VendorType:    vendor.Type,
		OrderDate:     time.Now(),
		ExpectedDate:  time.Now().AddDate(0, 0, 14),
		Status:        "draft",
		Priority:      "normal",
		Currency:      "IDR",
		ExchangeRate:  1.0,
		Subtotal:      4000000,
		TotalDiscount: 100000,
		TotalTax:      429000,
		ShippingCost:  50000,
		OtherCost:     20000,
		GrandTotal:    4399000,
		PaymentStatus: "unpaid",
		ReceiptStatus: "pending",
		ReceivedQty:   0,
		InvoicedQty:   0,
		Active:        true,
		CreatedBy:     "system",
		UpdatedBy:     "system",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	orderItem := &domain.PurchaseOrderItem{
		ID:           uuid.New().String(),
		OrderID:      order.ID,
		TenantID:     tenantID,
		LineNumber:   1,
		ProductID:    product.ID,
		ProductCode:  product.Code,
		ProductName:  product.Name,
		ProductType:  product.Type,
		Quantity:     50,
		ReceivedQty:  0,
		InvoicedQty:  0,
		UOM:          product.UOM,
		UnitPrice:    80000,
		DiscountType: "percentage",
		Discount:     2.5,
		TaxRate:      11,
		TaxAmount:    429000,
		Subtotal:     4000000,
		Total:        4329000,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	order.Items = []domain.PurchaseOrderItem{*orderItem}
	db.Create(order)
	db.Create(orderItem)

	return order
}

// teardownTestDB closes the test database connection
func teardownTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// assertNoError is a test helper to check for no error
func assertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: %v", msg, err)
	}
}

// assertError is a test helper to check for error
func assertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected error but got nil", msg)
	}
}

// assertEqual is a generic test helper for equality
func assertEqual(t *testing.T, got, want interface{}, msg string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", msg, got, want)
	}
}

// assertNotNil is a test helper to check for non-nil values
func assertNotNil(t *testing.T, val interface{}, msg string) {
	t.Helper()
	if val == nil {
		t.Fatalf("%s: expected non-nil value", msg)
	}
}
