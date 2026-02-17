package domain

import "time"

// ============================================================================
// MASTER DATA MODELS
// ============================================================================

// Customer represents a customer/client in the system
type Customer struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_customer_code_tenant;not null"` // CUST-001
	Name     string `json:"name" gorm:"not null"`
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_customer_code_tenant;index;not null"`
	Type     string `json:"type" gorm:"not null"` // individual, company
	Email    string `json:"email" gorm:"index"`
	Phone    string `json:"phone"`
	Mobile   string `json:"mobile"`
	TaxID    string `json:"tax_id"` // NPWP or Tax ID
	Website  string `json:"website"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Address
	BillingAddress string `json:"billing_address"`
	BillingCity    string `json:"billing_city"`
	BillingState   string `json:"billing_state"`
	BillingZip     string `json:"billing_zip"`
	BillingCountry string `json:"billing_country" gorm:"default:'Indonesia'"`

	ShippingAddress string `json:"shipping_address"`
	ShippingCity    string `json:"shipping_city"`
	ShippingState   string `json:"shipping_state"`
	ShippingZip     string `json:"shipping_zip"`
	ShippingCountry string `json:"shipping_country" gorm:"default:'Indonesia'"`

	// Financial
	CreditLimit     float64 `json:"credit_limit" gorm:"default:0"`
	PaymentTermDays int     `json:"payment_term_days" gorm:"default:30"` // Net 30, 45, etc
	Currency        string  `json:"currency" gorm:"default:'IDR'"`

	// Contact Person
	ContactPerson    string `json:"contact_person"`
	ContactPersonJob string `json:"contact_person_job"`

	// Metadata
	Notes     string    `json:"notes" gorm:"type:text"`
	Tags      []string  `json:"tags" gorm:"serializer:json"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Customer) TableName() string {
	return "customers"
}

// Vendor represents a supplier/vendor in the system
type Vendor struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_vendor_code_tenant;not null"` // VEND-001
	Name     string `json:"name" gorm:"not null"`
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_vendor_code_tenant;index;not null"`
	Type     string `json:"type" gorm:"not null"` // individual, company
	Email    string `json:"email" gorm:"index"`
	Phone    string `json:"phone"`
	Mobile   string `json:"mobile"`
	TaxID    string `json:"tax_id"`
	Website  string `json:"website"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Address
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country" gorm:"default:'Indonesia'"`

	// Financial
	PaymentTermDays int    `json:"payment_term_days" gorm:"default:30"`
	Currency        string `json:"currency" gorm:"default:'IDR'"`
	BankName        string `json:"bank_name"`
	BankAccount     string `json:"bank_account"`
	BankAccountName string `json:"bank_account_name"`

	// Contact Person
	ContactPerson    string `json:"contact_person"`
	ContactPersonJob string `json:"contact_person_job"`

	// Rating & Category
	Rating   int    `json:"rating" gorm:"default:0"` // 1-5 stars
	Category string `json:"category"`                // raw-material, service, etc

	// Metadata
	Notes     string    `json:"notes" gorm:"type:text"`
	Tags      []string  `json:"tags" gorm:"serializer:json"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Vendor) TableName() string {
	return "vendors"
}

// Product represents a product/item in the system
type Product struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_product_code_tenant;not null"` // PROD-001
	Name     string `json:"name" gorm:"not null"`
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_product_code_tenant;index;not null"`
	Type     string `json:"type" gorm:"not null"` // product, service, consumable
	Category string `json:"category"`             // electronics, food, etc
	UOM      string `json:"uom" gorm:"not null"`  // Unit of Measure: pcs, kg, liter, etc
	Barcode  string `json:"barcode" gorm:"index"`
	SKU      string `json:"sku" gorm:"index"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Inventory
	CanBeSold      bool    `json:"can_be_sold" gorm:"default:true"`
	CanBePurchased bool    `json:"can_be_purchased" gorm:"default:true"`
	TrackInventory bool    `json:"track_inventory" gorm:"default:true"`
	MinStock       float64 `json:"min_stock" gorm:"default:0"`
	MaxStock       float64 `json:"max_stock" gorm:"default:0"`
	ReorderLevel   float64 `json:"reorder_level" gorm:"default:0"`

	// Pricing
	SalePrice     float64 `json:"sale_price" gorm:"default:0"`
	PurchasePrice float64 `json:"purchase_price" gorm:"default:0"`
	Cost          float64 `json:"cost" gorm:"default:0"` // Standard cost
	Currency      string  `json:"currency" gorm:"default:'IDR'"`

	// Tax
	TaxCategory string  `json:"tax_category"` // taxable, non-taxable
	TaxRate     float64 `json:"tax_rate" gorm:"default:0"`

	// Physical
	Weight       float64 `json:"weight" gorm:"default:0"`        // in kg
	Volume       float64 `json:"volume" gorm:"default:0"`        // in m3
	Length       float64 `json:"length" gorm:"default:0"`        // in cm
	Width        float64 `json:"width" gorm:"default:0"`         // in cm
	Height       float64 `json:"height" gorm:"default:0"`        // in cm
	WarrantyDays int     `json:"warranty_days" gorm:"default:0"` // warranty period

	// Description
	Description      string   `json:"description" gorm:"type:text"`
	InternalNotes    string   `json:"internal_notes" gorm:"type:text"`
	ImageURL         string   `json:"image_url"`
	AdditionalImages []string `json:"additional_images" gorm:"serializer:json"`

	// Accounting
	IncomeAccountID  string `json:"income_account_id"`  // Chart of Account ID
	ExpenseAccountID string `json:"expense_account_id"` // Chart of Account ID
	AssetAccountID   string `json:"asset_account_id"`   // Chart of Account ID

	// Metadata
	Tags      []string  `json:"tags" gorm:"serializer:json"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "products"
}

