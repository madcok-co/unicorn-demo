package handlers

import (
	"context"
	"testing"

	ucontext "github.com/madcok-co/unicorn/core/pkg/context"
)

func TestHealthCheck(t *testing.T) {
	// Create a test context
	ctx := ucontext.New(context.Background())

	// Call the handler with empty request
	response, err := HealthCheck(ctx, struct{}{})

	// Check if the handler executed without error
	if err != nil {
		t.Errorf("HealthCheck() error = %v, want nil", err)
	}

	// Check response
	if response == nil {
		t.Fatal("HealthCheck() response is nil")
	}

	if response.Status != "healthy" && response.Status != "unhealthy" {
		t.Errorf("HealthCheck() status = %v, want 'healthy' or 'unhealthy'", response.Status)
	}

	if response.Version == "" {
		t.Error("HealthCheck() version is empty")
	}
}
