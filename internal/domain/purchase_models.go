package domain

import "time"

// ============================================================================
// PURCHASE MODULE MODELS
// ============================================================================

// PurchaseRequest represents a purchase request/requisition
type PurchaseRequest struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_purchase_request_code_tenant;not null"` // PR-2024-001
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_purchase_request_code_tenant;index;not null"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Requester Information
	RequesterID    string `json:"requester_id" gorm:"index;not null"`
	RequesterName  string `json:"requester_name" gorm:"not null"`
	DepartmentID   string `json:"department_id" gorm:"index"`
	DepartmentName string `json:"department_name"`

	// Document Details
	RequestDate  time.Time  `json:"request_date" gorm:"not null"`
	RequiredDate time.Time  `json:"required_date" gorm:"not null"`
	ApprovedDate *time.Time `json:"approved_date"`
	Status       string     `json:"status" gorm:"default:'draft'"`    // draft, submitted, approved, rejected, cancelled, completed
	Priority     string     `json:"priority" gorm:"default:'normal'"` // low, normal, high, urgent
	ReferenceNo  string     `json:"reference_no"`
	Purpose      string     `json:"purpose" gorm:"type:text"`

	// Approval
	ApproverID   string `json:"approver_id" gorm:"index"`
	ApproverName string `json:"approver_name"`

	// Additional Info
	Notes         string   `json:"notes" gorm:"type:text"`
	InternalNotes string   `json:"internal_notes" gorm:"type:text"`
	Tags          []string `json:"tags" gorm:"serializer:json"`

	// Line Items
	Items []PurchaseRequestItem `json:"items" gorm:"foreignKey:RequestID;constraint:OnDelete:CASCADE"`

	// Metadata
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Workflow tracking
	SubmittedAt *time.Time `json:"submitted_at"`
	RejectedAt  *time.Time `json:"rejected_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

func (PurchaseRequest) TableName() string {
	return "purchase_requests"
}

// PurchaseRequestItem represents a line item in purchase request
type PurchaseRequestItem struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	RequestID      string    `json:"request_id" gorm:"index;not null"`
	TenantID       string    `json:"tenant_id" gorm:"index;not null"`
	LineNumber     int       `json:"line_number" gorm:"not null"`
	ProductID      string    `json:"product_id" gorm:"index;not null"`
	ProductCode    string    `json:"product_code" gorm:"not null"`
	ProductName    string    `json:"product_name" gorm:"not null"`
	ProductType    string    `json:"product_type"`
	Description    string    `json:"description" gorm:"type:text"`
	Quantity       float64   `json:"quantity" gorm:"not null"`
	OrderedQty     float64   `json:"ordered_qty" gorm:"default:0"` // Qty converted to PO
	UOM            string    `json:"uom" gorm:"not null"`
	EstimatedPrice float64   `json:"estimated_price"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (PurchaseRequestItem) TableName() string {
	return "purchase_request_items"
}

