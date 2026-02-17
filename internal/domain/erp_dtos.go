package domain

// ============================================================================
// CUSTOMER DTOs
// ============================================================================

type CreateCustomerDTO struct {
	Code   string `json:"code" validate:"required,max=50"`
	Name   string `json:"name" validate:"required,max=200"`
	Type   string `json:"type" validate:"required,oneof=individual company"`
	Email  string `json:"email" validate:"omitempty,email"`
	Phone  string `json:"phone" validate:"omitempty,max=50"`
	Mobile string `json:"mobile" validate:"omitempty,max=50"`
	TaxID  string `json:"tax_id" validate:"omitempty,max:50"`

	BillingAddress string `json:"billing_address"`
	BillingCity    string `json:"billing_city"`
	BillingState   string `json:"billing_state"`
	BillingZip     string `json:"billing_zip"`
	BillingCountry string `json:"billing_country"`

	ShippingAddress string `json:"shipping_address"`
	ShippingCity    string `json:"shipping_city"`
	ShippingState   string `json:"shipping_state"`
	ShippingZip     string `json:"shipping_zip"`
	ShippingCountry string `json:"shipping_country"`

	CreditLimit     float64 `json:"credit_limit"`
	PaymentTermDays int     `json:"payment_term_days"`
	Currency        string  `json:"currency"`

	ContactPerson    string   `json:"contact_person"`
	ContactPersonJob string   `json:"contact_person_job"`
	Notes            string   `json:"notes"`
	Tags             []string `json:"tags"`
}

type UpdateCustomerDTO struct {
	Name   string `json:"name" validate:"omitempty,max=200"`
	Type   string `json:"type" validate:"omitempty,oneof=individual company"`
	Email  string `json:"email" validate:"omitempty,email"`
	Phone  string `json:"phone"`
	Mobile string `json:"mobile"`
	TaxID  string `json:"tax_id"`
	Active *bool  `json:"active"`

	BillingAddress string `json:"billing_address"`
	BillingCity    string `json:"billing_city"`
	BillingState   string `json:"billing_state"`
	BillingZip     string `json:"billing_zip"`
	BillingCountry string `json:"billing_country"`

	ShippingAddress string `json:"shipping_address"`
	ShippingCity    string `json:"shipping_city"`
	ShippingState   string `json:"shipping_state"`
	ShippingZip     string `json:"shipping_zip"`
	ShippingCountry string `json:"shipping_country"`

	CreditLimit     *float64 `json:"credit_limit"`
	PaymentTermDays *int     `json:"payment_term_days"`
	Currency        string   `json:"currency"`

	ContactPerson    string   `json:"contact_person"`
	ContactPersonJob string   `json:"contact_person_job"`
	Notes            string   `json:"notes"`
	Tags             []string `json:"tags"`
}

type ListCustomersDTO struct {
	Page   int    `json:"page" validate:"omitempty,min=1"`
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Search string `json:"search"`
	Type   string `json:"type" validate:"omitempty,oneof=individual company"`
	Active *bool  `json:"active"`
}

// ============================================================================
// VENDOR DTOs
// ============================================================================

type CreateVendorDTO struct {
	Code   string `json:"code" validate:"required,max=50"`
	Name   string `json:"name" validate:"required,max=200"`
	Type   string `json:"type" validate:"required,oneof=individual company"`
	Email  string `json:"email" validate:"omitempty,email"`
	Phone  string `json:"phone" validate:"omitempty,max=50"`
	Mobile string `json:"mobile" validate:"omitempty,max=50"`
	TaxID  string `json:"tax_id" validate:"omitempty,max=50"`

	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`

	PaymentTermDays int    `json:"payment_term_days"`
	Currency        string `json:"currency"`
	BankName        string `json:"bank_name"`
	BankAccount     string `json:"bank_account"`
	BankAccountName string `json:"bank_account_name"`

	ContactPerson    string   `json:"contact_person"`
	ContactPersonJob string   `json:"contact_person_job"`
	Rating           int      `json:"rating" validate:"omitempty,min=0,max=5"`
	Category         string   `json:"category"`
	Notes            string   `json:"notes"`
	Tags             []string `json:"tags"`
}

type UpdateVendorDTO struct {
	Name   string `json:"name" validate:"omitempty,max=200"`
	Type   string `json:"type" validate:"omitempty,oneof=individual company"`
	Email  string `json:"email" validate:"omitempty,email"`
	Phone  string `json:"phone"`
	Mobile string `json:"mobile"`
	TaxID  string `json:"tax_id"`
	Active *bool  `json:"active"`

	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`

	PaymentTermDays *int   `json:"payment_term_days"`
	Currency        string `json:"currency"`
	BankName        string `json:"bank_name"`
	BankAccount     string `json:"bank_account"`
	BankAccountName string `json:"bank_account_name"`

	ContactPerson    string   `json:"contact_person"`
	ContactPersonJob string   `json:"contact_person_job"`
	Rating           *int     `json:"rating" validate:"omitempty,min=0,max=5"`
	Category         string   `json:"category"`
	Notes            string   `json:"notes"`
	Tags             []string `json:"tags"`
}

