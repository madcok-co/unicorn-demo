package handlers

import (
	"time"

	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

var startTime = time.Now()

// HealthCheck returns the health status of the application
func HealthCheck(ctx *context.Context, req struct{}) (*domain.HealthResponse, error) {
	db := getDB(ctx)

	// Check database connection
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			err = sqlDB.Ping()
			if err != nil {
				return &domain.HealthResponse{
					Status:    "unhealthy",
					Timestamp: time.Now(),
					Version:   "1.0.0",
					Uptime:    time.Since(startTime).String(),
				}, nil
			}
		}
	}

	return &domain.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    time.Since(startTime).String(),
	}, nil
}
