package domain

import "time"

// ============================================================================
// SALES QUOTATION DTOs
// ============================================================================

type CreateSalesQuotationDTO struct {
	Code            string                        `json:"code" validate:"required,max=50"`
	CustomerID      string                        `json:"customer_id" validate:"required"`
	QuotationDate   time.Time                     `json:"quotation_date" validate:"required"`
	ValidUntil      time.Time                     `json:"valid_until" validate:"required"`
	ExpectedDate    *time.Time                    `json:"expected_date"`
	ReferenceNo     string                        `json:"reference_no"`
	SalespersonID   string                        `json:"salesperson_id"`
	Currency        string                        `json:"currency" validate:"required"`
	ExchangeRate    float64                       `json:"exchange_rate" validate:"required,min=0"`
	PaymentTermID   string                        `json:"payment_term_id"`
	PaymentTermDays int                           `json:"payment_term_days"`
	ShippingAddress string                        `json:"shipping_address"`
	ShippingCity    string                        `json:"shipping_city"`
	ShippingState   string                        `json:"shipping_state"`
	ShippingZip     string                        `json:"shipping_zip"`
	ShippingCountry string                        `json:"shipping_country"`
	ShippingCost    float64                       `json:"shipping_cost"`
	OtherCost       float64                       `json:"other_cost"`
	Notes           string                        `json:"notes"`
	InternalNotes   string                        `json:"internal_notes"`
	Terms           string                        `json:"terms"`
	Tags            []string                      `json:"tags"`
	Items           []CreateSalesQuotationItemDTO `json:"items" validate:"required,min=1,dive"`
}

type CreateSalesQuotationItemDTO struct {
	ProductID    string  `json:"product_id" validate:"required"`
	Description  string  `json:"description"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	UnitPrice    float64 `json:"unit_price" validate:"required,min=0"`
	DiscountType string  `json:"discount_type" validate:"omitempty,oneof=percentage fixed"`
	Discount     float64 `json:"discount" validate:"min=0"`
	TaxID        string  `json:"tax_id"`
	Notes        string  `json:"notes"`
}

type UpdateSalesQuotationDTO struct {
	ValidUntil      *time.Time `json:"valid_until"`
	ExpectedDate    *time.Time `json:"expected_date"`
	Status          string     `json:"status" validate:"omitempty,oneof=draft sent accepted rejected expired cancelled"`
	ReferenceNo     string     `json:"reference_no"`
	SalespersonID   string     `json:"salesperson_id"`
	PaymentTermID   string     `json:"payment_term_id"`
	PaymentTermDays *int       `json:"payment_term_days"`
	ShippingAddress string     `json:"shipping_address"`
	ShippingCity    string     `json:"shipping_city"`
	ShippingState   string     `json:"shipping_state"`
	ShippingZip     string     `json:"shipping_zip"`
	ShippingCountry string     `json:"shipping_country"`
	ShippingCost    *float64   `json:"shipping_cost"`
	OtherCost       *float64   `json:"other_cost"`
	Notes           string     `json:"notes"`
	InternalNotes   string     `json:"internal_notes"`
	Terms           string     `json:"terms"`
	Tags            []string   `json:"tags"`
}

type ListSalesQuotationsDTO struct {
	Page          int        `json:"page" validate:"omitempty,min=1"`
	Limit         int        `json:"limit" validate:"omitempty,min=1,max=100"`
	Search        string     `json:"search"`
	CustomerID    string     `json:"customer_id"`
	Status        string     `json:"status" validate:"omitempty,oneof=draft sent accepted rejected expired cancelled"`
	SalespersonID string     `json:"salesperson_id"`
	DateFrom      *time.Time `json:"date_from"`
	DateTo        *time.Time `json:"date_to"`
	Active        *bool      `json:"active"`
}

// ============================================================================
// SALES ORDER DTOs
// ============================================================================

type CreateSalesOrderDTO struct {
	Code            string                    `json:"code" validate:"required,max=50"`
	QuotationID     *string                   `json:"quotation_id"`
	CustomerID      string                    `json:"customer_id" validate:"required"`
	OrderDate       time.Time                 `json:"order_date" validate:"required"`
	ExpectedDate    time.Time                 `json:"expected_date" validate:"required"`
	Priority        string                    `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	ReferenceNo     string                    `json:"reference_no"`
	CustomerPO      string                    `json:"customer_po"`
	SalespersonID   string                    `json:"salesperson_id"`
	Currency        string                    `json:"currency" validate:"required"`
	ExchangeRate    float64                   `json:"exchange_rate" validate:"required,min=0"`
	PaymentTermID   string                    `json:"payment_term_id"`
	PaymentTermDays int                       `json:"payment_term_days"`
	WarehouseID     string                    `json:"warehouse_id"`
	ShippingMethod  string                    `json:"shipping_method"`
	ShippingAddress string                    `json:"shipping_address"`
	ShippingCity    string                    `json:"shipping_city"`
	ShippingState   string                    `json:"shipping_state"`
	ShippingZip     string                    `json:"shipping_zip"`
	ShippingCountry string                    `json:"shipping_country"`
	ShippingCost    float64                   `json:"shipping_cost"`
	OtherCost       float64                   `json:"other_cost"`
	Notes           string                    `json:"notes"`
	InternalNotes   string                    `json:"internal_notes"`
	Terms           string                    `json:"terms"`
	Tags            []string                  `json:"tags"`
	Items           []CreateSalesOrderItemDTO `json:"items" validate:"required,min=1,dive"`
}

