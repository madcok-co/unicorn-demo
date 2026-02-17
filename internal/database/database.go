package database

import (
	"fmt"
	"log"
	"time"

	"github.com/madcok-co/unicorn-demo/internal/database/seeders"
	"github.com/madcok-co/unicorn-demo/internal/domain"
	gormDriver "github.com/madcok-co/unicorn/contrib/database/gorm"
	"github.com/madcok-co/unicorn/core/pkg/contracts"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Driver   string // postgres, sqlite
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// Connect creates a new database connection and returns Unicorn's Database interface
func Connect(cfg *Config) (contracts.Database, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("Connected to %s database successfully", cfg.Driver)

	// Wrap GORM with Unicorn's database driver
	return gormDriver.NewDriver(db), nil
}

// ConnectRaw creates a raw GORM connection (for migrations and seeding)
func ConnectRaw(cfg *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB for connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Printf("Connected to %s database successfully", cfg.Driver)
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Core tables
	err := db.AutoMigrate(
		&domain.User{},
		&domain.Project{},
		&domain.Task{},
		&domain.Comment{},
		&domain.AuditLog{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate core tables: %w", err)
	}

	// ERP Master Data tables
	err = db.AutoMigrate(
		&domain.Customer{},
		&domain.Vendor{},
		&domain.Product{},
		&domain.Warehouse{},
		&domain.ChartOfAccount{},
		&domain.Tax{},
		&domain.PaymentTerm{},
		&domain.Currency{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate ERP tables: %w", err)
	}

	// Sales Module tables
	err = db.AutoMigrate(
		&domain.SalesQuotation{},
		&domain.SalesQuotationItem{},
		&domain.SalesOrder{},
		&domain.SalesOrderItem{},
		&domain.DeliveryOrder{},
		&domain.DeliveryOrderItem{},
		&domain.SalesInvoice{},
		&domain.SalesInvoiceItem{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate Sales tables: %w", err)
	}

	// Purchase Module tables
	err = db.AutoMigrate(
		&domain.PurchaseRequest{},
		&domain.PurchaseRequestItem{},
		&domain.PurchaseOrder{},
		&domain.PurchaseOrderItem{},
		&domain.GoodsReceipt{},
		&domain.GoodsReceiptItem{},
		&domain.PurchaseInvoice{},
		&domain.PurchaseInvoiceItem{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate Purchase tables: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func SeedData(db *gorm.DB) error {
	log.Println("Seeding initial data...")

	// Check if data already exists
	var count int64
	db.Model(&domain.User{}).Count(&count)
	if count > 0 {
		log.Println("Data already exists, skipping seed")
		return nil
	}

	// Seed users for different tenants
	users := []domain.User{
		{
			ID:       "user-admin-acme",
			Email:    "admin@acme.com",
			Name:     "Admin Acme",
			TenantID: "acme",
			Roles:    []string{"tenant_admin"},
			Active:   true,
		},
		{
			ID:       "user-pm-acme",
			Email:    "pm@acme.com",
			Name:     "Project Manager Acme",
			TenantID: "acme",
			Roles:    []string{"project_manager"},
			Active:   true,
		},
		{
			ID:       "user-dev-acme",
			Email:    "dev@acme.com",
			Name:     "Developer Acme",
			TenantID: "acme",
			Roles:    []string{"developer"},
			Active:   true,
		},
		{
			ID:       "user-admin-techcorp",
			Email:    "admin@techcorp.com",
			Name:     "Admin TechCorp",
			TenantID: "techcorp",
			Roles:    []string{"tenant_admin"},
			Active:   true,
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to seed user: %w", err)
		}
	}

	// Seed projects
	projects := []domain.Project{
		{
			ID:          "proj-1",
			Name:        "Project Alpha",
			Description: "First demo project for Acme Corp",
			TenantID:    "acme",
			OwnerID:     "user-admin-acme",
			Status:      "active",
			Tags:        []string{"demo", "important"},
		},
		{
			ID:          "proj-2",
			Name:        "Project Beta",
			Description: "Second demo project for Acme Corp",
			TenantID:    "acme",
			OwnerID:     "user-pm-acme",
			Status:      "active",
			Tags:        []string{"demo"},
		},
		{
			ID:          "proj-3",
			Name:        "Project Gamma",
			Description: "Demo project for TechCorp",
			TenantID:    "techcorp",
			OwnerID:     "user-admin-techcorp",
			Status:      "active",
			Tags:        []string{"demo"},
		},
	}

	for _, project := range projects {
		if err := db.Create(&project).Error; err != nil {
			return fmt.Errorf("failed to seed project: %w", err)
		}
	}

	// Seed tasks
	tasks := []domain.Task{
		{
			ID:          "task-1",
			ProjectID:   "proj-1",
			Title:       "Setup development environment",
			Description: "Install all necessary tools and dependencies",
			AssigneeID:  "user-dev-acme",
			Priority:    "high",
			Status:      "done",
		},
		{
			ID:          "task-2",
			ProjectID:   "proj-1",
			Title:       "Implement authentication",
			Description: "Add OAuth2 and JWT authentication",
			AssigneeID:  "user-dev-acme",
			Priority:    "high",
			Status:      "in_progress",
		},
		{
			ID:          "task-3",
			ProjectID:   "proj-1",
			Title:       "Write documentation",
			Description: "Document all API endpoints",
			AssigneeID:  "user-pm-acme",
			Priority:    "medium",
			Status:      "todo",
		},
	}

	for _, task := range tasks {
		if err := db.Create(&task).Error; err != nil {
			return fmt.Errorf("failed to seed task: %w", err)
		}
	}

	// Seed ERP master data for acme tenant
	log.Println("Seeding ERP master data for acme tenant...")
	if err := seeders.SeedERPMasterData(db, "acme", "user-admin-acme"); err != nil {
		return fmt.Errorf("failed to seed ERP data for acme: %w", err)
	}

	// Seed ERP master data for techcorp tenant
	log.Println("Seeding ERP master data for techcorp tenant...")
	if err := seeders.SeedERPMasterData(db, "techcorp", "user-admin-techcorp"); err != nil {
		return fmt.Errorf("failed to seed ERP data for techcorp: %w", err)
	}

	// Seed Sales data for acme tenant
	log.Println("Seeding Sales data for acme tenant...")
	if err := seeders.SeedSalesData(db, "acme", "user-admin-acme"); err != nil {
		return fmt.Errorf("failed to seed Sales data for acme: %w", err)
	}

	// Seed Sales data for techcorp tenant
	log.Println("Seeding Sales data for techcorp tenant...")
	if err := seeders.SeedSalesData(db, "techcorp", "user-admin-techcorp"); err != nil {
		return fmt.Errorf("failed to seed Sales data for techcorp: %w", err)
	}

	// Seed Purchase data for acme tenant
	log.Println("Seeding Purchase data for acme tenant...")
	if err := seeders.SeedPurchaseData(db, "acme", "user-admin-acme"); err != nil {
		return fmt.Errorf("failed to seed Purchase data for acme: %w", err)
	}

	// Seed Purchase data for techcorp tenant
	log.Println("Seeding Purchase data for techcorp tenant...")
	if err := seeders.SeedPurchaseData(db, "techcorp", "user-admin-techcorp"); err != nil {
		return fmt.Errorf("failed to seed Purchase data for techcorp: %w", err)
	}

	log.Println("Initial data seeded successfully")
	return nil
}
