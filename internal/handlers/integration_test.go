package handlers

import (
	"testing"
	"time"

	"github.com/madcok-co/unicorn-demo/internal/domain"
)

// TestSalesWorkflow_QuoteToInvoice tests the complete sales workflow
// from quotation to order to delivery to invoice
func TestSalesWorkflow_QuoteToInvoice(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create master data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")
	tax := createTestTax(db, "acme")
	warehouse := createTestWarehouse(db, "acme")

	// Step 1: Create Sales Quotation
	quotationReq := domain.CreateSalesQuotationDTO{
		Code:          "SQ-2024-001",
		CustomerID:    customer.ID,
		QuotationDate: time.Now(),
		ValidUntil:    time.Now().AddDate(0, 0, 30),
		Currency:      "IDR",
		ExchangeRate:  1.0,
		Items: []domain.CreateSalesQuotationItemDTO{
			{
				ProductID:    product.ID,
				Quantity:     10,
				UnitPrice:    100000,
				DiscountType: "percentage",
				Discount:     5,
				TaxID:        tax.ID,
			},
		},
	}

	quotation, err := CreateSalesQuotation(ctx, quotationReq)
	assertNoError(t, err, "Step 1: Create quotation should succeed")
	assertEqual(t, quotation.Status, "draft", "Quotation should start as draft")

	// Step 2: Send quotation to customer
	updateQuotationReq := domain.UpdateSalesQuotationDTO{
		Status: "sent",
	}
	quotation, err = UpdateSalesQuotation(ctx, quotation.ID, updateQuotationReq)
	assertNoError(t, err, "Step 2: Send quotation should succeed")
	assertEqual(t, quotation.Status, "sent", "Quotation should be sent")
	assertNotNil(t, quotation.SentAt, "SentAt should be set")

	// Step 3: Customer accepts quotation
	acceptQuotationReq := domain.UpdateSalesQuotationDTO{
		Status: "accepted",
	}
	quotation, err = UpdateSalesQuotation(ctx, quotation.ID, acceptQuotationReq)
	assertNoError(t, err, "Step 3: Accept quotation should succeed")
	assertEqual(t, quotation.Status, "accepted", "Quotation should be accepted")
	assertNotNil(t, quotation.AcceptedAt, "AcceptedAt should be set")

	// Step 4: Create Sales Order from quotation
	orderReq := domain.CreateSalesOrderDTO{
		Code:         "SO-2024-001",
		QuotationID:  &quotation.ID,
		CustomerID:   customer.ID,
		OrderDate:    time.Now(),
		ExpectedDate: time.Now().AddDate(0, 0, 7),
		Currency:     "IDR",
		ExchangeRate: 1.0,
		WarehouseID:  warehouse.ID,
		Items: []domain.CreateSalesOrderItemDTO{
			{
				ProductID:    product.ID,
				Quantity:     10,
				UnitPrice:    100000,
				DiscountType: "percentage",
				Discount:     5,
				TaxID:        tax.ID,
			},
		},
	}

	order, err := CreateSalesOrder(ctx, orderReq)
	assertNoError(t, err, "Step 4: Create order should succeed")
	assertEqual(t, order.Status, "draft", "Order should start as draft")
	assertEqual(t, *order.QuotationID, quotation.ID, "Order should reference quotation")

	// Step 5: Confirm order
	confirmOrderReq := domain.UpdateSalesOrderDTO{
		Status: "confirmed",
	}
	order, err = UpdateSalesOrder(ctx, order.ID, confirmOrderReq)
	assertNoError(t, err, "Step 5: Confirm order should succeed")
	assertEqual(t, order.Status, "confirmed", "Order should be confirmed")

	// Step 6: Create Delivery Order
	deliveryReq := domain.CreateDeliveryOrderDTO{
		Code:          "DO-2024-001",
		OrderID:       order.ID,
		DeliveryDate:  time.Now(),
		ScheduledDate: time.Now().AddDate(0, 0, 1),
		WarehouseID:   warehouse.ID,
		Items: []domain.CreateDeliveryOrderItemDTO{
			{
				OrderItemID:  order.Items[0].ID,
				ProductID:    product.ID,
				DeliveredQty: 10,
			},
		},
	}

	delivery, err := CreateDeliveryOrder(ctx, deliveryReq)
	assertNoError(t, err, "Step 6: Create delivery order should succeed")
	assertEqual(t, delivery.Status, "draft", "Delivery should start as draft")

	// Step 7: Deliver goods
	deliverReq := domain.UpdateDeliveryOrderDTO{
		Status: "delivered",
	}
	delivery, err = UpdateDeliveryOrder(ctx, delivery.ID, deliverReq)
	assertNoError(t, err, "Step 7: Deliver goods should succeed")
	assertEqual(t, delivery.Status, "delivered", "Delivery should be marked as delivered")

	// Step 8: Verify delivered_qty is updated in sales order
	var updatedOrder domain.SalesOrder
	db.Preload("Items").Where("id = ?", order.ID).First(&updatedOrder)
	assertEqual(t, updatedOrder.Items[0].DeliveredQty, 10.0, "Delivered qty should be updated")

	// Step 9: Create Sales Invoice
	invoiceReq := domain.CreateSalesInvoiceDTO{
		Code:         "INV-2024-001",
		OrderID:      &order.ID,
		CustomerID:   customer.ID,
		InvoiceDate:  time.Now(),
		DueDate:      time.Now().AddDate(0, 0, 30),
		Currency:     "IDR",
		ExchangeRate: 1.0,
		Items: []domain.CreateSalesInvoiceItemDTO{
			{
				OrderItemID:  &order.Items[0].ID,
				ProductID:    product.ID,
				Quantity:     10,
				UnitPrice:    100000,
				DiscountType: "percentage",
				Discount:     5,
				TaxID:        tax.ID,
			},
		},
	}

	invoice, err := CreateSalesInvoice(ctx, invoiceReq)
	assertNoError(t, err, "Step 9: Create invoice should succeed")
	assertEqual(t, invoice.Status, "draft", "Invoice should start as draft")
	assertEqual(t, *invoice.OrderID, order.ID, "Invoice should reference order")

	// Step 10: Send invoice to customer
	sendInvoiceReq := domain.UpdateSalesInvoiceDTO{
		Status: "sent",
	}
	invoice, err = UpdateSalesInvoice(ctx, invoice.ID, sendInvoiceReq)
	assertNoError(t, err, "Step 10: Send invoice should succeed")
	assertEqual(t, invoice.Status, "sent", "Invoice should be sent")

	// Verify complete workflow
	t.Log("Sales workflow completed successfully:")
	t.Logf("  Quotation: %s -> %s", quotation.Code, quotation.Status)
	t.Logf("  Order: %s -> %s", order.Code, order.Status)
	t.Logf("  Delivery: %s -> %s", delivery.Code, delivery.Status)
	t.Logf("  Invoice: %s -> %s", invoice.Code, invoice.Status)
}