// PurchaseOrder represents a purchase order to vendor
type PurchaseOrder struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_purchase_order_code_tenant;not null"` // PO-2024-001
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_purchase_order_code_tenant;index;not null"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Reference
	RequestID   *string `json:"request_id" gorm:"index"`
	RequestCode *string `json:"request_code"`

	// Vendor Information
	VendorID   string `json:"vendor_id" gorm:"index;not null"`
	VendorName string `json:"vendor_name" gorm:"not null"`
	VendorType string `json:"vendor_type"`

	// Document Details
	OrderDate     time.Time  `json:"order_date" gorm:"not null"`
	ExpectedDate  time.Time  `json:"expected_date" gorm:"not null"`
	ConfirmedDate *time.Time `json:"confirmed_date"`
	Status        string     `json:"status" gorm:"default:'draft'"`    // draft, sent, confirmed, receiving, completed, cancelled
	Priority      string     `json:"priority" gorm:"default:'normal'"` // low, normal, high, urgent
	ReferenceNo   string     `json:"reference_no"`
	VendorQuoteNo string     `json:"vendor_quote_no"` // Vendor's quotation number
	BuyerID       string     `json:"buyer_id" gorm:"index"`
	BuyerName     string     `json:"buyer_name"`

	// Financial
	Currency      string  `json:"currency" gorm:"default:'IDR'"`
	ExchangeRate  float64 `json:"exchange_rate" gorm:"default:1"`
	Subtotal      float64 `json:"subtotal" gorm:"default:0"`
	TotalDiscount float64 `json:"total_discount" gorm:"default:0"`
	TotalTax      float64 `json:"total_tax" gorm:"default:0"`
	ShippingCost  float64 `json:"shipping_cost" gorm:"default:0"`
	OtherCost     float64 `json:"other_cost" gorm:"default:0"`
	GrandTotal    float64 `json:"grandtotal" gorm:"default:0"`

	// Payment
	PaymentTermID   string `json:"payment_term_id"`
	PaymentTermDays int    `json:"payment_term_days" gorm:"default:30"`
	PaymentStatus   string `json:"payment_status" gorm:"default:'unpaid'"` // unpaid, partial, paid

	// Receipt
	WarehouseID   string  `json:"warehouse_id" gorm:"index"`
	WarehouseName string  `json:"warehouse_name"`
	ReceiptStatus string  `json:"receipt_status" gorm:"default:'pending'"` // pending, partial, completed
	ReceivedQty   float64 `json:"received_qty" gorm:"default:0"`
	InvoicedQty   float64 `json:"invoiced_qty" gorm:"default:0"`

	// Shipping
	ShippingMethod  string `json:"shipping_method"`
	ShippingAddress string `json:"shipping_address"`
	ShippingCity    string `json:"shipping_city"`
	ShippingState   string `json:"shipping_state"`
	ShippingZip     string `json:"shipping_zip"`
	ShippingCountry string `json:"shipping_country"`

	// Additional Info
	Notes         string   `json:"notes" gorm:"type:text"`
	InternalNotes string   `json:"internal_notes" gorm:"type:text"`
	Terms         string   `json:"terms" gorm:"type:text"`
	Tags          []string `json:"tags" gorm:"serializer:json"`

	// Line Items
	Items []PurchaseOrderItem `json:"items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`

	// Metadata
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}

// PurchaseOrderItem represents a line item in purchase order
type PurchaseOrderItem struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	OrderID      string    `json:"order_id" gorm:"index;not null"`
	TenantID     string    `json:"tenant_id" gorm:"index;not null"`
	LineNumber   int       `json:"line_number" gorm:"not null"`
	ProductID    string    `json:"product_id" gorm:"index;not null"`
	ProductCode  string    `json:"product_code" gorm:"not null"`
	ProductName  string    `json:"product_name" gorm:"not null"`
	ProductType  string    `json:"product_type"`
	Description  string    `json:"description" gorm:"type:text"`
	Quantity     float64   `json:"quantity" gorm:"not null"`
	ReceivedQty  float64   `json:"received_qty" gorm:"default:0"`
	InvoicedQty  float64   `json:"invoiced_qty" gorm:"default:0"`
	UOM          string    `json:"uom" gorm:"not null"`
	UnitPrice    float64   `json:"unit_price" gorm:"not null"`
	DiscountType string    `json:"discount_type" gorm:"default:'percentage'"`
	Discount     float64   `json:"discount" gorm:"default:0"`
	TaxID        string    `json:"tax_id"`
	TaxRate      float64   `json:"tax_rate" gorm:"default:0"`
	TaxAmount    float64   `json:"tax_amount" gorm:"default:0"`
	Subtotal     float64   `json:"subtotal" gorm:"not null"`
	Total        float64   `json:"total" gorm:"not null"`
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (PurchaseOrderItem) TableName() string {
	return "purchase_order_items"
}

