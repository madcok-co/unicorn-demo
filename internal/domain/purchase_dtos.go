package domain

import "time"

// ============================================================================
// PURCHASE REQUEST DTOs
// ============================================================================

type CreatePurchaseRequestDTO struct {
	Code          string                         `json:"code" validate:"required,max=50"`
	RequesterID   string                         `json:"requester_id" validate:"required"`
	DepartmentID  string                         `json:"department_id"`
	RequestDate   time.Time                      `json:"request_date" validate:"required"`
	RequiredDate  time.Time                      `json:"required_date" validate:"required"`
	Priority      string                         `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	ReferenceNo   string                         `json:"reference_no"`
	Purpose       string                         `json:"purpose"`
	ApproverID    string                         `json:"approver_id"`
	Notes         string                         `json:"notes"`
	InternalNotes string                         `json:"internal_notes"`
	Tags          []string                       `json:"tags"`
	Items         []CreatePurchaseRequestItemDTO `json:"items" validate:"required,min=1,dive"`
}

type CreatePurchaseRequestItemDTO struct {
	ProductID      string  `json:"product_id" validate:"required"`
	Description    string  `json:"description"`
	Quantity       float64 `json:"quantity" validate:"required,gt=0"`
	EstimatedPrice float64 `json:"estimated_price" validate:"min=0"`
	Notes          string  `json:"notes"`
}

type UpdatePurchaseRequestDTO struct {
	RequiredDate  *time.Time `json:"required_date"`
	Status        string     `json:"status" validate:"omitempty,oneof=draft submitted approved rejected cancelled completed"`
	Priority      string     `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	ReferenceNo   string     `json:"reference_no"`
	Purpose       string     `json:"purpose"`
	ApproverID    string     `json:"approver_id"`
	Notes         string     `json:"notes"`
	InternalNotes string     `json:"internal_notes"`
	Tags          []string   `json:"tags"`
}

type ListPurchaseRequestsDTO struct {
	Page         int        `json:"page" validate:"omitempty,min=1"`
	Limit        int        `json:"limit" validate:"omitempty,min=1,max=100"`
	Search       string     `json:"search"`
	Status       string     `json:"status" validate:"omitempty,oneof=draft submitted approved rejected cancelled completed"`
	RequesterID  string     `json:"requester_id"`
	ApproverID   string     `json:"approver_id"`
	DepartmentID string     `json:"department_id"`
	Priority     string     `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	DateFrom     *time.Time `json:"date_from"`
	DateTo       *time.Time `json:"date_to"`
	Active       *bool      `json:"active"`
}

// ============================================================================
// PURCHASE ORDER DTOs
// ============================================================================

type CreatePurchaseOrderDTO struct {
	Code            string                       `json:"code" validate:"required,max=50"`
	RequestID       *string                      `json:"request_id"`
	VendorID        string                       `json:"vendor_id" validate:"required"`
	OrderDate       time.Time                    `json:"order_date" validate:"required"`
	ExpectedDate    time.Time                    `json:"expected_date" validate:"required"`
	Priority        string                       `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	ReferenceNo     string                       `json:"reference_no"`
	VendorQuoteNo   string                       `json:"vendor_quote_no"`
	BuyerID         string                       `json:"buyer_id"`
	Currency        string                       `json:"currency" validate:"required"`
	ExchangeRate    float64                      `json:"exchange_rate" validate:"required,min=0"`
	PaymentTermID   string                       `json:"payment_term_id"`
	PaymentTermDays int                          `json:"payment_term_days"`
	WarehouseID     string                       `json:"warehouse_id"`
	ShippingMethod  string                       `json:"shipping_method"`
	ShippingAddress string                       `json:"shipping_address"`
	ShippingCity    string                       `json:"shipping_city"`
	ShippingState   string                       `json:"shipping_state"`
	ShippingZip     string                       `json:"shipping_zip"`
	ShippingCountry string                       `json:"shipping_country"`
	ShippingCost    float64                      `json:"shipping_cost"`
	OtherCost       float64                      `json:"other_cost"`
	Notes           string                       `json:"notes"`
	InternalNotes   string                       `json:"internal_notes"`
	Terms           string                       `json:"terms"`
	Tags            []string                     `json:"tags"`
	Items           []CreatePurchaseOrderItemDTO `json:"items" validate:"required,min=1,dive"`
}