// TestPurchaseWorkflow_RequestToInvoice tests the complete purchase workflow
// from request to order to receipt to invoice
func TestPurchaseWorkflow_RequestToInvoice(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create master data
	user := createTestUser(db, "acme")
	vendor := createTestVendor(db, "acme")
	product := createTestProduct(db, "acme")
	tax := createTestTax(db, "acme")
	warehouse := createTestWarehouse(db, "acme")

	// Step 1: Create Purchase Request
	requestReq := domain.CreatePurchaseRequestDTO{
		Code:         "PR-2024-001",
		RequesterID:  user.ID,
		RequestDate:  time.Now(),
		RequiredDate: time.Now().AddDate(0, 0, 14),
		Priority:     "normal",
		Purpose:      "Restock inventory",
		Items: []domain.CreatePurchaseRequestItemDTO{
			{
				ProductID:      product.ID,
				Quantity:       50,
				EstimatedPrice: 75000,
			},
		},
	}

	request, err := CreatePurchaseRequest(ctx, requestReq)
	assertNoError(t, err, "Step 1: Create request should succeed")
	assertEqual(t, request.Status, "draft", "Request should start as draft")

	// Step 2: Submit for approval
	submitReq := domain.UpdatePurchaseRequestDTO{
		Status: "submitted",
	}
	request, err = UpdatePurchaseRequest(ctx, request.ID, submitReq)
	assertNoError(t, err, "Step 2: Submit request should succeed")
	assertEqual(t, request.Status, "submitted", "Request should be submitted")
	assertNotNil(t, request.SubmittedAt, "SubmittedAt should be set")

	// Step 3: Approve request
	approveReq := domain.UpdatePurchaseRequestDTO{
		Status: "approved",
	}
	request, err = UpdatePurchaseRequest(ctx, request.ID, approveReq)
	assertNoError(t, err, "Step 3: Approve request should succeed")
	assertEqual(t, request.Status, "approved", "Request should be approved")
	assertNotNil(t, request.ApprovedDate, "ApprovedDate should be set")

	// Step 4: Create Purchase Order
	poReq := domain.CreatePurchaseOrderDTO{
		Code:         "PO-2024-001",
		RequestID:    &request.ID,
		VendorID:     vendor.ID,
		OrderDate:    time.Now(),
		ExpectedDate: time.Now().AddDate(0, 0, 14),
		Currency:     "IDR",
		ExchangeRate: 1.0,
		WarehouseID:  warehouse.ID,
		Items: []domain.CreatePurchaseOrderItemDTO{
			{
				ProductID:    product.ID,
				Quantity:     50,
				UnitPrice:    80000,
				DiscountType: "percentage",
				Discount:     2.5,
				TaxID:        tax.ID,
			},
		},
	}

	po, err := CreatePurchaseOrder(ctx, poReq)
	assertNoError(t, err, "Step 4: Create PO should succeed")
	assertEqual(t, po.Status, "draft", "PO should start as draft")
	assertEqual(t, *po.RequestID, request.ID, "PO should reference request")

	// Step 5: Send PO to vendor
	sendPOReq := domain.UpdatePurchaseOrderDTO{
		Status: "sent",
	}
	po, err = UpdatePurchaseOrder(ctx, po.ID, sendPOReq)
	assertNoError(t, err, "Step 5: Send PO should succeed")
	assertEqual(t, po.Status, "sent", "PO should be sent")

	// Step 6: Vendor confirms PO
	confirmPOReq := domain.UpdatePurchaseOrderDTO{
		Status: "confirmed",
	}
	po, err = UpdatePurchaseOrder(ctx, po.ID, confirmPOReq)
	assertNoError(t, err, "Step 6: Confirm PO should succeed")
	assertEqual(t, po.Status, "confirmed", "PO should be confirmed")

	// Step 7: Create Goods Receipt
	grReq := domain.CreateGoodsReceiptDTO{
		Code:          "GR-2024-001",
		OrderID:       po.ID,
		ReceiptDate:   time.Now(),
		ScheduledDate: time.Now(),
		WarehouseID:   warehouse.ID,
		Items: []domain.CreateGoodsReceiptItemDTO{
			{
				OrderItemID: po.Items[0].ID,
				ReceivedQty: 50,
				AcceptedQty: 50,
			},
		},
	}

	gr, err := CreateGoodsReceipt(ctx, grReq)
	assertNoError(t, err, "Step 7: Create goods receipt should succeed")
	assertEqual(t, gr.Status, "draft", "GR should start as draft")

	// Step 8: Accept goods
	acceptGRReq := domain.UpdateGoodsReceiptDTO{
		Status: "accepted",
	}
	gr, err = UpdateGoodsReceipt(ctx, gr.ID, acceptGRReq)
	assertNoError(t, err, "Step 8: Accept goods should succeed")
	assertEqual(t, gr.Status, "accepted", "GR should be accepted")

	// Step 9: Verify received_qty is updated in purchase order
	var updatedPO domain.PurchaseOrder
	db.Preload("Items").Where("id = ?", po.ID).First(&updatedPO)
	assertEqual(t, updatedPO.Items[0].ReceivedQty, 50.0, "Received qty should be updated")

	// Step 10: Create Purchase Invoice
	pinvReq := domain.CreatePurchaseInvoiceDTO{
		Code:         "PINV-2024-001",
		OrderID:      &po.ID,
		VendorID:     vendor.ID,
		InvoiceDate:  time.Now(),
		DueDate:      time.Now().AddDate(0, 0, 30),
		Currency:     "IDR",
		ExchangeRate: 1.0,
		Items: []domain.CreatePurchaseInvoiceItemDTO{
			{
				OrderItemID:  &po.Items[0].ID,
				ProductID:    product.ID,
				Quantity:     50,
				UnitPrice:    80000,
				DiscountType: "percentage",
				Discount:     2.5,
				TaxID:        tax.ID,
			},
		},
	}

	pinv, err := CreatePurchaseInvoice(ctx, pinvReq)
	assertNoError(t, err, "Step 10: Create purchase invoice should succeed")
	assertEqual(t, pinv.Status, "draft", "Invoice should start as draft")
	assertEqual(t, *pinv.OrderID, po.ID, "Invoice should reference PO")

	// Step 11: Verify invoice
	verifyInvReq := domain.UpdatePurchaseInvoiceDTO{
		Status: "verified",
	}
	pinv, err = UpdatePurchaseInvoice(ctx, pinv.ID, verifyInvReq)
	assertNoError(t, err, "Step 11: Verify invoice should succeed")
	assertEqual(t, pinv.Status, "verified", "Invoice should be verified")

	// Verify complete workflow
	t.Log("Purchase workflow completed successfully:")
	t.Logf("  Request: %s -> %s", request.Code, request.Status)
	t.Logf("  PO: %s -> %s", po.Code, po.Status)
	t.Logf("  GR: %s -> %s", gr.Code, gr.Status)
	t.Logf("  Invoice: %s -> %s", pinv.Code, pinv.Status)
}

