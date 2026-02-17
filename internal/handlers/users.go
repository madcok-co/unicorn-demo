package handlers

import (
	"fmt"
	"time"

	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
	"github.com/madcok-co/unicorn/core/pkg/contracts"
)

// ListUsers returns paginated list of users
func ListUsers(ctx *context.Context, req domain.ListUsersRequest) (*pagination.OffsetResult, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Check authorization
	identity := &contracts.Identity{ID: userID}
	allowed, err := ctx.Authz().Authorize(ctx.Context(), identity, "read", "users")
	if err != nil || !allowed {
		return nil, fmt.Errorf("forbidden: insufficient permissions")
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Sort == "" {
		req.Sort = "created_at"
	}
	if req.Order == "" {
		req.Order = "desc"
	}

	// Build query
	query := getDB(ctx).Where("tenant_id = ?", tenantID)

	if req.Role != "" {
		query = query.Where("roles @> ?", fmt.Sprintf(`["%s"]`, req.Role))
	}
	if req.Active != nil {
		query = query.Where("active = ?", *req.Active)
	}

	// Get total count
	var total int64
	query.Model(&domain.User{}).Count(&total)

	// Get paginated results
	var users []domain.User
	offset := (req.Page - 1) * req.Limit
	query.Order(fmt.Sprintf("%s %s", req.Sort, req.Order)).
		Limit(req.Limit).
		Offset(offset).
		Find(&users)

	params := pagination.ParseOffsetParams(req.Page, req.Limit, req.Sort, req.Order)
	return pagination.NewOffsetResult(users, total, params), nil
}

// CreateUser creates a new user
func CreateUser(ctx *context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Check authorization
	identity := &contracts.Identity{ID: userID}
	allowed, err := ctx.Authz().Authorize(ctx.Context(), identity, "create", "users")
	if err != nil || !allowed {
		return nil, fmt.Errorf("forbidden: insufficient permissions")
	}

	// Check if email already exists
	var existingUser domain.User
	err = getDB(ctx).Where("email = ? AND tenant_id = ?", req.Email, tenantID).First(&existingUser).Error
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &domain.User{
		ID:       fmt.Sprintf("user-%d", time.Now().UnixNano()),
		Email:    req.Email,
		Name:     req.Name,
		TenantID: tenantID,
		Roles:    req.Roles,
		Active:   true,
	}

	err = getDB(ctx).Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create audit log
	createAuditLog(ctx, userID, "create", "user", user.ID, map[string]interface{}{
		"email": user.Email,
		"roles": user.Roles,
	})

	return user, nil
}

// GetUser returns a single user by ID
func GetUser(ctx *context.Context, req struct{}) (*domain.User, error) {
	userIDParam := getParam(ctx, "id")
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	currentUserID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Check authorization
	identity := &contracts.Identity{ID: currentUserID}
	allowed, err := ctx.Authz().Authorize(ctx.Context(), identity, "read", "users")
	if err != nil || !allowed {
		return nil, fmt.Errorf("forbidden: insufficient permissions")
	}

	var user domain.User
	err = getDB(ctx).Where("id = ? AND tenant_id = ?", userIDParam, tenantID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

// UpdateUser updates an existing user
func UpdateUser(ctx *context.Context, req domain.UpdateUserRequest) (*domain.User, error) {
	userIDParam := getParam(ctx, "id")
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	currentUserID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Check authorization
	identity := &contracts.Identity{ID: currentUserID}
	allowed, err := ctx.Authz().Authorize(ctx.Context(), identity, "update", "users")
	if err != nil || !allowed {
		return nil, fmt.Errorf("forbidden: insufficient permissions")
	}

	var user domain.User
	err = getDB(ctx).Where("id = ? AND tenant_id = ?", userIDParam, tenantID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Track changes
	changes := make(map[string]interface{})

	if req.Name != "" {
		changes["name"] = req.Name
		user.Name = req.Name
	}
	if len(req.Roles) > 0 {
		changes["roles"] = req.Roles
		user.Roles = req.Roles
	}
	if req.Active != nil {
		changes["active"] = *req.Active
		user.Active = *req.Active
	}

	err = getDB(ctx).Save(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	createAuditLog(ctx, currentUserID, "update", "user", user.ID, changes)

	return &user, nil
}

// DeleteUser deletes a user
func DeleteUser(ctx *context.Context, req struct{}) (*domain.MessageResponse, error) {
	userIDParam := getParam(ctx, "id")
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	currentUserID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Check authorization
	identity := &contracts.Identity{ID: currentUserID}
	allowed, err := ctx.Authz().Authorize(ctx.Context(), identity, "delete", "users")
	if err != nil || !allowed {
		return nil, fmt.Errorf("forbidden: insufficient permissions")
	}

	// Prevent self-deletion
	if userIDParam == currentUserID {
		return nil, fmt.Errorf("cannot delete your own account")
	}

	var user domain.User
	err = getDB(ctx).Where("id = ? AND tenant_id = ?", userIDParam, tenantID).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	err = getDB(ctx).Delete(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	// Create audit log
	createAuditLog(ctx, currentUserID, "delete", "user", user.ID, map[string]interface{}{
		"email": user.Email,
	})

	return &domain.MessageResponse{
		Message: "User deleted successfully",
	}, nil
}

// GetUserStats returns user statistics
func GetUserStats(ctx *context.Context, req struct{}) (*domain.UserStatsResponse, error) {
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	// Check authorization
	identity := &contracts.Identity{ID: userID}
	allowed, err := ctx.Authz().Authorize(ctx.Context(), identity, "read", "users")
	if err != nil || !allowed {
		return nil, fmt.Errorf("forbidden: insufficient permissions")
	}

	stats := &domain.UserStatsResponse{
		UsersByRole: make(map[string]int64),
	}

	// Total users
	getDB(ctx).Model(&domain.User{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalUsers)

	// Active users
	getDB(ctx).Model(&domain.User{}).Where("tenant_id = ? AND active = ?", tenantID, true).Count(&stats.ActiveUsers)

	// Users by role (simplified - in real app would use JSON queries)
	var users []domain.User
	getDB(ctx).Where("tenant_id = ?", tenantID).Find(&users)

	roleCounts := make(map[string]int64)
	for _, user := range users {
		for _, role := range user.Roles {
			roleCounts[role]++
		}
	}
	stats.UsersByRole = roleCounts

	// Recent signups
	getDB(ctx).Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(5).
		Find(&stats.RecentSignups)

	return stats, nil
}
