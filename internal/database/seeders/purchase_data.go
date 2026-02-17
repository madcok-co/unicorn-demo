package seeders

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"gorm.io/gorm"
)

// SeedPurchaseData seeds purchase module data for a specific tenant
func SeedPurchaseData(db *gorm.DB, tenantID, userID string) error {
	// Get existing master data
	var vendor domain.Vendor
	if err := db.Where("tenant_id = ? AND active = ? AND code = ?", tenantID, true, "VEND001").First(&vendor).Error; err != nil {
		return fmt.Errorf("vendor VEND001 not found: %w", err)
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
	requestDate := now.AddDate(0, 0, -15)   // 15 days ago
	approvedDate := now.AddDate(0, 0, -13)  // 13 days ago
	orderDate := now.AddDate(0, 0, -12)     // 12 days ago
	confirmedDate := now.AddDate(0, 0, -11) // 11 days ago
	receiptDate := now.AddDate(0, 0, -5)    // 5 days ago
	invoiceDate := now.AddDate(0, 0, -3)    // 3 days ago

	// ============================================================================
	// PURCHASE REQUEST
	// ============================================================================

	requestID := uuid.New().String()
	requestCode := "PR-2024-001"
	requiredDate := requestDate.AddDate(0, 0, 14)

	purchaseRequest := domain.PurchaseRequest{
		ID:             requestID,
		Code:           requestCode,
		TenantID:       tenantID,
		RequesterID:    userID,
		RequesterName:  "Requester User",
		DepartmentID:   "",
		DepartmentName: "Operations",
		RequestDate:    requestDate,
		RequiredDate:   requiredDate,
		ApprovedDate:   &approvedDate,
		Status:         "approved",
		Priority:       "normal",
		ReferenceNo:    "REF-PR-001",
		Purpose:        "Stock replenishment",
		ApproverID:     userID,
		ApproverName:   "Approver User",
		Notes:          "Sample purchase request for demonstration",
		InternalNotes:  "Urgent stock needed",
		Tags:           []string{"demo", "stock-replenishment"},
		Active:         true,
		CreatedBy:      userID,
		UpdatedBy:      userID,
		CreatedAt:      requestDate,
		UpdatedAt:      requestDate,
		SubmittedAt:    &requestDate,
	}

	requestItems := []domain.PurchaseRequestItem{
		{
			ID:             uuid.New().String(),
			RequestID:      requestID,
			TenantID:       tenantID,
			LineNumber:     1,
			ProductID:      product1.ID,
			ProductCode:    product1.Code,
			ProductName:    product1.Name,
			ProductType:    product1.Type,
			Description:    "Quality product needed for stock",
			Quantity:       5.0,
			OrderedQty:     5.0, // Fully converted to PO
			UOM:            product1.UOM,
			EstimatedPrice: 75.00,
			Notes:          "",
			CreatedAt:      requestDate,
			UpdatedAt:      requestDate,
		},
		{
			ID:             uuid.New().String(),
			RequestID:      requestID,
			TenantID:       tenantID,
			LineNumber:     2,
			ProductID:      product2.ID,
			ProductCode:    product2.Code,
			ProductName:    product2.Name,
			ProductType:    product2.Type,
			Description:    "Service required",
			Quantity:       2.0,
			OrderedQty:     2.0, // Fully converted to PO
			UOM:            product2.UOM,
			EstimatedPrice: 100.00,
			Notes:          "",
			CreatedAt:      requestDate,
			UpdatedAt:      requestDate,
		},
	}

	purchaseRequest.Items = requestItems

	if err := db.Create(&purchaseRequest).Error; err != nil {
		return fmt.Errorf("failed to seed purchase request: %w", err)
	}

	// ============================================================================
	// PURCHASE ORDER (converted from purchase request)
	// ============================================================================

	orderID := uuid.New().String()
	orderCode := "PO-2024-001"
	expectedDate := orderDate.AddDate(0, 0, 14)

	// Calculate order totals
	// Item 1: qty=5, price=75, discount=5%, tax=10%
	p1Subtotal := 5.0 * 75.00                  // 375.00
	p1Discount := p1Subtotal * 0.05            // 18.75 (5%)
	p1AfterDiscount := p1Subtotal - p1Discount // 356.25
	p1Tax := p1AfterDiscount * 0.10            // 35.625 (10%)
	p1Total := p1AfterDiscount + p1Tax         // 391.875

	// Item 2: qty=2, price=100, no discount, tax=10%
	p2Subtotal := 2.0 * 100.00    // 200.00
	p2Tax := p2Subtotal * 0.10    // 20.00 (10%)
	p2Total := p2Subtotal + p2Tax // 220.00

	pGrandTotal := p1Total + p2Total + 50.0 // 661.875 + shipping

	purchaseOrder := domain.PurchaseOrder{
		ID:              orderID,
		Code:            orderCode,
		TenantID:        tenantID,
		RequestID:       &requestID,
		RequestCode:     &requestCode,
		VendorID:        vendor.ID,
		VendorName:      vendor.Name,
		VendorType:      vendor.Type,
		OrderDate:       orderDate,
		ExpectedDate:    expectedDate,
		ConfirmedDate:   &confirmedDate,
		Status:          "confirmed",
		Priority:        "normal",
		ReferenceNo:     "REF-PO-001",
		VendorQuoteNo:   "VQ-2024-001",
		BuyerID:         userID,
		BuyerName:       "Buyer User",
		Currency:        currency.Code,
		ExchangeRate:    1.0,
		Subtotal:        p1Subtotal + p2Subtotal,
		TotalDiscount:   p1Discount,
		TotalTax:        p1Tax + p2Tax,
		ShippingCost:    50.0,
		OtherCost:       0.0,
		GrandTotal:      pGrandTotal,
		PaymentTermID:   "",
		PaymentTermDays: 30,
		PaymentStatus:   "unpaid",
		WarehouseID:     warehouse.ID,
		WarehouseName:   warehouse.Name,
		ReceiptStatus:   "partial", // Will be set after partial receipt
		ReceivedQty:     3.0,       // Partial receipt for item 1
		InvoicedQty:     5.0,       // Partial invoice for item 1
		ShippingMethod:  "Standard Shipping",
		ShippingAddress: warehouse.Address,
		ShippingCity:    warehouse.City,
		ShippingState:   warehouse.State,
		ShippingZip:     warehouse.Zip,
		ShippingCountry: warehouse.Country,
		Notes:           "Converted from purchase request PR-2024-001",
		InternalNotes:   "Preferred vendor",
		Terms:           "Payment terms: Net 30",
		Tags:            []string{"demo", "converted"},
		Active:          true,
		CreatedBy:       userID,
		UpdatedBy:       userID,
		CreatedAt:       orderDate,
		UpdatedAt:       orderDate,
	}

	orderItem1ID := uuid.New().String()
	orderItem2ID := uuid.New().String()

	orderItems := []domain.PurchaseOrderItem{
		{
			ID:           orderItem1ID,
			OrderID:      orderID,
			TenantID:     tenantID,
			LineNumber:   1,
			ProductID:    product1.ID,
			ProductCode:  product1.Code,
			ProductName:  product1.Name,
			ProductType:  product1.Type,
			Description:  "Quality product needed for stock",
			Quantity:     5.0,
			ReceivedQty:  3.0, // Partial receipt
			InvoicedQty:  5.0, // Fully invoiced
			UOM:          product1.UOM,
			UnitPrice:    75.00,
			DiscountType: "percentage",
			Discount:     5.0,
			TaxID:        tax.ID,
			TaxRate:      tax.Rate,
			TaxAmount:    p1Tax,
			Subtotal:     p1Subtotal,
			Total:        p1Total,
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
			Description:  "Service required",
			Quantity:     2.0,
			ReceivedQty:  0.0, // Service not received yet
			InvoicedQty:  0.0,
			UOM:          product2.UOM,
			UnitPrice:    100.00,
			DiscountType: "fixed",
			Discount:     0.0,
			TaxID:        tax.ID,
			TaxRate:      tax.Rate,
			TaxAmount:    p2Tax,
			Subtotal:     p2Subtotal,
			Total:        p2Total,
			Notes:        "",
			CreatedAt:    orderDate,
			UpdatedAt:    orderDate,
		},
	}

	purchaseOrder.Items = orderItems

	if err := db.Create(&purchaseOrder).Error; err != nil {
		return fmt.Errorf("failed to seed purchase order: %w", err)
	}

	// ============================================================================
	// GOODS RECEIPT (partial receipt for product1)
	// ============================================================================

	receiptID := uuid.New().String()
	actualDate := receiptDate.Add(time.Hour * 2)

	goodsReceipt := domain.GoodsReceipt{
		ID:              receiptID,
		Code:            "GR-2024-001",
		TenantID:        tenantID,
		OrderID:         orderID,
		OrderCode:       orderCode,
		VendorID:        vendor.ID,
		VendorName:      vendor.Name,
		ReceiptDate:     receiptDate,
		ScheduledDate:   receiptDate,
		ActualDate:      &actualDate,
		Status:          "accepted",
		ReferenceNo:     "REF-GR-001",
		DeliveryNote:    "DN-VEND-123",
		PackingList:     "PL-VEND-123",
		WarehouseID:     warehouse.ID,
		WarehouseName:   warehouse.Name,
		ReceiverID:      userID,
		ReceiverName:    "Receiver User",
		InspectorID:     userID,
		InspectorName:   "Quality Inspector",
		QualityStatus:   "passed",
		RejectedQty:     0.0,
		RejectionReason: "",
		Notes:           "Partial receipt - 3 units of product",
		InternalNotes:   "Remaining 2 units to be delivered later",
		Tags:            []string{"demo", "partial"},
		Active:          true,
		DeliveryProof:   []string{},
		CreatedBy:       userID,
		UpdatedBy:       userID,
		CreatedAt:       receiptDate,
		UpdatedAt:       receiptDate,
	}

	receiptItems := []domain.GoodsReceiptItem{
		{
			ID:            uuid.New().String(),
			ReceiptID:     receiptID,
			TenantID:      tenantID,
			OrderItemID:   orderItem1ID,
			LineNumber:    1,
			ProductID:     product1.ID,
			ProductCode:   product1.Code,
			ProductName:   product1.Name,
			Description:   "Quality product needed for stock",
			OrderedQty:    5.0,
			ReceivedQty:   3.0,
			AcceptedQty:   3.0,
			RejectedQty:   0.0,
			UOM:           product1.UOM,
			SerialNumbers: []string{"SN-P001", "SN-P002", "SN-P003"},
			BatchNumbers:  []string{"BATCH-P-001"},
			ExpiryDates:   []string{},
			Notes:         "Remaining 2 units to be delivered later",
			CreatedAt:     receiptDate,
			UpdatedAt:     receiptDate,
		},
	}

	goodsReceipt.Items = receiptItems

	if err := db.Create(&goodsReceipt).Error; err != nil {
		return fmt.Errorf("failed to seed goods receipt: %w", err)
	}

	// ============================================================================
	// PURCHASE INVOICE (invoice for full quantity of product1)
	// ============================================================================

	invoiceID := uuid.New().String()
	dueDate := invoiceDate.AddDate(0, 0, 30)

	purchaseInvoice := domain.PurchaseInvoice{
		ID:               invoiceID,
		Code:             "PINV-2024-001",
		TenantID:         tenantID,
		OrderID:          &orderID,
		OrderCode:        &orderCode,
		ReceiptID:        &receiptID,
		ReceiptCode:      &goodsReceipt.Code,
		VendorID:         vendor.ID,
		VendorName:       vendor.Name,
		VendorType:       vendor.Type,
		VendorTaxID:      vendor.TaxID,
		VendorInvoiceNo:  "VINV-2024-001",
		InvoiceDate:      invoiceDate,
		DueDate:          dueDate,
		Status:           "received",
		ReferenceNo:      "REF-PINV-001",
		TaxInvoiceNo:     "FP-VEND-001-2024",
		Currency:         currency.Code,
		ExchangeRate:     1.0,
		Subtotal:         p1Subtotal,
		TotalDiscount:    p1Discount,
		TotalTax:         p1Tax,
		ShippingCost:     50.0,
		OtherCost:        0.0,
		GrandTotal:       p1Total + 50.0,
		AmountPaid:       0.0,
		AmountDue:        p1Total + 50.0,
		PaymentTermID:    "",
		PaymentTermDays:  30,
		PaymentStatus:    "unpaid",
		ExpenseAccountID: "",
		APAccountID:      "",
		JournalEntryID:   nil,
		Notes:            "Invoice for full quantity (5 units)",
		InternalNotes:    "Payment due within 30 days",
		Terms:            "Payment terms: Net 30. Late payment subject to 2% interest per month.",
		Tags:             []string{"demo", "partial-receipt"},
		Active:           true,
		CreatedBy:        userID,
		UpdatedBy:        userID,
		CreatedAt:        invoiceDate,
		UpdatedAt:        invoiceDate,
		ReceivedAt:       &invoiceDate,
	}

	invoiceItems := []domain.PurchaseInvoiceItem{
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
			Description:  "Quality product needed for stock",
			Quantity:     5.0,
			UOM:          product1.UOM,
			UnitPrice:    75.00,
			DiscountType: "percentage",
			Discount:     5.0,
			TaxID:        tax.ID,
			TaxRate:      tax.Rate,
			TaxAmount:    p1Tax,
			Subtotal:     p1Subtotal,
			Total:        p1Total,
			AccountID:    "",
			Notes:        "Full quantity invoiced (partial receipt: 3 of 5)",
			CreatedAt:    invoiceDate,
			UpdatedAt:    invoiceDate,
		},
	}

	purchaseInvoice.Items = invoiceItems

	if err := db.Create(&purchaseInvoice).Error; err != nil {
		return fmt.Errorf("failed to seed purchase invoice: %w", err)
	}

	return nil
}