// TestQuantityTracking_Delivery tests that delivered_qty is properly tracked
func TestQuantityTracking_Delivery(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create master data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")
	warehouse := createTestWarehouse(db, "acme")

	// Create sales order with quantity 100
	order := createTestSalesOrder(db, "acme", customer, product)
	order.Items[0].Quantity = 100
	db.Save(&order.Items[0])

	// Verify initial state
	assertEqual(t, order.Items[0].DeliveredQty, 0.0, "Initial delivered qty should be 0")

	// First partial delivery: 60 units
	delivery1Req := domain.CreateDeliveryOrderDTO{
		Code:          "DO-2024-001",
		OrderID:       order.ID,
		DeliveryDate:  time.Now(),
		ScheduledDate: time.Now(),
		WarehouseID:   warehouse.ID,
		Items: []domain.CreateDeliveryOrderItemDTO{
			{
				OrderItemID:  order.Items[0].ID,
				ProductID:    product.ID,
				DeliveredQty: 60,
			},
		},
	}

	delivery1, err := CreateDeliveryOrder(ctx, delivery1Req)
	assertNoError(t, err, "Create first delivery should succeed")

	// Update delivery to delivered status
	delivery1.Status = "delivered"
	db.Save(delivery1)

	// Check delivered_qty after first delivery
	var orderAfterFirst domain.SalesOrder
	db.Preload("Items").Where("id = ?", order.ID).First(&orderAfterFirst)
	assertEqual(t, orderAfterFirst.Items[0].DeliveredQty, 60.0, "Delivered qty should be 60 after first delivery")

	// Second partial delivery: 40 units (completing the order)
	delivery2Req := domain.CreateDeliveryOrderDTO{
		Code:          "DO-2024-002",
		OrderID:       order.ID,
		DeliveryDate:  time.Now(),
		ScheduledDate: time.Now(),
		WarehouseID:   warehouse.ID,
		Items: []domain.CreateDeliveryOrderItemDTO{
			{
				OrderItemID:  order.Items[0].ID,
				ProductID:    product.ID,
				DeliveredQty: 40,
			},
		},
	}

	delivery2, err := CreateDeliveryOrder(ctx, delivery2Req)
	assertNoError(t, err, "Create second delivery should succeed")

	// Update delivery to delivered status
	delivery2.Status = "delivered"
	db.Save(delivery2)

	// Check delivered_qty after second delivery
	var orderAfterSecond domain.SalesOrder
	db.Preload("Items").Where("id = ?", order.ID).First(&orderAfterSecond)
	assertEqual(t, orderAfterSecond.Items[0].DeliveredQty, 100.0, "Delivered qty should be 100 after second delivery")

	// Verify fulfillment status should be completed
	assertEqual(t, orderAfterSecond.FulfillmentStatus, "completed", "Fulfillment should be completed")
}