type ListVendorsDTO struct {
	Page     int    `json:"page" validate:"omitempty,min=1"`
	Limit    int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Search   string `json:"search"`
	Type     string `json:"type" validate:"omitempty,oneof=individual company"`
	Active   *bool  `json:"active"`
	Category string `json:"category"`
}

// ============================================================================
// PRODUCT DTOs
// ============================================================================

type CreateProductDTO struct {
	Code     string `json:"code" validate:"required,max=50"`
	Name     string `json:"name" validate:"required,max=200"`
	Type     string `json:"type" validate:"required,oneof=product service consumable"`
	Category string `json:"category"`
	UOM      string `json:"uom" validate:"required,max=50"`
	Barcode  string `json:"barcode"`
	SKU      string `json:"sku"`

	CanBeSold      bool    `json:"can_be_sold"`
	CanBePurchased bool    `json:"can_be_purchased"`
	TrackInventory bool    `json:"track_inventory"`
	MinStock       float64 `json:"min_stock"`
	MaxStock       float64 `json:"max_stock"`
	ReorderLevel   float64 `json:"reorder_level"`

	SalePrice     float64 `json:"sale_price"`
	PurchasePrice float64 `json:"purchase_price"`
	Cost          float64 `json:"cost"`
	Currency      string  `json:"currency"`

	TaxCategory string  `json:"tax_category"`
	TaxRate     float64 `json:"tax_rate"`

	Weight       float64 `json:"weight"`
	Volume       float64 `json:"volume"`
	Length       float64 `json:"length"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
	WarrantyDays int     `json:"warranty_days"`

	Description      string   `json:"description"`
	InternalNotes    string   `json:"internal_notes"`
	ImageURL         string   `json:"image_url"`
	AdditionalImages []string `json:"additional_images"`

	IncomeAccountID  string `json:"income_account_id"`
	ExpenseAccountID string `json:"expense_account_id"`
	AssetAccountID   string `json:"asset_account_id"`

	Tags []string `json:"tags"`
}

type UpdateProductDTO struct {
	Name     string `json:"name" validate:"omitempty,max=200"`
	Type     string `json:"type" validate:"omitempty,oneof=product service consumable"`
	Category string `json:"category"`
	UOM      string `json:"uom"`
	Barcode  string `json:"barcode"`
	SKU      string `json:"sku"`
	Active   *bool  `json:"active"`

	CanBeSold      *bool    `json:"can_be_sold"`
	CanBePurchased *bool    `json:"can_be_purchased"`
	TrackInventory *bool    `json:"track_inventory"`
	MinStock       *float64 `json:"min_stock"`
	MaxStock       *float64 `json:"max_stock"`
	ReorderLevel   *float64 `json:"reorder_level"`

	SalePrice     *float64 `json:"sale_price"`
	PurchasePrice *float64 `json:"purchase_price"`
	Cost          *float64 `json:"cost"`
	Currency      string   `json:"currency"`

	TaxCategory string   `json:"tax_category"`
	TaxRate     *float64 `json:"tax_rate"`

	Weight       *float64 `json:"weight"`
	Volume       *float64 `json:"volume"`
	Length       *float64 `json:"length"`
	Width        *float64 `json:"width"`
	Height       *float64 `json:"height"`
	WarrantyDays *int     `json:"warranty_days"`

	Description      string   `json:"description"`
	InternalNotes    string   `json:"internal_notes"`
	ImageURL         string   `json:"image_url"`
	AdditionalImages []string `json:"additional_images"`

	IncomeAccountID  string `json:"income_account_id"`
	ExpenseAccountID string `json:"expense_account_id"`
	AssetAccountID   string `json:"asset_account_id"`

	Tags []string `json:"tags"`
}

type ListProductsDTO struct {
	Page           int    `json:"page" validate:"omitempty,min=1"`
	Limit          int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Search         string `json:"search"`
	Type           string `json:"type" validate:"omitempty,oneof=product service consumable"`
	Category       string `json:"category"`
	Active         *bool  `json:"active"`
	CanBeSold      *bool  `json:"can_be_sold"`
	CanBePurchased *bool  `json:"can_be_purchased"`
	TrackInventory *bool  `json:"track_inventory"`
}

// ============================================================================
// WAREHOUSE DTOs
// ============================================================================

type CreateWarehouseDTO struct {
	Code string `json:"code" validate:"required,max=50"`
	Name string `json:"name" validate:"required,max=200"`
	Type string `json:"type" validate:"required,oneof=physical virtual transit"`

	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`

	Phone       string `json:"phone"`
	Email       string `json:"email" validate:"omitempty,email"`
	ManagerID   string `json:"manager_id"`
	ManagerName string `json:"manager_name"`

	TotalArea     float64 `json:"total_area"`
	TotalCapacity float64 `json:"total_capacity"`

	AllowNegativeStock bool   `json:"allow_negative_stock"`
	IsDefault          bool   `json:"is_default"`
	Description        string `json:"description"`
	Notes              string `json:"notes"`
}