// Warehouse represents a storage location
type Warehouse struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_warehouse_code_tenant;not null"` // WH-001
	Name     string `json:"name" gorm:"not null"`
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_warehouse_code_tenant;index;not null"`
	Type     string `json:"type" gorm:"not null"` // physical, virtual, transit
	Active   bool   `json:"active" gorm:"default:true"`

	// Location
	Address string `json:"address"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country" gorm:"default:'Indonesia'"`

	// Contact
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	ManagerID   string `json:"manager_id"` // User ID
	ManagerName string `json:"manager_name"`

	// Capacity
	TotalArea     float64 `json:"total_area" gorm:"default:0"`     // in m2
	UsedArea      float64 `json:"used_area" gorm:"default:0"`      // in m2
	TotalCapacity float64 `json:"total_capacity" gorm:"default:0"` // in units
	UsedCapacity  float64 `json:"used_capacity" gorm:"default:0"`  // in units

	// Settings
	AllowNegativeStock bool `json:"allow_negative_stock" gorm:"default:false"`
	IsDefault          bool `json:"is_default" gorm:"default:false"`

	// Metadata
	Description string    `json:"description" gorm:"type:text"`
	Notes       string    `json:"notes" gorm:"type:text"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Warehouse) TableName() string {
	return "warehouses"
}

// ChartOfAccount represents accounting accounts (COA)
type ChartOfAccount struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_coa_code_tenant;not null"` // 1-1001
	Name     string `json:"name" gorm:"not null"`
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_coa_code_tenant;index;not null"`
	Type     string `json:"type" gorm:"not null"`          // asset, liability, equity, income, expense
	Category string `json:"category"`                      // current_asset, fixed_asset, etc
	ParentID string `json:"parent_id" gorm:"index"`        // For hierarchical COA
	Level    int    `json:"level" gorm:"default:1"`        // 1, 2, 3 (account depth)
	IsGroup  bool   `json:"is_group" gorm:"default:false"` // Is this a group/header account
	Active   bool   `json:"active" gorm:"default:true"`

	// Financial
	Currency       string  `json:"currency" gorm:"default:'IDR'"`
	CurrentBalance float64 `json:"current_balance" gorm:"default:0"` // Calculated field
	DebitBalance   float64 `json:"debit_balance" gorm:"default:0"`
	CreditBalance  float64 `json:"credit_balance" gorm:"default:0"`

	// Settings
	AllowReconciliation bool `json:"allow_reconciliation" gorm:"default:false"` // For bank accounts
	RequireTaxReporting bool `json:"require_tax_reporting" gorm:"default:false"`

	// Metadata
	Description string    `json:"description" gorm:"type:text"`
	Notes       string    `json:"notes" gorm:"type:text"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ChartOfAccount) TableName() string {
	return "chart_of_accounts"
}

// Tax represents tax configuration
type Tax struct {
	ID       string  `json:"id" gorm:"primaryKey"`
	Code     string  `json:"code" gorm:"uniqueIndex:idx_tax_code_tenant;not null"` // PPN, PPH
	Name     string  `json:"name" gorm:"not null"`
	TenantID string  `json:"tenant_id" gorm:"uniqueIndex:idx_tax_code_tenant;index;not null"`
	Type     string  `json:"type" gorm:"not null"` // sales, purchase, both
	Rate     float64 `json:"rate" gorm:"not null"` // 11% = 11.0
	Active   bool    `json:"active" gorm:"default:true"`

	// Accounting
	TaxAccountID string `json:"tax_account_id"` // COA for tax payable/receivable

	// Settings
	IsDefault       bool `json:"is_default" gorm:"default:false"`
	IncludedInPrice bool `json:"included_in_price" gorm:"default:false"` // Tax-inclusive pricing

	// Metadata
	Description string    `json:"description" gorm:"type:text"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Tax) TableName() string {
	return "taxes"
}

