package domain

import "time"

// ============================================================================
// SALES MODULE MODELS
// ============================================================================

// SalesQuotation represents a sales quotation/quote
type SalesQuotation struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_sales_quotation_code_tenant;not null"` // SQ-2024-001
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_sales_quotation_code_tenant;index;not null"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Customer Information
	CustomerID   string `json:"customer_id" gorm:"index;not null"`
	CustomerName string `json:"customer_name" gorm:"not null"`
	CustomerType string `json:"customer_type"`

	// Document Details
	QuotationDate   time.Time  `json:"quotation_date" gorm:"not null"`
	ValidUntil      time.Time  `json:"valid_until" gorm:"not null"`
	ExpectedDate    *time.Time `json:"expected_date"`
	Status          string     `json:"status" gorm:"default:'draft'"` // draft, sent, accepted, rejected, expired, cancelled
	ReferenceNo     string     `json:"reference_no"`
	SalespersonID   string     `json:"salesperson_id" gorm:"index"`
	SalespersonName string     `json:"salesperson_name"`

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

	// Shipping
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
	Items []SalesQuotationItem `json:"items" gorm:"foreignKey:QuotationID;constraint:OnDelete:CASCADE"`

	// Metadata
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Workflow tracking
	SentAt     *time.Time `json:"sent_at"`
	AcceptedAt *time.Time `json:"accepted_at"`
	RejectedAt *time.Time `json:"rejected_at"`
	ExpiredAt  *time.Time `json:"expired_at"`
}

func (SalesQuotation) TableName() string {
	return "sales_quotations"
}

// SalesQuotationItem represents a line item in sales quotation
type SalesQuotationItem struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	QuotationID  string    `json:"quotation_id" gorm:"index;not null"`
	TenantID     string    `json:"tenant_id" gorm:"index;not null"`
	LineNumber   int       `json:"line_number" gorm:"not null"`
	ProductID    string    `json:"product_id" gorm:"index;not null"`
	ProductCode  string    `json:"product_code" gorm:"not null"`
	ProductName  string    `json:"product_name" gorm:"not null"`
	ProductType  string    `json:"product_type"` // product, service, consumable
	Description  string    `json:"description" gorm:"type:text"`
	Quantity     float64   `json:"quantity" gorm:"not null"`
	UOM          string    `json:"uom" gorm:"not null"`
	UnitPrice    float64   `json:"unit_price" gorm:"not null"`
	DiscountType string    `json:"discount_type" gorm:"default:'percentage'"` // percentage, fixed
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

func (SalesQuotationItem) TableName() string {
	return "sales_quotation_items"
}

// SalesOrder represents a confirmed sales order
type SalesOrder struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_sales_order_code_tenant;not null"` // SO-2024-001
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_sales_order_code_tenant;index;not null"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Reference
	QuotationID   *string `json:"quotation_id" gorm:"index"`
	QuotationCode *string `json:"quotation_code"`

	// Customer Information
	CustomerID   string `json:"customer_id" gorm:"index;not null"`
	CustomerName string `json:"customer_name" gorm:"not null"`
	CustomerType string `json:"customer_type"`

	// Document Details
	OrderDate       time.Time  `json:"order_date" gorm:"not null"`
	ExpectedDate    time.Time  `json:"expected_date" gorm:"not null"`
	ConfirmedDate   *time.Time `json:"confirmed_date"`
	Status          string     `json:"status" gorm:"default:'draft'"`    // draft, confirmed, processing, completed, cancelled
	Priority        string     `json:"priority" gorm:"default:'normal'"` // low, normal, high, urgent
	ReferenceNo     string     `json:"reference_no"`
	CustomerPO      string     `json:"customer_po"` // Customer Purchase Order number
	SalespersonID   string     `json:"salesperson_id" gorm:"index"`
	SalespersonName string     `json:"salesperson_name"`

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

	// Fulfillment
	WarehouseID       string  `json:"warehouse_id" gorm:"index"`
	WarehouseName     string  `json:"warehouse_name"`
	FulfillmentStatus string  `json:"fulfillment_status" gorm:"default:'pending'"` // pending, partial, completed
	DeliveredQty      float64 `json:"delivered_qty" gorm:"default:0"`
	InvoicedQty       float64 `json:"invoiced_qty" gorm:"default:0"`

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
	Items []SalesOrderItem `json:"items" gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`

	// Metadata
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SalesOrder) TableName() string {
	return "sales_orders"
}

