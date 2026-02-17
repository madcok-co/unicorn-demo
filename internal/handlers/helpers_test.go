package handlers

import (
	"context"
	"testing"

	ucontext "github.com/madcok-co/unicorn/core/pkg/context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetDB(t *testing.T) {
	// Create in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create context and set DB
	ctx := ucontext.New(context.Background())
	ctx.Set("db", db)

	// Test getDB helper
	result := getDB(ctx)
	if result == nil {
		t.Error("getDB() returned nil, want *gorm.DB")
	}
	if result != db {
		t.Error("getDB() returned different DB instance")
	}
}

func TestGetDB_NotSet(t *testing.T) {
	// Create context without DB
	ctx := ucontext.New(context.Background())

	// Test getDB helper
	result := getDB(ctx)
	if result != nil {
		t.Error("getDB() returned non-nil, want nil when DB not set")
	}
}

func TestGetTenantID(t *testing.T) {
	tests := []struct {
		name      string
		tenantID  interface{}
		shouldSet bool
		wantErr   bool
		wantID    string
	}{
		{
			name:      "valid tenant ID",
			tenantID:  "acme",
			shouldSet: true,
			wantErr:   false,
			wantID:    "acme",
		},
		{
			name:      "tenant ID not set",
			shouldSet: false,
			wantErr:   true,
			wantID:    "",
		},
		{
			name:      "tenant ID is nil",
			tenantID:  nil,
			shouldSet: true,
			wantErr:   true,
			wantID:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ucontext.New(context.Background())
			if tt.shouldSet {
				ctx.Set("tenant_id", tt.tenantID)
			}

			gotID, err := getTenantID(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTenantID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("getTenantID() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name      string
		userID    interface{}
		shouldSet bool
		wantErr   bool
		wantID    string
	}{
		{
			name:      "valid user ID",
			userID:    "user-123",
			shouldSet: true,
			wantErr:   false,
			wantID:    "user-123",
		},
		{
			name:      "user ID not set",
			shouldSet: false,
			wantErr:   true,
			wantID:    "",
		},
		{
			name:      "user ID is nil",
			userID:    nil,
			shouldSet: true,
			wantErr:   true,
			wantID:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ucontext.New(context.Background())
			if tt.shouldSet {
				ctx.Set("user_id", tt.userID)
			}

			gotID, err := getUserID(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("getUserID() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}

func TestGetParam(t *testing.T) {
	tests := []struct {
		name      string
		paramKey  string
		paramVal  string
		wantValue string
	}{
		{
			name:      "param exists",
			paramKey:  "id",
			paramVal:  "123",
			wantValue: "123",
		},
		{
			name:      "param does not exist",
			paramKey:  "id",
			paramVal:  "",
			wantValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ucontext.New(context.Background())
			if tt.paramVal != "" {
				ctx.Set("param_"+tt.paramKey, tt.paramVal)
			}

			gotValue := getParam(ctx, tt.paramKey)
			if gotValue != tt.wantValue {
				t.Errorf("getParam() = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}
