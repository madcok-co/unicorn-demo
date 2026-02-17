package database

import (
	"testing"

	"github.com/madcok-co/unicorn-demo/internal/domain"
)

func TestConnect_SQLite(t *testing.T) {
	cfg := &Config{
		Driver:   "sqlite",
		Database: ":memory:",
	}

	db, err := ConnectRaw(cfg)
	if err != nil {
		t.Fatalf("Connect() error = %v, want nil", err)
	}
	if db == nil {
		t.Fatal("Connect() returned nil database")
	}

	// Test connection via SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("db.DB() error = %v, want nil", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("sqlDB.Ping() error = %v, want nil", err)
	}
}

func TestConnect_UnsupportedDriver(t *testing.T) {
	cfg := &Config{
		Driver:   "mysql",
		Database: "test",
	}

	_, err := ConnectRaw(cfg)
	if err == nil {
		t.Error("Connect() error = nil, want error for unsupported driver")
	}
}

func TestAutoMigrate(t *testing.T) {
	// Create in-memory database
	cfg := &Config{
		Driver:   "sqlite",
		Database: ":memory:",
	}

	db, err := ConnectRaw(cfg)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	// Run migrations
	err = AutoMigrate(db)
	if err != nil {
		t.Fatalf("AutoMigrate() error = %v, want nil", err)
	}

	// Verify tables exist
	tables := []string{"users", "projects", "tasks", "comments", "audit_logs"}
	for _, tableName := range tables {
		if !db.Migrator().HasTable(tableName) {
			t.Errorf("Table %s does not exist after migration", tableName)
		}
	}
}

func TestSeedData(t *testing.T) {
	// Create in-memory database
	cfg := &Config{
		Driver:   "sqlite",
		Database: ":memory:",
	}

	db, err := ConnectRaw(cfg)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	// Run migrations
	if err := AutoMigrate(db); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}

	// Seed data
	err = SeedData(db)
	if err != nil {
		t.Fatalf("SeedData() error = %v, want nil", err)
	}

	// Verify users were seeded
	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)
	if userCount == 0 {
		t.Error("No users were seeded")
	}

	// Verify projects were seeded
	var projectCount int64
	db.Model(&domain.Project{}).Count(&projectCount)
	if projectCount == 0 {
		t.Error("No projects were seeded")
	}

	// Verify tasks were seeded
	var taskCount int64
	db.Model(&domain.Task{}).Count(&taskCount)
	if taskCount == 0 {
		t.Error("No tasks were seeded")
	}

	// Test that seeding again skips (data already exists)
	err = SeedData(db)
	if err != nil {
		t.Fatalf("SeedData() second call error = %v, want nil", err)
	}

	// Count should be the same
	var userCount2 int64
	db.Model(&domain.User{}).Count(&userCount2)
	if userCount2 != userCount {
		t.Errorf("User count after second seed = %v, want %v (should skip)", userCount2, userCount)
	}
}

func TestSeedData_VerifyTenantData(t *testing.T) {
	// Create in-memory database
	cfg := &Config{
		Driver:   "sqlite",
		Database: ":memory:",
	}

	db, err := ConnectRaw(cfg)
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	if err := AutoMigrate(db); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}

	if err := SeedData(db); err != nil {
		t.Fatalf("SeedData() error = %v", err)
	}

	// Verify acme tenant data
	var acmeUsers []domain.User
	db.Where("tenant_id = ?", "acme").Find(&acmeUsers)
	if len(acmeUsers) < 3 {
		t.Errorf("Acme users count = %v, want at least 3", len(acmeUsers))
	}

	// Verify techcorp tenant data
	var techcorpUsers []domain.User
	db.Where("tenant_id = ?", "techcorp").Find(&techcorpUsers)
	if len(techcorpUsers) < 1 {
		t.Errorf("TechCorp users count = %v, want at least 1", len(techcorpUsers))
	}

	// Verify project ownership
	var projects []domain.Project
	db.Where("tenant_id = ?", "acme").Find(&projects)
	if len(projects) < 2 {
		t.Errorf("Acme projects count = %v, want at least 2", len(projects))
	}
}
