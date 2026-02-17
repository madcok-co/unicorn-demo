package handlers

import (
	"fmt"
	"time"

	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/core/pkg/context"
	"gorm.io/gorm"
)

// getDB returns the underlying GORM database instance
func getDB(ctx *context.Context) *gorm.DB {
	// First try to get DB directly from context (for testing)
	if db, exists := ctx.Get("db"); exists && db != nil {
		if gormDB, ok := db.(*gorm.DB); ok {
			return gormDB
		}
	}

	// Access DB via context's DB() method (lazy loaded from AppAdapters)
	dbInterface := ctx.DB()
	if dbInterface == nil {
		return nil
	}

	// Unwrap GORM DB from Unicorn's database driver
	if gormDriver, ok := dbInterface.(interface{ DB() *gorm.DB }); ok {
		return gormDriver.DB()
	}

	return nil
}

// getTenantID gets tenant ID from context
func getTenantID(ctx *context.Context) (string, error) {
	tenantID, exists := ctx.Get("tenant_id")
	if !exists || tenantID == nil {
		return "", fmt.Errorf("tenant not found")
	}
	return tenantID.(string), nil
}

// getUserID gets user ID from context
func getUserID(ctx *context.Context) (string, error) {
	userID, exists := ctx.Get("user_id")
	if !exists || userID == nil {
		return "", fmt.Errorf("user not authenticated")
	}
	return userID.(string), nil
}

// getParam gets URL parameter from context
func getParam(ctx *context.Context, key string) string {
	if param, exists := ctx.Get("param_" + key); exists && param != nil {
		return param.(string)
	}
	return ""
}

// createAuditLog creates an audit log entry
func createAuditLog(ctx *context.Context, userID, action, resource, resourceID string, changes map[string]interface{}) {
	db := getDB(ctx)
	if db == nil {
		return
	}

	tenantID, _ := getTenantID(ctx)

	// Get request info safely
	ipAddress := "unknown"
	userAgent := "unknown"

	if req := ctx.Request(); req != nil {
		// Get remote address
		if addr := req.Header("X-Real-IP"); addr != "" {
			ipAddress = addr
		} else if addr := req.Header("X-Forwarded-For"); addr != "" {
			ipAddress = addr
		}

		// Get user agent
		if ua := req.Header("User-Agent"); ua != "" {
			userAgent = ua
		}
	}

	log := domain.AuditLog{
		ID:         fmt.Sprintf("audit-%d", time.Now().UnixNano()),
		TenantID:   tenantID,
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Changes:    changes,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		CreatedAt:  time.Now(),
	}

	db.Create(&log)
}