type CreateSalesOrderItemDTO struct {
	ProductID    string  `json:"product_id" validate:"required"`
	Description  string  `json:"description"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	UnitPrice    float64 `json:"unit_price" validate:"required,min=0"`
	DiscountType string  `json:"discount_type" validate:"omitempty,oneof=percentage fixed"`
	Discount     float64 `json:"discount" validate:"min=0"`
	TaxID        string  `json:"tax_id"`
	Notes        string  `json:"notes"`
}

type UpdateSalesOrderDTO struct {
	ExpectedDate      *time.Time `json:"expected_date"`
	Status            string     `json:"status" validate:"omitempty,oneof=draft confirmed processing completed cancelled"`
	Priority          string     `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	ReferenceNo       string     `json:"reference_no"`
	CustomerPO        string     `json:"customer_po"`
	SalespersonID     string     `json:"salesperson_id"`
	PaymentTermID     string     `json:"payment_term_id"`
	PaymentTermDays   *int       `json:"payment_term_days"`
	PaymentStatus     string     `json:"payment_status" validate:"omitempty,oneof=unpaid partial paid"`
	WarehouseID       string     `json:"warehouse_id"`
	FulfillmentStatus string     `json:"fulfillment_status" validate:"omitempty,oneof=pending partial completed"`
	ShippingMethod    string     `json:"shipping_method"`
	ShippingAddress   string     `json:"shipping_address"`
	ShippingCity      string     `json:"shipping_city"`
	ShippingState     string     `json:"shipping_state"`
	ShippingZip       string     `json:"shipping_zip"`
	ShippingCountry   string     `json:"shipping_country"`
	ShippingCost      *float64   `json:"shipping_cost"`
	OtherCost         *float64   `json:"other_cost"`
	Notes             string     `json:"notes"`
	InternalNotes     string     `json:"internal_notes"`
	Terms             string     `json:"terms"`
	Tags              []string   `json:"tags"`
}