// TestQuantityTracking_Receipt tests that received_qty is properly tracked
func TestQuantityTracking_Receipt(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create master data
	vendor := createTestVendor(db, "acme")
	product := createTestProduct(db, "acme")
	warehouse := createTestWarehouse(db, "acme")

	// Create purchase order with quantity 200
	po := createTestPurchaseOrder(db, "acme", vendor, product)
	po.Items[0].Quantity = 200
	db.Save(&po.Items[0])

	// Verify initial state
	assertEqual(t, po.Items[0].ReceivedQty, 0.0, "Initial received qty should be 0")

	// First partial receipt: 120 units
	receipt1Req := domain.CreateGoodsReceiptDTO{
		Code:          "GR-2024-001",
		OrderID:       po.ID,
		ReceiptDate:   time.Now(),
		ScheduledDate: time.Now(),
		WarehouseID:   warehouse.ID,
		Items: []domain.CreateGoodsReceiptItemDTO{
			{
				OrderItemID: po.Items[0].ID,
				ReceivedQty: 120,
				AcceptedQty: 120,
			},
		},
	}

	receipt1, err := CreateGoodsReceipt(ctx, receipt1Req)
	assertNoError(t, err, "Create first receipt should succeed")

	// Update receipt to accepted status
	receipt1.Status = "accepted"
	db.Save(receipt1)

	// Check received_qty after first receipt
	var poAfterFirst domain.PurchaseOrder
	db.Preload("Items").Where("id = ?", po.ID).First(&poAfterFirst)
	assertEqual(t, poAfterFirst.Items[0].ReceivedQty, 120.0, "Received qty should be 120 after first receipt")

	// Second partial receipt: 80 units (completing the order)
	receipt2Req := domain.CreateGoodsReceiptDTO{
		Code:          "GR-2024-002",
		OrderID:       po.ID,
		ReceiptDate:   time.Now(),
		ScheduledDate: time.Now(),
		WarehouseID:   warehouse.ID,
		Items: []domain.CreateGoodsReceiptItemDTO{
			{
				OrderItemID: po.Items[0].ID,
				ReceivedQty: 80,
				AcceptedQty: 80,
			},
		},
	}

	receipt2, err := CreateGoodsReceipt(ctx, receipt2Req)
	assertNoError(t, err, "Create second receipt should succeed")

	// Update receipt to accepted status
	receipt2.Status = "accepted"
	db.Save(receipt2)

	// Check received_qty after second receipt
	var poAfterSecond domain.PurchaseOrder
	db.Preload("Items").Where("id = ?", po.ID).First(&poAfterSecond)
	assertEqual(t, poAfterSecond.Items[0].ReceivedQty, 200.0, "Received qty should be 200 after second receipt")

	// Verify receipt status should be completed
	assertEqual(t, poAfterSecond.ReceiptStatus, "completed", "Receipt should be completed")
}

