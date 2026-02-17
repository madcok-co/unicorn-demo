package seeders

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"gorm.io/gorm"
)

// SeedSalesData seeds sales module data for a specific tenant
func SeedSalesData(db *gorm.DB, tenantID, userID string) error {
	// Get existing master data
	var customer domain.Customer
	if err := db.Where("tenant_id = ? AND active = ? AND code = ?", tenantID, true, "CUST001").First(&customer).Error; err != nil {
		return fmt.Errorf("customer CUST001 not found: %w", err)
	}

	var product1 domain.Product
	if err := db.Where("tenant_id = ? AND active = ? AND code = ?", tenantID, true, "PROD001").First(&product1).Error; err != nil {
		return fmt.Errorf("product PROD001 not found: %w", err)
	}

	var product2 domain.Product
	if err := db.Where("tenant_id = ? AND active = ? AND code = ?", tenantID, true, "SVC001").First(&product2).Error; err != nil {
		// If service not found, use product1 for both items
		product2 = product1
	}

	var warehouse domain.Warehouse
	if err := db.Where("tenant_id = ? AND active = ?", tenantID, true).First(&warehouse).Error; err != nil {
		return fmt.Errorf("warehouse not found: %w", err)
	}

	var tax domain.Tax
	if err := db.Where("tenant_id = ? AND active = ?", tenantID, true).First(&tax).Error; err != nil {
		return fmt.Errorf("tax not found: %w", err)
	}

	var currency domain.Currency
	if err := db.Where("tenant_id = ? AND active = ? AND is_default = ?", tenantID, true, true).First(&currency).Error; err != nil {
		return fmt.Errorf("default currency not found: %w", err)
	}

	now := time.Now()
	quotationDate := now.AddDate(0, 0, -10) // 10 days ago
	orderDate := now.AddDate(0, 0, -7)      // 7 days ago
	deliveryDate := now.AddDate(0, 0, -3)   // 3 days ago
	invoiceDate := now.AddDate(0, 0, -2)    // 2 days ago

	// ============================================================================
	// SALES QUOTATION
	// ============================================================================

	quotationID := uuid.New().String()
	quotationCode := "SQ-2024-001"

	// Calculate quotation totals
	q1Subtotal := 5.0 * 99.99       // 499.95
	q1Discount := q1Subtotal * 0.10 // 49.995 (10%)
	q1AfterDiscount := q1Subtotal - q1Discount
	q1Tax := q1AfterDiscount * 0.10    // 45.00
	q1Total := q1AfterDiscount + q1Tax // 495.00

	q2Subtotal := 2.0 * 150.00    // 300.00
	q2Tax := q2Subtotal * 0.10    // 30.00
	q2Total := q2Subtotal + q2Tax // 330.00

	qGrandTotal := q1Total + q2Total + 50.0 // 875.00 + shipping

	quotation := domain.SalesQuotation{
		ID:              quotationID,
		Code:            quotationCode,
		TenantID:        tenantID,
		CustomerID:      customer.ID,
		CustomerName:    customer.Name,
		CustomerType:    customer.Type,
		QuotationDate:   quotationDate,
		ValidUntil:      quotationDate.AddDate(0, 0, 30),
		ExpectedDate:    &orderDate,
		Status:          "accepted",
		ReferenceNo:     "REF-Q-001",
		SalespersonID:   userID,
		SalespersonName: "Sales Person",
		Currency:        currency.Code,
		ExchangeRate:    1.0,
		Subtotal:        q1Subtotal + q2Subtotal,
		TotalDiscount:   q1Discount,
		TotalTax:        q1Tax + q2Tax,
		ShippingCost:    50.0,
		OtherCost:       0.0,
		GrandTotal:      qGrandTotal,
		PaymentTermID:   "",
		PaymentTermDays: 30,
		ShippingAddress: customer.ShippingAddress,
		ShippingCity:    customer.ShippingCity,
		ShippingState:   customer.ShippingState,
		ShippingZip:     customer.ShippingZip,
		ShippingCountry: customer.ShippingCountry,
		Notes:           "Sample quotation for demonstration",
		InternalNotes:   "Customer is interested in bulk order",
		Terms:           "Payment terms: Net 30",
		Tags:            []string{"demo", "bulk"},
		Active:          true,
		CreatedBy:       userID,
		UpdatedBy:       userID,
		CreatedAt:       quotationDate,
		UpdatedAt:       quotationDate,
		AcceptedAt:      &orderDate,
	}

	quotationItems := []domain.SalesQuotationItem{
		{
			ID:           uuid.New().String(),
			QuotationID:  quotationID,
			TenantID:     tenantID,
			LineNumber:   1,
			ProductID:    product1.ID,
			ProductCode:  product1.Code,
			ProductName:  product1.Name,
			ProductType:  product1.Type,
			Description:  "High quality product",
			Quantity:     5.0,
			UOM:          product1.UOM,
			UnitPrice:    99.99,
			DiscountType: "percentage",
			Discount:     10.0,
			TaxID:        tax.ID,
			TaxRate:      tax.Rate,
			TaxAmount:    q1Tax,
			Subtotal:     q1Subtotal,
			Total:        q1Total,
			Notes:        "",
			CreatedAt:    quotationDate,
			UpdatedAt:    quotationDate,
		},
		{
			ID:           uuid.New().String(),
			QuotationID:  quotationID,
			TenantID:     tenantID,
			LineNumber:   2,
			ProductID:    product2.ID,
			ProductCode:  product2.Code,
			ProductName:  product2.Name,
			ProductType:  product2.Type,
			Description:  "Professional service",
			Quantity:     2.0,
			UOM:          product2.UOM,
			UnitPrice:    150.00,
			DiscountType: "fixed",
			Discount:     0.0,
			TaxID:        tax.ID,
			TaxRate:      tax.Rate,
			TaxAmount:    q2Tax,
			Subtotal:     q2Subtotal,
			Total:        q2Total,
			Notes:        "",
			CreatedAt:    quotationDate,
			UpdatedAt:    quotationDate,
		},
	}

	quotation.Items = quotationItems

	if err := db.Create(&quotation).Error; err != nil {
		return fmt.Errorf("failed to seed quotation: %w", err)
	}

	// ============================================================================
	// SALES ORDER (converted from quotation)
	// ============================================================================

	orderID := uuid.New().String()
	orderCode := "SO-2024-001"
	confirmedDate := orderDate.Add(time.Hour)

	salesOrder := domain.SalesOrder{
		ID:                orderID,
		Code:              orderCode,
		TenantID:          tenantID,
		QuotationID:       &quotationID,
		QuotationCode:     &quotationCode,
		CustomerID:        customer.ID,
		CustomerName:      customer.Name,
		CustomerType:      customer.Type,
		OrderDate:         orderDate,
		ExpectedDate:      orderDate.AddDate(0, 0, 14),
		ConfirmedDate:     &confirmedDate,
		Status:            "confirmed",
		Priority:          "normal",
		ReferenceNo:       "REF-SO-001",
		CustomerPO:        "PO-CUST-123",
		SalespersonID:     userID,
		SalespersonName:   "Sales Person",
		Currency:          currency.Code,
		ExchangeRate:      1.0,
		Subtotal:          q1Subtotal + q2Subtotal,
		TotalDiscount:     q1Discount,
		TotalTax:          q1Tax + q2Tax,
		ShippingCost:      50.0,
		OtherCost:         0.0,
		GrandTotal:        qGrandTotal,
		PaymentTermID:     "",
		PaymentTermDays:   30,
		PaymentStatus:     "unpaid",
		WarehouseID:       warehouse.ID,
		WarehouseName:     warehouse.Name,
		FulfillmentStatus: "partial", // Will be set after delivery
		DeliveredQty:      3.0,       // Partial delivery
		InvoicedQty:       5.0,       // Partial invoice
		ShippingMethod:    "Standard Shipping",
		ShippingAddress:   customer.ShippingAddress,
		ShippingCity:      customer.ShippingCity,
		ShippingState:     customer.ShippingState,
		ShippingZip:       customer.ShippingZip,
		ShippingCountry:   customer.ShippingCountry,
		Notes:             "Converted from quotation SQ-2024-001",
		InternalNotes:     "Priority customer",
		Terms:             "Payment terms: Net 30",
		Tags:              []string{"demo", "converted"},
		Active:            true,
		CreatedBy:         userID,
		UpdatedBy:         userID,
		CreatedAt:         orderDate,
		UpdatedAt:         orderDate,
	}

	orderItem1ID := uuid.New().String()
	orderItem2ID := uuid.New().String()

	orderItems := []domain.SalesOrderItem{
		{
			ID:           orderItem1ID,
			OrderID:      orderID,
			TenantID:     tenantID,
			LineNumber:   1,
			ProductID:    product1.ID,
			ProductCode:  product1.Code,
			ProductName:  product1.Name,
			ProductType:  product1.Type,
			Description:  "High quality product",
			Quantity:     5.0,
			DeliveredQty: 3.0, // Partial delivery
			InvoicedQty:  5.0, // Fully invoiced
			UOM:          product1.UOM,
			UnitPrice:    99.99,
			DiscountType: "percentage",
			Discount:     10.0,
			TaxID:        tax.ID,
			TaxRate:      tax.Rate,
			TaxAmount:    q1Tax,
			Subtotal:     q1Subtotal,
			Total:        q1Total,
			Notes:        "",
			CreatedAt:    orderDate,
			UpdatedAt:    orderDate,
		},
		{
			ID:           orderItem2ID,
			OrderID:      orderID,
			TenantID:     tenantID,
			LineNumber:   2,
			ProductID:    product2.ID,
			ProductCode:  product2.Code,
			ProductName:  product2.Name,
			ProductType:  product2.Type,
			Description:  "Professional service",
			Quantity:     2.0,
			DeliveredQty: 0.0, // Service not delivered yet
			InvoicedQty:  0.0,
			UOM:          product2.UOM,
			UnitPrice:    150.00,
			DiscountType: "fixed",
			Discount:     0.0,
			TaxID:        tax.ID,
			TaxRate:      tax.Rate,
			TaxAmount:    q2Tax,
			Subtotal:     q2Subtotal,
			Total:        q2Total,
			Notes:        "",
			CreatedAt:    orderDate,
			UpdatedAt:    orderDate,
		},
	}

	salesOrder.Items = orderItems

	if err := db.Create(&salesOrder).Error; err != nil {
		return fmt.Errorf("failed to seed sales order: %w", err)
	}

	// ============================================================================
	// DELIVERY ORDER (partial delivery for product1)
	// ============================================================================

	deliveryID := uuid.New().String()
	actualDate := deliveryDate.Add(time.Hour * 2)

	deliveryOrder := domain.DeliveryOrder{
		ID:              deliveryID,
		Code:            "DO-2024-001",
		TenantID:        tenantID,
		OrderID:         orderID,
		OrderCode:       orderCode,
		CustomerID:      customer.ID,
		CustomerName:    customer.Name,
		DeliveryDate:    deliveryDate,
		ScheduledDate:   deliveryDate,
		ActualDate:      &actualDate,
		Status:          "delivered",
		ReferenceNo:     "REF-DO-001",
		TrackingNumber:  "TRACK-123456",
		WarehouseID:     warehouse.ID,
		WarehouseName:   warehouse.Name,
		ShippingMethod:  "Standard Shipping",
		ShippingCost:    50.0,
		CourierName:     "Express Courier",
		DriverName:      "John Driver",
		VehicleNumber:   "B-1234-XYZ",
		ShippingAddress: customer.ShippingAddress,
		ShippingCity:    customer.ShippingCity,
		ShippingState:   customer.ShippingState,
		ShippingZip:     customer.ShippingZip,
		ShippingCountry: customer.ShippingCountry,
		RecipientName:   customer.ContactPerson,
		RecipientPhone:  customer.Phone,
		RecipientEmail:  customer.Email,
		Notes:           "Partial delivery - 3 units of product",
		InternalNotes:   "Customer requested partial delivery",
		Tags:            []string{"demo", "partial"},
		Active:          true,
		CreatedBy:       userID,
		UpdatedBy:       userID,
		CreatedAt:       deliveryDate,
		UpdatedAt:       deliveryDate,
		ReceivedBy:      customer.ContactPerson,
		ReceivedAt:      &actualDate,
	}

	deliveryItems := []domain.DeliveryOrderItem{
		{
			ID:              uuid.New().String(),
			DeliveryOrderID: deliveryID,
			TenantID:        tenantID,
			OrderItemID:     orderItem1ID,
			LineNumber:      1,
			ProductID:       product1.ID,
			ProductCode:     product1.Code,
			ProductName:     product1.Name,
			Description:     "High quality product",
			OrderedQty:      5.0,
			DeliveredQty:    3.0,
			UOM:             product1.UOM,
			SerialNumbers:   []string{"SN001", "SN002", "SN003"},
			BatchNumbers:    []string{"BATCH-2024-001"},
			Notes:           "Remaining 2 units to be delivered later",
			CreatedAt:       deliveryDate,
			UpdatedAt:       deliveryDate,
		},
	}

	deliveryOrder.Items = deliveryItems

	if err := db.Create(&deliveryOrder).Error; err != nil {
		return fmt.Errorf("failed to seed delivery order: %w", err)
	}

	// ============================================================================
	// SALES INVOICE (partial invoice for product1)
	// ============================================================================

	invoiceID := uuid.New().String()
	dueDate := invoiceDate.AddDate(0, 0, 30)

	salesInvoice := domain.SalesInvoice{
		ID:                invoiceID,
		Code:              "INV-2024-001",
		TenantID:          tenantID,
		OrderID:           &orderID,
		OrderCode:         &orderCode,
		DeliveryOrderID:   &deliveryID,
		DeliveryOrderCode: &deliveryOrder.Code,
		CustomerID:        customer.ID,
		CustomerName:      customer.Name,
		CustomerType:      customer.Type,
		CustomerTaxID:     customer.TaxID,
		InvoiceDate:       invoiceDate,
		DueDate:           dueDate,
		Status:            "sent",
		ReferenceNo:       "REF-INV-001",
		CustomerPO:        "PO-CUST-123",
		TaxInvoiceNo:      "FP-001-2024",
		Currency:          currency.Code,
		ExchangeRate:      1.0,
		Subtotal:          q1Subtotal,
		TotalDiscount:     q1Discount,
		TotalTax:          q1Tax,
		ShippingCost:      50.0,
		OtherCost:         0.0,
		GrandTotal:        q1Total + 50.0,
		AmountPaid:        0.0,
		AmountDue:         q1Total + 50.0,
		PaymentTermID:     "",
		PaymentTermDays:   30,
		PaymentStatus:     "unpaid",
		IncomeAccountID:   "",
		ARAccountID:       "",
		BillingAddress:    customer.BillingAddress,
		BillingCity:       customer.BillingCity,
		BillingState:      customer.BillingState,
		BillingZip:        customer.BillingZip,
		BillingCountry:    customer.BillingCountry,
		Notes:             "Invoice for partial delivery",
		InternalNotes:     "Payment expected within 30 days",
		Terms:             "Payment terms: Net 30. Late payment subject to 2% interest per month.",
		Tags:              []string{"demo", "partial"},
		Active:            true,
		CreatedBy:         userID,
		UpdatedBy:         userID,
		CreatedAt:         invoiceDate,
		UpdatedAt:         invoiceDate,
		SentAt:            &invoiceDate,
	}

	invoiceItems := []domain.SalesInvoiceItem{
		{
			ID:           uuid.New().String(),
			InvoiceID:    invoiceID,
			TenantID:     tenantID,
			OrderItemID:  &orderItem1ID,
			LineNumber:   1,
			ProductID:    product1.ID,
			ProductCode:  product1.Code,
			ProductName:  product1.Name,
			ProductType:  product1.Type,
			Description:  "High quality product",
			Quantity:     5.0,
			UOM:          product1.UOM,
			UnitPrice:    99.99,
			DiscountType: "percentage",
			Discount:     10.0,
			TaxID:        tax.ID,
			TaxRate:      tax.Rate,
			TaxAmount:    q1Tax,
			Subtotal:     q1Subtotal,
			Total:        q1Total,
			AccountID:    "",
			Notes:        "Full quantity invoiced",
			CreatedAt:    invoiceDate,
			UpdatedAt:    invoiceDate,
		},
	}

	salesInvoice.Items = invoiceItems

	if err := db.Create(&salesInvoice).Error; err != nil {
		return fmt.Errorf("failed to seed sales invoice: %w", err)
	}

	return nil
}