// GoodsReceipt represents received goods from vendor
type GoodsReceipt struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_goods_receipt_code_tenant;not null"` // GR-2024-001
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_goods_receipt_code_tenant;index;not null"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Reference
	OrderID   string `json:"order_id" gorm:"index;not null"`
	OrderCode string `json:"order_code" gorm:"not null"`

	// Vendor Information
	VendorID   string `json:"vendor_id" gorm:"index;not null"`
	VendorName string `json:"vendor_name" gorm:"not null"`

	// Document Details
	ReceiptDate   time.Time  `json:"receipt_date" gorm:"not null"`
	ScheduledDate time.Time  `json:"scheduled_date" gorm:"not null"`
	ActualDate    *time.Time `json:"actual_date"`
	Status        string     `json:"status" gorm:"default:'draft'"` // draft, received, inspected, accepted, rejected
	ReferenceNo   string     `json:"reference_no"`
	DeliveryNote  string     `json:"delivery_note"` // Vendor's delivery note number
	PackingList   string     `json:"packing_list"`

	// Warehouse & Receipt
	WarehouseID   string `json:"warehouse_id" gorm:"index;not null"`
	WarehouseName string `json:"warehouse_name" gorm:"not null"`
	ReceiverID    string `json:"receiver_id" gorm:"index"`
	ReceiverName  string `json:"receiver_name"`
	InspectorID   string `json:"inspector_id" gorm:"index"`
	InspectorName string `json:"inspector_name"`

	// Quality Check
	QualityStatus   string  `json:"quality_status" gorm:"default:'pending'"` // pending, passed, failed, partial
	RejectedQty     float64 `json:"rejected_qty" gorm:"default:0"`
	RejectionReason string  `json:"rejection_reason" gorm:"type:text"`

	// Additional Info
	Notes         string   `json:"notes" gorm:"type:text"`
	InternalNotes string   `json:"internal_notes" gorm:"type:text"`
	Tags          []string `json:"tags" gorm:"serializer:json"`

	// Line Items
	Items []GoodsReceiptItem `json:"items" gorm:"foreignKey:ReceiptID;constraint:OnDelete:CASCADE"`

	// Attachments
	DeliveryProof []string `json:"delivery_proof" gorm:"serializer:json"` // URLs to proof photos

	// Metadata
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (GoodsReceipt) TableName() string {
	return "goods_receipts"
}