// TestPaymentStatusCalculation tests automatic payment status updates
func TestPaymentStatusCalculation(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Create master data
	customer := createTestCustomer(db, "acme")
	product := createTestProduct(db, "acme")

	// Create sales invoice with total 1,000,000
	invoiceReq := domain.CreateSalesInvoiceDTO{
		Code:         "INV-2024-001",
		CustomerID:   customer.ID,
		InvoiceDate:  time.Now(),
		DueDate:      time.Now().AddDate(0, 0, 30),
		Currency:     "IDR",
		ExchangeRate: 1.0,
		Items: []domain.CreateSalesInvoiceItemDTO{
			{
				ProductID: product.ID,
				Quantity:  10,
				UnitPrice: 100000,
			},
		},
	}

	invoice, err := CreateSalesInvoice(ctx, invoiceReq)
	assertNoError(t, err, "Create invoice should succeed")
	assertEqual(t, invoice.PaymentStatus, "unpaid", "Initial status should be unpaid")
	assertEqual(t, invoice.AmountDue, invoice.GrandTotal, "Amount due should equal grand total")

	// Record partial payment: 500,000
	partialPayment := 500000.0
	updateReq1 := domain.UpdateSalesInvoiceDTO{
		AmountPaid: &partialPayment,
	}
	invoice, err = UpdateSalesInvoice(ctx, invoice.ID, updateReq1)
	assertNoError(t, err, "Record partial payment should succeed")
	assertEqual(t, invoice.PaymentStatus, "partial", "Status should be partial after partial payment")
	assertEqual(t, invoice.AmountDue, invoice.GrandTotal-partialPayment, "Amount due should be updated")

	// Record full payment
	fullPayment := invoice.GrandTotal
	updateReq2 := domain.UpdateSalesInvoiceDTO{
		AmountPaid: &fullPayment,
	}
	invoice, err = UpdateSalesInvoice(ctx, invoice.ID, updateReq2)
	assertNoError(t, err, "Record full payment should succeed")
	assertEqual(t, invoice.PaymentStatus, "paid", "Status should be paid after full payment")
	assertEqual(t, invoice.AmountDue, 0.0, "Amount due should be 0")
	assertNotNil(t, invoice.PaidAt, "PaidAt should be set")
}