// SalesOrderItem represents a line item in sales order
type SalesOrderItem struct {
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
	DeliveredQty float64   `json:"delivered_qty" gorm:"default:0"`
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

func (SalesOrderItem) TableName() string {
	return "sales_order_items"
}

// DeliveryOrder represents a delivery/shipment order
type DeliveryOrder struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_delivery_order_code_tenant;not null"` // DO-2024-001
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_delivery_order_code_tenant;index;not null"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Reference
	OrderID   string `json:"order_id" gorm:"index;not null"`
	OrderCode string `json:"order_code" gorm:"not null"`

	// Customer Information
	CustomerID   string `json:"customer_id" gorm:"index;not null"`
	CustomerName string `json:"customer_name" gorm:"not null"`

	// Document Details
	DeliveryDate   time.Time  `json:"delivery_date" gorm:"not null"`
	ScheduledDate  time.Time  `json:"scheduled_date" gorm:"not null"`
	ActualDate     *time.Time `json:"actual_date"`
	Status         string     `json:"status" gorm:"default:'draft'"` // draft, ready, in_transit, delivered, cancelled
	ReferenceNo    string     `json:"reference_no"`
	TrackingNumber string     `json:"tracking_number"`

	// Warehouse & Shipping
	WarehouseID    string  `json:"warehouse_id" gorm:"index;not null"`
	WarehouseName  string  `json:"warehouse_name" gorm:"not null"`
	ShippingMethod string  `json:"shipping_method"`
	ShippingCost   float64 `json:"shipping_cost" gorm:"default:0"`
	CourierName    string  `json:"courier_name"`
	DriverName     string  `json:"driver_name"`
	VehicleNumber  string  `json:"vehicle_number"`

	// Shipping Address
	ShippingAddress string `json:"shipping_address"`
	ShippingCity    string `json:"shipping_city"`
	ShippingState   string `json:"shipping_state"`
	ShippingZip     string `json:"shipping_zip"`
	ShippingCountry string `json:"shipping_country"`

	// Recipient
	RecipientName  string `json:"recipient_name"`
	RecipientPhone string `json:"recipient_phone"`
	RecipientEmail string `json:"recipient_email"`

	// Additional Info
	Notes         string   `json:"notes" gorm:"type:text"`
	InternalNotes string   `json:"internal_notes" gorm:"type:text"`
	Tags          []string `json:"tags" gorm:"serializer:json"`

	// Line Items
	Items []DeliveryOrderItem `json:"items" gorm:"foreignKey:DeliveryOrderID;constraint:OnDelete:CASCADE"`

	// Proof of Delivery
	ReceivedBy        string     `json:"received_by"`
	ReceivedAt        *time.Time `json:"received_at"`
	ReceiverSignature string     `json:"receiver_signature"`                    // URL to signature image
	DeliveryProof     []string   `json:"delivery_proof" gorm:"serializer:json"` // URLs to proof photos

	// Metadata
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DeliveryOrder) TableName() string {
	return "delivery_orders"
}

// DeliveryOrderItem represents a line item in delivery order
type DeliveryOrderItem struct {
	ID              string    `json:"id" gorm:"primaryKey"`
	DeliveryOrderID string    `json:"delivery_order_id" gorm:"index;not null"`
	TenantID        string    `json:"tenant_id" gorm:"index;not null"`
	OrderItemID     string    `json:"order_item_id" gorm:"index;not null"`
	LineNumber      int       `json:"line_number" gorm:"not null"`
	ProductID       string    `json:"product_id" gorm:"index;not null"`
	ProductCode     string    `json:"product_code" gorm:"not null"`
	ProductName     string    `json:"product_name" gorm:"not null"`
	Description     string    `json:"description"`
	OrderedQty      float64   `json:"ordered_qty" gorm:"not null"`
	DeliveredQty    float64   `json:"delivered_qty" gorm:"not null"`
	UOM             string    `json:"uom" gorm:"not null"`
	SerialNumbers   []string  `json:"serial_numbers" gorm:"serializer:json"`
	BatchNumbers    []string  `json:"batch_numbers" gorm:"serializer:json"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (DeliveryOrderItem) TableName() string {
	return "delivery_order_items"
}

// SalesInvoice represents a sales invoice
type SalesInvoice struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_sales_invoice_code_tenant;not null"` // INV-2024-001
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_sales_invoice_code_tenant;index;not null"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Reference
	OrderID           *string `json:"order_id" gorm:"index"`
	OrderCode         *string `json:"order_code"`
	DeliveryOrderID   *string `json:"delivery_order_id" gorm:"index"`
	DeliveryOrderCode *string `json:"delivery_order_code"`

	// Customer Information
	CustomerID    string `json:"customer_id" gorm:"index;not null"`
	CustomerName  string `json:"customer_name" gorm:"not null"`
	CustomerType  string `json:"customer_type"`
	CustomerTaxID string `json:"customer_tax_id"`

	// Document Details
	InvoiceDate  time.Time `json:"invoice_date" gorm:"not null"`
	DueDate      time.Time `json:"due_date" gorm:"not null"`
	Status       string    `json:"status" gorm:"default:'draft'"` // draft, sent, partial_paid, paid, overdue, cancelled
	ReferenceNo  string    `json:"reference_no"`
	CustomerPO   string    `json:"customer_po"`
	TaxInvoiceNo string    `json:"tax_invoice_no"` // Faktur Pajak number

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
	IncomeAccountID string  `json:"income_account_id" gorm:"index"`
	ARAccountID     string  `json:"ar_account_id" gorm:"index"`    // Accounts Receivable
	JournalEntryID  *string `json:"journal_entry_id" gorm:"index"` // Link to journal entry (future)

	// Billing Address
	BillingAddress string `json:"billing_address"`
	BillingCity    string `json:"billing_city"`
	BillingState   string `json:"billing_state"`
	BillingZip     string `json:"billing_zip"`
	BillingCountry string `json:"billing_country"`

	// Additional Info
	Notes         string   `json:"notes" gorm:"type:text"`
	InternalNotes string   `json:"internal_notes" gorm:"type:text"`
	Terms         string   `json:"terms" gorm:"type:text"`
	Tags          []string `json:"tags" gorm:"serializer:json"`

	// Line Items
	Items []SalesInvoiceItem `json:"items" gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE"`

	// Metadata
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Workflow tracking
	SentAt *time.Time `json:"sent_at"`
	PaidAt *time.Time `json:"paid_at"`
}

func (SalesInvoice) TableName() string {
	return "sales_invoices"
}

// SalesInvoiceItem represents a line item in sales invoice
type SalesInvoiceItem struct {
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
	AccountID    string    `json:"account_id" gorm:"index"` // COA for revenue
	Notes        string    `json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (SalesInvoiceItem) TableName() string {
	return "sales_invoice_items"
}
