package handlers

import (
	"fmt"
	"time"

	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
	"github.com/madcok-co/unicorn/core/pkg/contracts"
)

// ListTasks returns paginated list of tasks
func ListTasks(ctx *context.Context, req domain.ListTasksRequest) (*pagination.OffsetResult, error) {
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
	allowed, err := ctx.Authz().Authorize(ctx.Context(), identity, "read", "tasks")
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

	// Build query - join with projects to filter by tenant
	query := getDB(ctx).Model(&domain.Task{}).
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("projects.tenant_id = ?", tenantID)

	if req.ProjectID != "" {
		query = query.Where("tasks.project_id = ?", req.ProjectID)
	}
	if req.AssigneeID != "" {
		query = query.Where("tasks.assignee_id = ?", req.AssigneeID)
	}
	if req.Status != "" {
		query = query.Where("tasks.status = ?", req.Status)
	}
	if req.Priority != "" {
		query = query.Where("tasks.priority = ?", req.Priority)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get paginated results
	var tasks []domain.Task
	offset := (req.Page - 1) * req.Limit
	query.Order(fmt.Sprintf("tasks.%s %s", req.Sort, req.Order)).
		Limit(req.Limit).
		Offset(offset).
		Find(&tasks)

	params := pagination.ParseOffsetParams(req.Page, req.Limit, req.Sort, req.Order)
	return pagination.NewOffsetResult(tasks, total, params), nil
}

// CreateTask creates a new task
func CreateTask(ctx *context.Context, req domain.CreateTaskRequest) (*domain.Task, error) {
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
	allowed, err := ctx.Authz().Authorize(ctx.Context(), identity, "create", "tasks")
	if err != nil || !allowed {
		return nil, fmt.Errorf("forbidden: insufficient permissions")
	}

	// Verify project exists and belongs to tenant
	var project domain.Project
	err = getDB(ctx).Where("id = ? AND tenant_id = ?", req.ProjectID, tenantID).First(&project).Error
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	task := &domain.Task{
		ID:          fmt.Sprintf("task-%d", time.Now().UnixNano()),
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
		AssigneeID:  req.AssigneeID,
		Priority:    req.Priority,
		Status:      "todo",
		DueDate:     req.DueDate,
	}

	if task.Priority == "" {
		task.Priority = "medium"
	}

	err = getDB(ctx).Create(task).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Create audit log
	createAuditLog(ctx, userID, "create", "task", task.ID, map[string]interface{}{
		"title":      task.Title,
		"project_id": task.ProjectID,
	})

	return task, nil
}

// GetTask returns a single task by ID
func GetTask(ctx *context.Context, req struct{}) (*domain.Task, error) {
	taskID := getParam(ctx, "id")
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
	allowed, err := ctx.Authz().Authorize(ctx.Context(), identity, "read", "tasks")
	if err != nil || !allowed {
		return nil, fmt.Errorf("forbidden: insufficient permissions")
	}

	var task domain.Task
	err = getDB(ctx).Model(&domain.Task{}).
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("tasks.id = ? AND projects.tenant_id = ?", taskID, tenantID).
		First(&task).Error
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}

	return &task, nil
}

// UpdateTask updates an existing task
func UpdateTask(ctx *context.Context, req domain.UpdateTaskRequest) (*domain.Task, error) {
	taskID := getParam(ctx, "id")
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
	allowed, err := ctx.Authz().Authorize(ctx.Context(), identity, "update", "tasks")
	if err != nil || !allowed {
		return nil, fmt.Errorf("forbidden: insufficient permissions")
	}

	var task domain.Task
	err = getDB(ctx).Model(&domain.Task{}).
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("tasks.id = ? AND projects.tenant_id = ?", taskID, tenantID).
		First(&task).Error
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}

	// Track changes
	changes := make(map[string]interface{})

	if req.Title != "" {
		changes["title"] = req.Title
		task.Title = req.Title
	}
	if req.Description != "" {
		changes["description"] = req.Description
		task.Description = req.Description
	}
	if req.AssigneeID != "" {
		changes["assignee_id"] = req.AssigneeID
		task.AssigneeID = req.AssigneeID
	}
	if req.Priority != "" {
		changes["priority"] = req.Priority
		task.Priority = req.Priority
	}
	if req.Status != "" {
		changes["status"] = req.Status
		task.Status = req.Status
	}
	if req.DueDate != nil {
		changes["due_date"] = req.DueDate
		task.DueDate = req.DueDate
	}

	err = getDB(ctx).Save(&task).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Create audit log
	createAuditLog(ctx, userID, "update", "task", task.ID, changes)

	return &task, nil
}

// DeleteTask deletes a task
func DeleteTask(ctx *context.Context, req struct{}) (*domain.MessageResponse, error) {
	taskID := getParam(ctx, "id")
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
	allowed, err := ctx.Authz().Authorize(ctx.Context(), identity, "delete", "tasks")
	if err != nil || !allowed {
		return nil, fmt.Errorf("forbidden: insufficient permissions")
	}

	var task domain.Task
	err = getDB(ctx).Model(&domain.Task{}).
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("tasks.id = ? AND projects.tenant_id = ?", taskID, tenantID).
		First(&task).Error
	if err != nil {
		return nil, fmt.Errorf("task not found")
	}

	err = getDB(ctx).Delete(&task).Error
	if err != nil {
		return nil, fmt.Errorf("failed to delete task: %w", err)
	}

	// Create audit log
	createAuditLog(ctx, userID, "delete", "task", task.ID, map[string]interface{}{
		"title": task.Title,
	})

	return &domain.MessageResponse{
		Message: "Task deleted successfully",
	}, nil
}