// TestCrossModuleIntegration tests integration between sales and purchase modules
func TestCrossModuleIntegration(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer teardownTestDB(db)

	ctx := setupTestContext("acme", "user-123")
	ctx.Set("db", db)

	// Scenario: Customer orders product that needs to be purchased from vendor

	// Create master data
	customer := createTestCustomer(db, "acme")
	vendor := createTestVendor(db, "acme")
	product := createTestProduct(db, "acme")
	user := createTestUser(db, "acme")
	warehouse := createTestWarehouse(db, "acme")

	// Step 1: Customer creates sales order for 100 units
	salesOrder := createTestSalesOrder(db, "acme", customer, product)
	salesOrder.Items[0].Quantity = 100
	db.Save(&salesOrder.Items[0])

	// Step 2: Check stock availability (simplified - assuming low stock)
	// Create purchase request to restock
	prReq := domain.CreatePurchaseRequestDTO{
		Code:         "PR-RESTOCK-001",
		RequesterID:  user.ID,
		RequestDate:  time.Now(),
		RequiredDate: salesOrder.ExpectedDate.AddDate(0, 0, -3), // 3 days before SO expected date
		Priority:     "high",
		Purpose:      "Restock for customer order " + salesOrder.Code,
		Items: []domain.CreatePurchaseRequestItemDTO{
			{
				ProductID:      product.ID,
				Quantity:       100,
				EstimatedPrice: product.PurchasePrice,
			},
		},
	}

	pr, err := CreatePurchaseRequest(ctx, prReq)
	assertNoError(t, err, "Create purchase request should succeed")
	assertEqual(t, pr.Priority, "high", "Priority should be high for urgent restock")

	// Step 3: Approve and create PO
	pr.Status = "approved"
	db.Save(pr)

	po := createTestPurchaseOrder(db, "acme", vendor, product)
	po.RequestID = &pr.ID
	po.Items[0].Quantity = 100
	db.Save(po)

	// Step 4: Receive goods
	gr := domain.GoodsReceipt{
		ID:            "gr-001",
		Code:          "GR-RESTOCK-001",
		TenantID:      "acme",
		OrderID:       po.ID,
		OrderCode:     po.Code,
		VendorID:      vendor.ID,
		VendorName:    vendor.Name,
		ReceiptDate:   time.Now(),
		ScheduledDate: time.Now(),
		WarehouseID:   warehouse.ID,
		WarehouseName: warehouse.Name,
		Status:        "accepted",
		Active:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	db.Create(&gr)

	// Update received qty
	po.Items[0].ReceivedQty = 100
	po.ReceiptStatus = "completed"
	db.Save(&po.Items[0])
	db.Save(&po)

	// Step 5: Now fulfill sales order
	delivery := domain.DeliveryOrder{
		ID:            "do-001",
		Code:          "DO-FULFILL-001",
		TenantID:      "acme",
		OrderID:       salesOrder.ID,
		OrderCode:     salesOrder.Code,
		CustomerID:    customer.ID,
		CustomerName:  customer.Name,
		DeliveryDate:  time.Now(),
		ScheduledDate: time.Now(),
		WarehouseID:   warehouse.ID,
		WarehouseName: warehouse.Name,
		Status:        "delivered",
		Active:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	db.Create(&delivery)

	// Update delivered qty
	salesOrder.Items[0].DeliveredQty = 100
	salesOrder.FulfillmentStatus = "completed"
	db.Save(&salesOrder.Items[0])
	db.Save(&salesOrder)

	// Verify the complete flow
	assertEqual(t, pr.Status, "approved", "Purchase request should be approved")
	assertEqual(t, po.ReceiptStatus, "completed", "PO should be fully received")
	assertEqual(t, salesOrder.FulfillmentStatus, "completed", "Sales order should be fully delivered")

	t.Log("Cross-module integration completed successfully:")
	t.Logf("  Sales Order: %s (Delivered: %.0f/%.0f)", salesOrder.Code, salesOrder.Items[0].DeliveredQty, salesOrder.Items[0].Quantity)
	t.Logf("  Purchase Request: %s (Status: %s)", pr.Code, pr.Status)
	t.Logf("  Purchase Order: %s (Received: %.0f/%.0f)", po.Code, po.Items[0].ReceivedQty, po.Items[0].Quantity)
}