// GoodsReceiptItem represents a line item in goods receipt
type GoodsReceiptItem struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	ReceiptID     string    `json:"receipt_id" gorm:"index;not null"`
	TenantID      string    `json:"tenant_id" gorm:"index;not null"`
	OrderItemID   string    `json:"order_item_id" gorm:"index;not null"`
	LineNumber    int       `json:"line_number" gorm:"not null"`
	ProductID     string    `json:"product_id" gorm:"index;not null"`
	ProductCode   string    `json:"product_code" gorm:"not null"`
	ProductName   string    `json:"product_name" gorm:"not null"`
	Description   string    `json:"description"`
	OrderedQty    float64   `json:"ordered_qty" gorm:"not null"`
	ReceivedQty   float64   `json:"received_qty" gorm:"not null"`
	AcceptedQty   float64   `json:"accepted_qty" gorm:"default:0"`
	RejectedQty   float64   `json:"rejected_qty" gorm:"default:0"`
	UOM           string    `json:"uom" gorm:"not null"`
	SerialNumbers []string  `json:"serial_numbers" gorm:"serializer:json"`
	BatchNumbers  []string  `json:"batch_numbers" gorm:"serializer:json"`
	ExpiryDates   []string  `json:"expiry_dates" gorm:"serializer:json"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (GoodsReceiptItem) TableName() string {
	return "goods_receipt_items"
}

// PurchaseInvoice represents a vendor invoice
type PurchaseInvoice struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_purchase_invoice_code_tenant;not null"` // PINV-2024-001
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_purchase_invoice_code_tenant;index;not null"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Reference
	OrderID     *string `json:"order_id" gorm:"index"`
	OrderCode   *string `json:"order_code"`
	ReceiptID   *string `json:"receipt_id" gorm:"index"`
	ReceiptCode *string `json:"receipt_code"`

	// Vendor Information
	VendorID        string `json:"vendor_id" gorm:"index;not null"`
	VendorName      string `json:"vendor_name" gorm:"not null"`
	VendorType      string `json:"vendor_type"`
	VendorTaxID     string `json:"vendor_tax_id"`
	VendorInvoiceNo string `json:"vendor_invoice_no"` // Vendor's invoice number

	// Document Details
	InvoiceDate  time.Time `json:"invoice_date" gorm:"not null"`
	DueDate      time.Time `json:"due_date" gorm:"not null"`
	Status       string    `json:"status" gorm:"default:'draft'"` // draft, received, verified, approved, partial_paid, paid, cancelled
	ReferenceNo  string    `json:"reference_no"`
	TaxInvoiceNo string    `json:"tax_invoice_no"` // Faktur Pajak from vendor

	// Financial
	Currency      string  `json:"currency" gorm:"default:'IDR'"`
	ExchangeRate  float64 `json:"exchange_rate" gorm:"default:1"`
	Subtotal      float64 `json:"subtotal" gorm:"default:0"`
	TotalDiscount float64 `json:"total_discount" gorm:"default:0"`
	TotalTax      float64 `json:"total_tax" gorm:"default:0"`
	ShippingCost  float64 `json:"shipping_cost" gorm:"default:0"`
	OtherCost     float64 `json:"other_cost" gorm:"default:0"`
	GrandTotal    float64 `json:"grandtotal" gorm:"default:0"`
	AmountPaid    float64 `json:"amount_paid" gorm:"default:0"`
	AmountDue     float64 `json:"amount_due" gorm:"default:0"`

	// Payment
	PaymentTermID   string `json:"payment_term_id"`
	PaymentTermDays int    `json:"payment_term_days" gorm:"default:30"`
	PaymentStatus   string `json:"payment_status" gorm:"default:'unpaid'"` // unpaid, partial, paid, overdue

	// Accounting
	ExpenseAccountID string  `json:"expense_account_id" gorm:"index"`
	APAccountID      string  `json:"ap_account_id" gorm:"index"`    // Accounts Payable
	JournalEntryID   *string `json:"journal_entry_id" gorm:"index"` // Link to journal entry (future)

	// Additional Info
	Notes         string   `json:"notes" gorm:"type:text"`
	InternalNotes string   `json:"internal_notes" gorm:"type:text"`
	Terms         string   `json:"terms" gorm:"type:text"`
	Tags          []string `json:"tags" gorm:"serializer:json"`

	// Line Items
	Items []PurchaseInvoiceItem `json:"items" gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE"`

	// Metadata
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Workflow tracking
	ReceivedAt *time.Time `json:"received_at"`
	VerifiedAt *time.Time `json:"verified_at"`
	ApprovedAt *time.Time `json:"approved_at"`
	PaidAt     *time.Time `json:"paid_at"`
}

func (PurchaseInvoice) TableName() string {
	return "purchase_invoices"
}

// PurchaseInvoiceItem represents a line item in purchase invoice
type PurchaseInvoiceItem struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	InvoiceID    string    `json:"invoice_id" gorm:"index;not null"`
	TenantID     string    `json:"tenant_id" gorm:"index;not null"`
	OrderItemID  *string   `json:"order_item_id" gorm:"index"`
	LineNumber   int       `json:"line_number" gorm:"not null"`
	ProductID    string    `json:"product_id" gorm:"index;not null"`
	ProductCode  string    `json:"product_code" gorm:"not null"`
	ProductName  string    `json:"product_name" gorm:"not null"`
	ProductType  string    `json:"product_type"`
	Description  string    `json:"description" gorm:"type:text"`
	Quantity     float64   `json:"quantity" gorm:"not null"`
	UOM          string    `json:"uom" gorm:"not null"`
	UnitPrice    float64   `json:"unit_price" gorm:"not null"`
	DiscountType string    `json:"discount_type" gorm:"default:'percentage'"`
	Discount     float64   `json:"discount" gorm:"default:0"`
	TaxID        string    `json:"tax_id"`
	TaxRate      float64   `json:"tax_rate" gorm:"default:0"`
	TaxAmount    float64   `json:"tax_amount" gorm:"default:0"`
	Subtotal     float64   `json:"subtotal" gorm:"not null"`
	Total        float64   `json:"total" gorm:"not null"`
	AccountID    string    `json:"account_id" gorm:"index"` // COA for expense
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (PurchaseInvoiceItem) TableName() string {
	return "purchase_invoice_items"
}