type CreatePurchaseOrderItemDTO struct {
	ProductID    string  `json:"product_id" validate:"required"`
	Description  string  `json:"description"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	UnitPrice    float64 `json:"unit_price" validate:"required,min=0"`
	DiscountType string  `json:"discount_type" validate:"omitempty,oneof=percentage fixed"`
	Discount     float64 `json:"discount" validate:"min=0"`
	TaxID        string  `json:"tax_id"`
	Notes        string  `json:"notes"`
}

type UpdatePurchaseOrderDTO struct {
	ExpectedDate    *time.Time `json:"expected_date"`
	Status          string     `json:"status" validate:"omitempty,oneof=draft sent confirmed receiving completed cancelled"`
	Priority        string     `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	ReferenceNo     string     `json:"reference_no"`
	VendorQuoteNo   string     `json:"vendor_quote_no"`
	BuyerID         string     `json:"buyer_id"`
	PaymentTermID   string     `json:"payment_term_id"`
	PaymentTermDays *int       `json:"payment_term_days"`
	PaymentStatus   string     `json:"payment_status" validate:"omitempty,oneof=unpaid partial paid"`
	WarehouseID     string     `json:"warehouse_id"`
	ReceiptStatus   string     `json:"receipt_status" validate:"omitempty,oneof=pending partial completed"`
	ShippingMethod  string     `json:"shipping_method"`
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

type ListPurchaseOrdersDTO struct {
	Page          int        `json:"page" validate:"omitempty,min=1"`
	Limit         int        `json:"limit" validate:"omitempty,min=1,max=100"`
	Search        string     `json:"search"`
	VendorID      string     `json:"vendor_id"`
	Status        string     `json:"status" validate:"omitempty,oneof=draft sent confirmed receiving completed cancelled"`
	Priority      string     `json:"priority" validate:"omitempty,oneof=low normal high urgent"`
	PaymentStatus string     `json:"payment_status" validate:"omitempty,oneof=unpaid partial paid"`
	ReceiptStatus string     `json:"receipt_status" validate:"omitempty,oneof=pending partial completed"`
	WarehouseID   string     `json:"warehouse_id"`
	BuyerID       string     `json:"buyer_id"`
	DateFrom      *time.Time `json:"date_from"`
	DateTo        *time.Time `json:"date_to"`
	Active        *bool      `json:"active"`
}

// ============================================================================
// GOODS RECEIPT DTOs
// ============================================================================

type CreateGoodsReceiptDTO struct {
	Code          string                      `json:"code" validate:"required,max=50"`
	OrderID       string                      `json:"order_id" validate:"required"`
	ReceiptDate   time.Time                   `json:"receipt_date" validate:"required"`
	ScheduledDate time.Time                   `json:"scheduled_date" validate:"required"`
	ReferenceNo   string                      `json:"reference_no"`
	DeliveryNote  string                      `json:"delivery_note"`
	PackingList   string                      `json:"packing_list"`
	WarehouseID   string                      `json:"warehouse_id" validate:"required"`
	ReceiverID    string                      `json:"receiver_id"`
	InspectorID   string                      `json:"inspector_id"`
	Notes         string                      `json:"notes"`
	InternalNotes string                      `json:"internal_notes"`
	Tags          []string                    `json:"tags"`
	Items         []CreateGoodsReceiptItemDTO `json:"items" validate:"required,min=1,dive"`
}

type CreateGoodsReceiptItemDTO struct {
	OrderItemID   string   `json:"order_item_id" validate:"required"`
	ReceivedQty   float64  `json:"received_qty" validate:"required,gt=0"`
	AcceptedQty   float64  `json:"accepted_qty" validate:"required,min=0"`
	RejectedQty   float64  `json:"rejected_qty" validate:"min=0"`
	SerialNumbers []string `json:"serial_numbers"`
	BatchNumbers  []string `json:"batch_numbers"`
	ExpiryDates   []string `json:"expiry_dates"`
	Notes         string   `json:"notes"`
}

type UpdateGoodsReceiptDTO struct {
	ScheduledDate   *time.Time `json:"scheduled_date"`
	ActualDate      *time.Time `json:"actual_date"`
	Status          string     `json:"status" validate:"omitempty,oneof=draft received inspected accepted rejected"`
	QualityStatus   string     `json:"quality_status" validate:"omitempty,oneof=pending passed failed partial"`
	ReferenceNo     string     `json:"reference_no"`
	DeliveryNote    string     `json:"delivery_note"`
	PackingList     string     `json:"packing_list"`
	ReceiverID      string     `json:"receiver_id"`
	InspectorID     string     `json:"inspector_id"`
	RejectionReason string     `json:"rejection_reason"`
	Notes           string     `json:"notes"`
	InternalNotes   string     `json:"internal_notes"`
	Tags            []string   `json:"tags"`
	DeliveryProof   []string   `json:"delivery_proof"`
}

type ListGoodsReceiptsDTO struct {
	Page          int        `json:"page" validate:"omitempty,min=1"`
	Limit         int        `json:"limit" validate:"omitempty,min=1,max=100"`
	Search        string     `json:"search"`
	OrderID       string     `json:"order_id"`
	VendorID      string     `json:"vendor_id"`
	WarehouseID   string     `json:"warehouse_id"`
	Status        string     `json:"status" validate:"omitempty,oneof=draft received inspected accepted rejected"`
	QualityStatus string     `json:"quality_status" validate:"omitempty,oneof=pending passed failed partial"`
	DateFrom      *time.Time `json:"date_from"`
	DateTo        *time.Time `json:"date_to"`
	Active        *bool      `json:"active"`
}

// ============================================================================
// PURCHASE INVOICE DTOs
// ============================================================================

type CreatePurchaseInvoiceDTO struct {
	Code             string                         `json:"code" validate:"required,max=50"`
	OrderID          *string                        `json:"order_id"`
	ReceiptID        *string                        `json:"receipt_id"`
	VendorID         string                         `json:"vendor_id" validate:"required"`
	InvoiceDate      time.Time                      `json:"invoice_date" validate:"required"`
	DueDate          time.Time                      `json:"due_date" validate:"required"`
	ReferenceNo      string                         `json:"reference_no"`
	VendorInvoiceNo  string                         `json:"vendor_invoice_no"`
	TaxInvoiceNo     string                         `json:"tax_invoice_no"`
	Currency         string                         `json:"currency" validate:"required"`
	ExchangeRate     float64                        `json:"exchange_rate" validate:"required,min=0"`
	PaymentTermID    string                         `json:"payment_term_id"`
	PaymentTermDays  int                            `json:"payment_term_days"`
	ExpenseAccountID string                         `json:"expense_account_id"`
	APAccountID      string                         `json:"ap_account_id"`
	ShippingCost     float64                        `json:"shipping_cost"`
	OtherCost        float64                        `json:"other_cost"`
	Notes            string                         `json:"notes"`
	InternalNotes    string                         `json:"internal_notes"`
	Terms            string                         `json:"terms"`
	Tags             []string                       `json:"tags"`
	Items            []CreatePurchaseInvoiceItemDTO `json:"items" validate:"required,min=1,dive"`
}

type CreatePurchaseInvoiceItemDTO struct {
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

type UpdatePurchaseInvoiceDTO struct {
	DueDate          *time.Time `json:"due_date"`
	Status           string     `json:"status" validate:"omitempty,oneof=draft received verified approved partial_paid paid cancelled"`
	ReferenceNo      string     `json:"reference_no"`
	VendorInvoiceNo  string     `json:"vendor_invoice_no"`
	TaxInvoiceNo     string     `json:"tax_invoice_no"`
	PaymentTermID    string     `json:"payment_term_id"`
	PaymentTermDays  *int       `json:"payment_term_days"`
	PaymentStatus    string     `json:"payment_status" validate:"omitempty,oneof=unpaid partial paid overdue"`
	AmountPaid       *float64   `json:"amount_paid"`
	ExpenseAccountID string     `json:"expense_account_id"`
	APAccountID      string     `json:"ap_account_id"`
	ShippingCost     *float64   `json:"shipping_cost"`
	OtherCost        *float64   `json:"other_cost"`
	Notes            string     `json:"notes"`
	InternalNotes    string     `json:"internal_notes"`
	Terms            string     `json:"terms"`
	Tags             []string   `json:"tags"`
}

type ListPurchaseInvoicesDTO struct {
	Page          int        `json:"page" validate:"omitempty,min=1"`
	Limit         int        `json:"limit" validate:"omitempty,min=1,max=100"`
	Search        string     `json:"search"`
	VendorID      string     `json:"vendor_id"`
	OrderID       string     `json:"order_id"`
	Status        string     `json:"status" validate:"omitempty,oneof=draft received verified approved partial_paid paid cancelled"`
	PaymentStatus string     `json:"payment_status" validate:"omitempty,oneof=unpaid partial paid overdue"`
	DateFrom      *time.Time `json:"date_from"`
	DateTo        *time.Time `json:"date_to"`
	DueDateFrom   *time.Time `json:"due_date_from"`
	DueDateTo     *time.Time `json:"due_date_to"`
	Active        *bool      `json:"active"`
}