type UpdateWarehouseDTO struct {
	Name   string `json:"name" validate:"omitempty,max=200"`
	Type   string `json:"type" validate:"omitempty,oneof=physical virtual transit"`
	Active *bool  `json:"active"`

	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`

	Phone       string `json:"phone"`
	Email       string `json:"email" validate:"omitempty,email"`
	ManagerID   string `json:"manager_id"`
	ManagerName string `json:"manager_name"`

	TotalArea     *float64 `json:"total_area"`
	TotalCapacity *float64 `json:"total_capacity"`

	AllowNegativeStock *bool  `json:"allow_negative_stock"`
	IsDefault          *bool  `json:"is_default"`
	Description        string `json:"description"`
	Notes              string `json:"notes"`
}

type ListWarehousesDTO struct {
	Page   int    `json:"page" validate:"omitempty,min=1"`
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Search string `json:"search"`
	Type   string `json:"type" validate:"omitempty,oneof=physical virtual transit"`
	Active *bool  `json:"active"`
}

// ============================================================================
// CHART OF ACCOUNT DTOs
// ============================================================================

type CreateChartOfAccountDTO struct {
	Code     string `json:"code" validate:"required,max=50"`
	Name     string `json:"name" validate:"required,max=200"`
	Type     string `json:"type" validate:"required,oneof=asset liability equity income expense"`
	Category string `json:"category"`
	ParentID string `json:"parent_id"`
	IsGroup  bool   `json:"is_group"`

	Currency            string `json:"currency"`
	AllowReconciliation bool   `json:"allow_reconciliation"`
	RequireTaxReporting bool   `json:"require_tax_reporting"`
	Description         string `json:"description"`
	Notes               string `json:"notes"`
}

type UpdateChartOfAccountDTO struct {
	Name     string `json:"name" validate:"omitempty,max=200"`
	Category string `json:"category"`
	Active   *bool  `json:"active"`

	Currency            string `json:"currency"`
	AllowReconciliation *bool  `json:"allow_reconciliation"`
	RequireTaxReporting *bool  `json:"require_tax_reporting"`
	Description         string `json:"description"`
	Notes               string `json:"notes"`
}

type ListChartOfAccountsDTO struct {
	Page     int     `json:"page" validate:"omitempty,min=1"`
	Limit    int     `json:"limit" validate:"omitempty,min=1,max=100"`
	Search   string  `json:"search"`
	Type     string  `json:"type" validate:"omitempty,oneof=asset liability equity income expense"`
	Category string  `json:"category"`
	ParentID *string `json:"parent_id"`
	Active   *bool   `json:"active"`
	IsGroup  *bool   `json:"is_group"`
}

// ============================================================================
// TAX DTOs
// ============================================================================

type CreateTaxDTO struct {
	Code        string
	Name        string
	Description string
	Type        string
	Scope       string
	Rate        float64
	IsDefault   bool
}

type UpdateTaxDTO struct {
	Code        *string
	Name        *string
	Description *string
	Type        *string
	Scope       *string
	Rate        *float64
	IsDefault   *bool
}

type ListTaxesDTO struct {
	Page      int
	Limit     int
	Search    string
	Type      string
	Scope     string
	IsDefault *bool
}

// ============================================================================
// PAYMENT TERM DTOs
// ============================================================================

type CreatePaymentTermDTO struct {
	Code            string
	Name            string
	Description     string
	Days            int
	DiscountPercent float64
	DiscountDays    int
	IsDefault       bool
}

type UpdatePaymentTermDTO struct {
	Code            *string
	Name            *string
	Description     *string
	Days            *int
	DiscountPercent *float64
	DiscountDays    *int
	IsDefault       *bool
}

type ListPaymentTermsDTO struct {
	Page      int
	Limit     int
	Search    string
	IsDefault *bool
}

// ============================================================================
// CURRENCY DTOs
// ============================================================================

type CreateCurrencyDTO struct {
	Code          string
	Name          string
	Symbol        string
	ExchangeRate  float64
	DecimalPlaces int
	IsDefault     bool
}

type UpdateCurrencyDTO struct {
	Code          *string
	Name          *string
	Symbol        *string
	ExchangeRate  *float64
	DecimalPlaces *int
	IsDefault     *bool
}

type ListCurrenciesDTO struct {
	Page      int
	Limit     int
	Search    string
	IsDefault *bool
}