type ListSalesOrdersDTO struct {
	Page              int        `json:"page" validate:"omitempty,min=1"`
	Limit             int        `json:"limit" validate:"omitempty,min=1,max=100"`
	Search            string     `json:"search"`
	CustomerID        string     `json:"customer_id"`
	Status            string     `json:"status" validate:"omitempty,oneof=draft confirmed processing completed cancelled"`
	Priority          string     `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	PaymentStatus     string     `json:"payment_status" validate:"omitempty,oneof=unpaid partial paid"`
	FulfillmentStatus string     `json:"fulfillment_status" validate:"omitempty,oneof=pending partial completed"`
	SalespersonID     string     `json:"salesperson_id"`
	WarehouseID       string     `json:"warehouse_id"`
	DateFrom          *time.Time `json:"date_from"`
	DateTo            *time.Time `json:"date_to"`
	Active            *bool      `json:"active"`
}

// ============================================================================
// DELIVERY ORDER DTOs
// ============================================================================

type CreateDeliveryOrderDTO struct {
	Code            string                       `json:"code" validate:"required,max=50"`
	OrderID         string                       `json:"order_id" validate:"required"`
	DeliveryDate    time.Time                    `json:"delivery_date" validate:"required"`
	ScheduledDate   time.Time                    `json:"scheduled_date" validate:"required"`
	ReferenceNo     string                       `json:"reference_no"`
	TrackingNumber  string                       `json:"tracking_number"`
	WarehouseID     string                       `json:"warehouse_id" validate:"required"`
	ShippingMethod  string                       `json:"shipping_method"`
	ShippingCost    float64                      `json:"shipping_cost"`
	CourierName     string                       `json:"courier_name"`
	DriverName      string                       `json:"driver_name"`
	VehicleNumber   string                       `json:"vehicle_number"`
	ShippingAddress string                       `json:"shipping_address"`
	ShippingCity    string                       `json:"shipping_city"`
	ShippingState   string                       `json:"shipping_state"`
	ShippingZip     string                       `json:"shipping_zip"`
	ShippingCountry string                       `json:"shipping_country"`
	RecipientName   string                       `json:"recipient_name"`
	RecipientPhone  string                       `json:"recipient_phone"`
	RecipientEmail  string                       `json:"recipient_email"`
	Notes           string                       `json:"notes"`
	InternalNotes   string                       `json:"internal_notes"`
	Tags            []string                     `json:"tags"`
	Items           []CreateDeliveryOrderItemDTO `json:"items" validate:"required,min=1,dive"`
}

type CreateDeliveryOrderItemDTO struct {
	OrderItemID   string   `json:"order_item_id" validate:"required"`
	ProductID     string   `json:"product_id" validate:"required"`
	DeliveredQty  float64  `json:"delivered_qty" validate:"required,gt=0"`
	SerialNumbers []string `json:"serial_numbers"`
	BatchNumbers  []string `json:"batch_numbers"`
	Notes         string   `json:"notes"`
}

type UpdateDeliveryOrderDTO struct {
	ScheduledDate     *time.Time `json:"scheduled_date"`
	ActualDate        *time.Time `json:"actual_date"`
	Status            string     `json:"status" validate:"omitempty,oneof=draft ready in_transit delivered cancelled"`
	ReferenceNo       string     `json:"reference_no"`
	TrackingNumber    string     `json:"tracking_number"`
	ShippingMethod    string     `json:"shipping_method"`
	ShippingCost      *float64   `json:"shipping_cost"`
	CourierName       string     `json:"courier_name"`
	DriverName        string     `json:"driver_name"`
	VehicleNumber     string     `json:"vehicle_number"`
	RecipientName     string     `json:"recipient_name"`
	RecipientPhone    string     `json:"recipient_phone"`
	RecipientEmail    string     `json:"recipient_email"`
	Notes             string     `json:"notes"`
	InternalNotes     string     `json:"internal_notes"`
	Tags              []string   `json:"tags"`
	ReceivedBy        string     `json:"received_by"`
	ReceiverSignature string     `json:"receiver_signature"`
	DeliveryProof     []string   `json:"delivery_proof"`
}

type ListDeliveryOrdersDTO struct {
	Page        int        `json:"page" validate:"omitempty,min=1"`
	Limit       int        `json:"limit" validate:"omitempty,min=1,max=100"`
	Search      string     `json:"search"`
	OrderID     string     `json:"order_id"`
	CustomerID  string     `json:"customer_id"`
	WarehouseID string     `json:"warehouse_id"`
	Status      string     `json:"status" validate:"omitempty,oneof=draft ready in_transit delivered cancelled"`
	DateFrom    *time.Time `json:"date_from"`
	DateTo      *time.Time `json:"date_to"`
	Active      *bool      `json:"active"`
}

// ============================================================================
// SALES INVOICE DTOs
// ============================================================================

type CreateSalesInvoiceDTO struct {
	Code            string                      `json:"code" validate:"required,max=50"`
	OrderID         *string                     `json:"order_id"`
	DeliveryOrderID *string                     `json:"delivery_order_id"`
	CustomerID      string                      `json:"customer_id" validate:"required"`
	InvoiceDate     time.Time                   `json:"invoice_date" validate:"required"`
	DueDate         time.Time                   `json:"due_date" validate:"required"`
	ReferenceNo     string                      `json:"reference_no"`
	CustomerPO      string                      `json:"customer_po"`
	TaxInvoiceNo    string                      `json:"tax_invoice_no"`
	Currency        string                      `json:"currency" validate:"required"`
	ExchangeRate    float64                     `json:"exchange_rate" validate:"required,min=0"`
	PaymentTermID   string                      `json:"payment_term_id"`
	PaymentTermDays int                         `json:"payment_term_days"`
	IncomeAccountID string                      `json:"income_account_id"`
	ARAccountID     string                      `json:"ar_account_id"`
	BillingAddress  string                      `json:"billing_address"`
	BillingCity     string                      `json:"billing_city"`
	BillingState    string                      `json:"billing_state"`
	BillingZip      string                      `json:"billing_zip"`
	BillingCountry  string                      `json:"billing_country"`
	ShippingCost    float64                     `json:"shipping_cost"`
	OtherCost       float64                     `json:"other_cost"`
	Notes           string                      `json:"notes"`
	InternalNotes   string                      `json:"internal_notes"`
	Terms           string                      `json:"terms"`
	Tags            []string                    `json:"tags"`
	Items           []CreateSalesInvoiceItemDTO `json:"items" validate:"required,min=1,dive"`
}

type CreateSalesInvoiceItemDTO struct {
	OrderItemID  *string `json:"order_item_id"`
	ProductID    string  `json:"product_id" validate:"required"`
	Description  string  `json:"description"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	UnitPrice    float64 `json:"unit_price" validate:"required,min=0"`
	DiscountType string  `json:"discount_type" validate:"omitempty,oneof=percentage fixed"`
	Discount     float64 `json:"discount" validate:"min=0"`
	TaxID        string  `json:"tax_id"`
	AccountID    string  `json:"account_id"`
	Notes        string  `json:"notes"`
}

