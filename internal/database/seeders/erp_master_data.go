package seeders

import (
	"time"

	"github.com/google/uuid"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"gorm.io/gorm"
)

// SeedERPMasterData seeds initial ERP master data for a tenant
func SeedERPMasterData(db *gorm.DB, tenantID, userID string) error {
	// Seed default currency (USD)
	defaultCurrency := domain.Currency{
		ID:            uuid.New().String(),
		Code:          "USD",
		Name:          "US Dollar",
		Symbol:        "$",
		ExchangeRate:  1.0,
		DecimalPlaces: 2,
		IsDefault:     true,
		TenantID:      tenantID,
		Active:        true,
		CreatedBy:     userID,
		UpdatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := db.Create(&defaultCurrency).Error; err != nil {
		return err
	}

	// Seed additional currencies
	currencies := []domain.Currency{
		{
			ID:            uuid.New().String(),
			Code:          "EUR",
			Name:          "Euro",
			Symbol:        "€",
			ExchangeRate:  0.92,
			DecimalPlaces: 2,
			IsDefault:     false,
			TenantID:      tenantID,
			Active:        true,
			CreatedBy:     userID,
			UpdatedBy:     userID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            uuid.New().String(),
			Code:          "IDR",
			Name:          "Indonesian Rupiah",
			Symbol:        "Rp",
			ExchangeRate:  15750.0,
			DecimalPlaces: 0,
			IsDefault:     false,
			TenantID:      tenantID,
			Active:        true,
			CreatedBy:     userID,
			UpdatedBy:     userID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}
	for _, currency := range currencies {
		if err := db.Create(&currency).Error; err != nil {
			return err
		}
	}

	// Seed default payment terms
	paymentTerms := []domain.PaymentTerm{
		{
			ID:              uuid.New().String(),
			Code:            "NET30",
			Name:            "Net 30 Days",
			Description:     "Payment due in 30 days",
			Days:            30,
			DiscountPercent: 0,
			DiscountDays:    0,
			IsDefault:       true,
			TenantID:        tenantID,
			Active:          true,
			CreatedBy:       userID,
			UpdatedBy:       userID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Code:            "NET15",
			Name:            "Net 15 Days",
			Description:     "Payment due in 15 days",
			Days:            15,
			DiscountPercent: 0,
			DiscountDays:    0,
			IsDefault:       false,
			TenantID:        tenantID,
			Active:          true,
			CreatedBy:       userID,
			UpdatedBy:       userID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Code:            "2/10NET30",
			Name:            "2/10 Net 30",
			Description:     "2% discount if paid within 10 days, otherwise due in 30 days",
			Days:            30,
			DiscountPercent: 2.0,
			DiscountDays:    10,
			IsDefault:       false,
			TenantID:        tenantID,
			Active:          true,
			CreatedBy:       userID,
			UpdatedBy:       userID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Code:            "COD",
			Name:            "Cash on Delivery",
			Description:     "Payment due on delivery",
			Days:            0,
			DiscountPercent: 0,
			DiscountDays:    0,
			IsDefault:       false,
			TenantID:        tenantID,
			Active:          true,
			CreatedBy:       userID,
			UpdatedBy:       userID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}
	for _, term := range paymentTerms {
		if err := db.Create(&term).Error; err != nil {
			return err
		}
	}

	// Seed default taxes
	taxes := []domain.Tax{
		{
			ID:          uuid.New().String(),
			Code:        "VAT10",
			Name:        "VAT 10%",
			Description: "Value Added Tax 10%",
			Type:        "sales",
			Rate:        10.0,
			IsDefault:   true,
			TenantID:    tenantID,
			Active:      true,
			CreatedBy:   userID,
			UpdatedBy:   userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Code:        "PVAT10",
			Name:        "Purchase VAT 10%",
			Description: "Purchase Value Added Tax 10%",
			Type:        "purchase",
			Rate:        10.0,
			IsDefault:   true,
			TenantID:    tenantID,
			Active:      true,
			CreatedBy:   userID,
			UpdatedBy:   userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			Code:        "EXEMPT",
			Name:        "Tax Exempt",
			Description: "No tax applied",
			Type:        "exempt",
			Rate:        0.0,
			IsDefault:   false,
			TenantID:    tenantID,
			Active:      true,
			CreatedBy:   userID,
			UpdatedBy:   userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	for _, tax := range taxes {
		if err := db.Create(&tax).Error; err != nil {
			return err
		}
	}

	// Seed default warehouse
	defaultWarehouse := domain.Warehouse{
		ID:            uuid.New().String(),
		Code:          "WH01",
		Name:          "Main Warehouse",
		Type:          "distribution",
		Address:       "123 Main Street",
		City:          "Jakarta",
		State:         "DKI Jakarta",
		Country:       "Indonesia",
		Zip:           "12345",
		Phone:         "+62-21-1234567",
		Email:         "warehouse@company.com",
		ManagerName:   "Warehouse Manager",
		TotalCapacity: 10000.0,
		IsDefault:     true,
		TenantID:      tenantID,
		Active:        true,
		CreatedBy:     userID,
		UpdatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := db.Create(&defaultWarehouse).Error; err != nil {
		return err
	}

	// Seed Chart of Accounts (simplified structure)
	accounts := []domain.ChartOfAccount{
		// Assets
		{
			ID:             uuid.New().String(),
			Code:           "1000",
			Name:           "Assets",
			Type:           "asset",
			Category:       "current",
			Description:    "Current Assets",
			ParentID:       "",
			Level:          0,
			Currency:       "USD",
			CurrentBalance: 0,
			TenantID:       tenantID,
			Active:         true,
			CreatedBy:      userID,
			UpdatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		// Liabilities
		{
			ID:             uuid.New().String(),
			Code:           "2000",
			Name:           "Liabilities",
			Type:           "liability",
			Category:       "current",
			Description:    "Current Liabilities",
			ParentID:       "",
			Level:          0,
			Currency:       "USD",
			CurrentBalance: 0,
			TenantID:       tenantID,
			Active:         true,
			CreatedBy:      userID,
			UpdatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		// Equity
		{
			ID:             uuid.New().String(),
			Code:           "3000",
			Name:           "Equity",
			Type:           "equity",
			Category:       "owner_equity",
			Description:    "Owner's Equity",
			ParentID:       "",
			Level:          0,
			Currency:       "USD",
			CurrentBalance: 0,
			TenantID:       tenantID,
			Active:         true,
			CreatedBy:      userID,
			UpdatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		// Revenue
		{
			ID:             uuid.New().String(),
			Code:           "4000",
			Name:           "Revenue",
			Type:           "revenue",
			Category:       "sales",
			Description:    "Sales Revenue",
			ParentID:       "",
			Level:          0,
			Currency:       "USD",
			CurrentBalance: 0,
			TenantID:       tenantID,
			Active:         true,
			CreatedBy:      userID,
			UpdatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		// Expenses
		{
			ID:             uuid.New().String(),
			Code:           "5000",
			Name:           "Expenses",
			Type:           "expense",
			Category:       "operating",
			Description:    "Operating Expenses",
			ParentID:       "",
			Level:          0,
			Currency:       "USD",
			CurrentBalance: 0,
			TenantID:       tenantID,
			Active:         true,
			CreatedBy:      userID,
			UpdatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for i := range accounts {
		if err := db.Create(&accounts[i]).Error; err != nil {
			return err
		}
	}

	// Add some sub-accounts
	assetID := accounts[0].ID
	subAccounts := []domain.ChartOfAccount{
		{
			ID:             uuid.New().String(),
			Code:           "1100",
			Name:           "Cash and Bank",
			Type:           "asset",
			Category:       "current",
			Description:    "Cash and Bank Accounts",
			ParentID:       assetID,
			Level:          1,
			Currency:       "USD",
			CurrentBalance: 0,
			TenantID:       tenantID,
			Active:         true,
			CreatedBy:      userID,
			UpdatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New().String(),
			Code:           "1200",
			Name:           "Accounts Receivable",
			Type:           "asset",
			Category:       "current",
			Description:    "Customer Receivables",
			ParentID:       assetID,
			Level:          1,
			Currency:       "USD",
			CurrentBalance: 0,
			TenantID:       tenantID,
			Active:         true,
			CreatedBy:      userID,
			UpdatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New().String(),
			Code:           "1300",
			Name:           "Inventory",
			Type:           "asset",
			Category:       "current",
			Description:    "Product Inventory",
			ParentID:       assetID,
			Level:          1,
			Currency:       "USD",
			CurrentBalance: 0,
			TenantID:       tenantID,
			Active:         true,
			CreatedBy:      userID,
			UpdatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for i := range subAccounts {
		if err := db.Create(&subAccounts[i]).Error; err != nil {
			return err
		}
	}

	// Seed sample customers
	customers := []domain.Customer{
		{
			ID:              uuid.New().String(),
			Code:            "CUST001",
			Name:            "ABC Corporation",
			Type:            "company",
			Email:           "contact@abc-corp.com",
			Phone:           "+1-555-0001",
			BillingAddress:  "456 Business Ave",
			BillingCity:     "New York",
			BillingState:    "NY",
			BillingCountry:  "USA",
			BillingZip:      "10001",
			ShippingAddress: "456 Business Ave",
			ShippingCity:    "New York",
			ShippingState:   "NY",
			ShippingCountry: "USA",
			ShippingZip:     "10001",
			PaymentTermDays: 30,
			ContactPerson:   "John Doe",
			TenantID:        tenantID,
			Active:          true,
			CreatedBy:       userID,
			UpdatedBy:       userID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			ID:              uuid.New().String(),
			Code:            "CUST002",
			Name:            "XYZ Retail",
			Type:            "company",
			Email:           "sales@xyz-retail.com",
			Phone:           "+1-555-0003",
			BillingAddress:  "789 Commerce St",
			BillingCity:     "Los Angeles",
			BillingState:    "CA",
			BillingCountry:  "USA",
			BillingZip:      "90001",
			ShippingAddress: "789 Commerce St",
			ShippingCity:    "Los Angeles",
			ShippingState:   "CA",
			ShippingCountry: "USA",
			ShippingZip:     "90001",
			PaymentTermDays: 15,
			ContactPerson:   "Jane Smith",
			TenantID:        tenantID,
			Active:          true,
			CreatedBy:       userID,
			UpdatedBy:       userID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	for i := range customers {
		if err := db.Create(&customers[i]).Error; err != nil {
			return err
		}
	}

	// Seed sample vendors
	vendors := []domain.Vendor{
		{
			ID:              uuid.New().String(),
			Code:            "VEND001",
			Name:            "Global Suppliers Inc",
			Type:            "company",
			Email:           "orders@global-suppliers.com",
			Phone:           "+1-555-1001",
			Address:         "111 Supply Chain Blvd",
			City:            "Chicago",
			State:           "IL",
			Country:         "USA",
			Zip:             "60601",
			PaymentTermDays: 30,
			Rating:          5,
			ContactPerson:   "Robert Johnson",
			TenantID:        tenantID,
			Active:          true,
			CreatedBy:       userID,
			UpdatedBy:       userID,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	for i := range vendors {
		if err := db.Create(&vendors[i]).Error; err != nil {
			return err
		}
	}

	// Seed sample products
	products := []domain.Product{
		{
			ID:             uuid.New().String(),
			Code:           "PROD-001",
			SKU:            "PROD-001",
			Barcode:        "1234567890123",
			Name:           "Premium Widget",
			Description:    "High-quality widget for industrial use",
			Type:           "goods",
			Category:       "Electronics",
			CanBeSold:      true,
			CanBePurchased: true,
			Weight:         1.5,
			Volume:         0.1,
			TenantID:       tenantID,
			Active:         true,
			CreatedBy:      userID,
			UpdatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New().String(),
			Code:           "SERV-001",
			SKU:            "SERV-001",
			Barcode:        "",
			Name:           "Consulting Service",
			Description:    "Professional consulting services per hour",
			Type:           "service",
			Category:       "Services",
			CanBeSold:      true,
			CanBePurchased: false,
			Weight:         0.0,
			Volume:         0.0,
			TenantID:       tenantID,
			Active:         true,
			CreatedBy:      userID,
			UpdatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for i := range products {
		if err := db.Create(&products[i]).Error; err != nil {
			return err
		}
	}

	return nil
}
