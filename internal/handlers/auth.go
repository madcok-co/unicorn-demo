package handlers

import (
	"fmt"
	"time"

	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// Login handles user login with email/password
func Login(ctx *context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	db := getDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	// In a real app, verify password hash
	// For demo, we'll simulate authentication
	var user domain.User
	err = db.Where("email = ? AND tenant_id = ?", req.Email, tenantID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.Active {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Create audit log
	createAuditLog(ctx, user.ID, "login", "user", user.ID, nil)

	// In a real app, generate JWT token
	// For demo, we'll return a mock token
	return &domain.LoginResponse{
		Token:        fmt.Sprintf("jwt-token-%s", user.ID),
		RefreshToken: fmt.Sprintf("refresh-token-%s", user.ID),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		User:         user,
	}, nil
}

// OAuth2Callback handles OAuth2 callback
func OAuth2Callback(ctx *context.Context, req domain.OAuth2CallbackRequest) (*domain.LoginResponse, error) {
	db := getDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	// Note: For demo, we'll simulate OAuth2 flow
	// In real app, you would:
	// 1. Exchange code for token with OAuth2 provider
	// 2. Get user info from provider
	// 3. Create or update user in database

	tenantID, _ := getTenantID(ctx)
	if tenantID == "" {
		tenantID = "default"
	}

	// Simulate getting user info from OAuth provider
	email := fmt.Sprintf("oauth-user-%d@example.com", time.Now().Unix())
	name := "OAuth User"

	var user domain.User
	err := db.Where("email = ? AND tenant_id = ?", email, tenantID).First(&user).Error
	if err != nil {
		// Create new user
		user = domain.User{
			ID:       fmt.Sprintf("user-%d", time.Now().Unix()),
			Email:    email,
			Name:     name,
			TenantID: tenantID,
			Roles:    []string{"viewer"}, // Default role
			Active:   true,
		}
		db.Create(&user)
	}

	// Create audit log
	createAuditLog(ctx, user.ID, "oauth_login", "user", user.ID, nil)

	return &domain.LoginResponse{
		Token:        fmt.Sprintf("jwt-token-%s", user.ID),
		RefreshToken: fmt.Sprintf("refresh-token-%s", user.ID),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		User:         user,
	}, nil
}

// RefreshToken handles token refresh
func RefreshToken(ctx *context.Context, req domain.RefreshTokenRequest) (*domain.LoginResponse, error) {
	db := getDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	// In real app, validate refresh token and generate new access token
	// For demo, we'll extract user ID from refresh token
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &domain.LoginResponse{
		Token:        fmt.Sprintf("jwt-token-%s", user.ID),
		RefreshToken: fmt.Sprintf("refresh-token-%s", user.ID),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		User:         user,
	}, nil
}

// GetCurrentUser returns the currently authenticated user
func GetCurrentUser(ctx *context.Context, req struct{}) (*domain.User, error) {
	db := getDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

// Logout handles user logout
func Logout(ctx *context.Context, req struct{}) (*domain.MessageResponse, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Create audit log
	createAuditLog(ctx, userID, "logout", "user", userID, nil)

	// In a real app, invalidate token (blacklist or delete from cache)
	return &domain.MessageResponse{
		Message: "Logged out successfully",
	}, nil
}