// PaymentTerm represents payment terms configuration
type PaymentTerm struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_payment_term_code_tenant;not null"` // NET30
	Name     string `json:"name" gorm:"not null"`                                          // Net 30 Days
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_payment_term_code_tenant;index;not null"`
	Days     int    `json:"days" gorm:"not null"` // 30, 45, 60 days
	Type     string `json:"type" gorm:"not null"` // fixed, eom (end of month)
	Active   bool   `json:"active" gorm:"default:true"`

	// Discount for early payment
	DiscountDays    int     `json:"discount_days" gorm:"default:0"`    // Discount if paid within X days
	DiscountPercent float64 `json:"discount_percent" gorm:"default:0"` // Discount percentage

	// Settings
	IsDefault bool `json:"is_default" gorm:"default:false"`

	// Metadata
	Description string    `json:"description" gorm:"type:text"`
	CreatedBy   string    `json:"created_by"`
	UpdatedBy   string    `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PaymentTerm) TableName() string {
	return "payment_terms"
}

// Currency represents multi-currency support
type Currency struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Code     string `json:"code" gorm:"uniqueIndex:idx_currency_code_tenant;not null"` // USD, EUR, IDR
	Name     string `json:"name" gorm:"not null"`                                      // US Dollar, Euro
	Symbol   string `json:"symbol" gorm:"not null"`                                    // $, €, Rp
	TenantID string `json:"tenant_id" gorm:"uniqueIndex:idx_currency_code_tenant;index;not null"`
	Active   bool   `json:"active" gorm:"default:true"`

	// Exchange rate and default
	ExchangeRate float64 `json:"exchange_rate" gorm:"not null;default:1"` // Exchange rate to base currency
	IsDefault    bool    `json:"is_default" gorm:"default:false"`         // Default currency for the tenant

	// Formatting
	DecimalPlaces int `json:"decimal_places" gorm:"default:2"`

	// Metadata
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Currency) TableName() string {
	return "currencies"
}