type UpdateSalesInvoiceDTO struct {
	DueDate         *time.Time `json:"due_date"`
	Status          string     `json:"status" validate:"omitempty,oneof=draft sent partial_paid paid overdue cancelled"`
	ReferenceNo     string     `json:"reference_no"`
	CustomerPO      string     `json:"customer_po"`
	TaxInvoiceNo    string     `json:"tax_invoice_no"`
	PaymentTermID   string     `json:"payment_term_id"`
	PaymentTermDays *int       `json:"payment_term_days"`
	PaymentStatus   string     `json:"payment_status" validate:"omitempty,oneof=unpaid partial paid overdue"`
	AmountPaid      *float64   `json:"amount_paid"`
	IncomeAccountID string     `json:"income_account_id"`
	ARAccountID     string     `json:"ar_account_id"`
	BillingAddress  string     `json:"billing_address"`
	BillingCity     string     `json:"billing_city"`
	BillingState    string     `json:"billing_state"`
	BillingZip      string     `json:"billing_zip"`
	BillingCountry  string     `json:"billing_country"`
	ShippingCost    *float64   `json:"shipping_cost"`
	OtherCost       *float64   `json:"other_cost"`
	Notes           string     `json:"notes"`
	InternalNotes   string     `json:"internal_notes"`
	Terms           string     `json:"terms"`
	Tags            []string   `json:"tags"`
}

type ListSalesInvoicesDTO struct {
	Page          int        `json:"page" validate:"omitempty,min=1"`
	Limit         int        `json:"limit" validate:"omitempty,min=1,max=100"`
	Search        string     `json:"search"`
	CustomerID    string     `json:"customer_id"`
	OrderID       string     `json:"order_id"`
	Status        string     `json:"status" validate:"omitempty,oneof=draft sent partial_paid paid overdue cancelled"`
	PaymentStatus string     `json:"payment_status" validate:"omitempty,oneof=unpaid partial paid overdue"`
	DateFrom      *time.Time `json:"date_from"`
	DateTo        *time.Time `json:"date_to"`
	DueDateFrom   *time.Time `json:"due_date_from"`
	DueDateTo     *time.Time `json:"due_date_to"`
	Active        *bool      `json:"active"`
}

// ============================================================================
// WORKFLOW ACTION DTOs
// ============================================================================

type ConvertQuotationToOrderDTO struct {
	QuotationID string `json:"quotation_id" validate:"required"`
	OrderCode   string `json:"order_code" validate:"required"`
	WarehouseID string `json:"warehouse_id"`
}

type CreateDeliveryFromOrderDTO struct {
	OrderID       string    `json:"order_id" validate:"required"`
	DeliveryCode  string    `json:"delivery_code" validate:"required"`
	ScheduledDate time.Time `json:"scheduled_date" validate:"required"`
	Items         []struct {
		OrderItemID  string  `json:"order_item_id" validate:"required"`
		DeliveredQty float64 `json:"delivered_qty" validate:"required,gt=0"`
	} `json:"items" validate:"required,min=1,dive"`
}

type CreateInvoiceFromOrderDTO struct {
	OrderID     string    `json:"order_id" validate:"required"`
	InvoiceCode string    `json:"invoice_code" validate:"required"`
	InvoiceDate time.Time `json:"invoice_date" validate:"required"`
	DueDate     time.Time `json:"due_date" validate:"required"`
}

type RecordPaymentDTO struct {
	InvoiceID     string    `json:"invoice_id" validate:"required"`
	PaymentDate   time.Time `json:"payment_date" validate:"required"`
	Amount        float64   `json:"amount" validate:"required,gt=0"`
	PaymentMethod string    `json:"payment_method" validate:"required"`
	ReferenceNo   string    `json:"reference_no"`
	Notes         string    `json:"notes"`
}
