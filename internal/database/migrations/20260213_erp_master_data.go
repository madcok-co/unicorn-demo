package migrations

import (
	"github.com/madcok-co/unicorn-demo/internal/domain"
	"gorm.io/gorm"
)

// MigrateERPMasterData runs all ERP master data migrations
func MigrateERPMasterData(db *gorm.DB) error {
	// Create all master data tables
	if err := db.AutoMigrate(
		&domain.Customer{},
		&domain.Vendor{},
		&domain.Product{},
		&domain.Warehouse{},
		&domain.ChartOfAccount{},
		&domain.Tax{},
		&domain.PaymentTerm{},
		&domain.Currency{},
	); err != nil {
		return err
	}

	return nil
}

// RollbackERPMasterData rolls back all ERP master data migrations
func RollbackERPMasterData(db *gorm.DB) error {
	// Drop all master data tables in reverse order
	if err := db.Migrator().DropTable(
		&domain.Currency{},
		&domain.PaymentTerm{},
		&domain.Tax{},
		&domain.ChartOfAccount{},
		&domain.Warehouse{},
		&domain.Product{},
		&domain.Vendor{},
		&domain.Customer{},
	); err != nil {
		return err
	}

	return nil
}
