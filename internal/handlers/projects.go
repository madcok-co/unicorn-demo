package handlers

import (
	"fmt"
	"time"

	"github.com/madcok-co/unicorn-demo/internal/domain"
	"github.com/madcok-co/unicorn/contrib/pagination"
	"github.com/madcok-co/unicorn/core/pkg/context"
)

// ListProjects returns paginated list of projects
func ListProjects(ctx *context.Context, req domain.ListProjectsRequest) (*pagination.OffsetResult, error) {
	db := getDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
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
	query := db.Where("tenant_id = ?", tenantID)
	
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Get total count
	var total int64
	query.Model(&domain.Project{}).Count(&total)

	// Get paginated results
	var projects []domain.Project
	offset := (req.Page - 1) * req.Limit
	query.Order(fmt.Sprintf("%s %s", req.Sort, req.Order)).
		Limit(req.Limit).
		Offset(offset).
		Find(&projects)

	params := pagination.ParseOffsetParams(req.Page, req.Limit, req.Sort, req.Order)
	return pagination.NewOffsetResult(projects, total, params), nil
}

// CreateProject creates a new project
func CreateProject(ctx *context.Context, req domain.CreateProjectRequest) (*domain.Project, error) {
	db := getDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	project := &domain.Project{
		ID:          fmt.Sprintf("proj-%d", time.Now().UnixNano()),
		Name:        req.Name,
		Description: req.Description,
		TenantID:    tenantID,
		OwnerID:     userID,
		Status:      "active",
		Tags:        req.Tags,
	}

	err = db.Create(project).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Create audit log
	createAuditLog(ctx, userID, "create", "project", project.ID, map[string]interface{}{
		"name": project.Name,
	})

	return project, nil
}

// GetProject returns a single project by ID
func GetProject(ctx *context.Context, req struct{}) (*domain.Project, error) {
	db := getDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	projectID := getParam(ctx, "id")
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	var project domain.Project
	err = db.Where("id = ? AND tenant_id = ?", projectID, tenantID).First(&project).Error
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	return &project, nil
}

// UpdateProject updates an existing project
func UpdateProject(ctx *context.Context, req domain.UpdateProjectRequest) (*domain.Project, error) {
	db := getDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	projectID := getParam(ctx, "id")
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	var project domain.Project
	err = db.Where("id = ? AND tenant_id = ?", projectID, tenantID).First(&project).Error
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	// Track changes
	changes := make(map[string]interface{})
	
	if req.Name != "" {
		changes["name"] = req.Name
		project.Name = req.Name
	}
	if req.Description != "" {
		changes["description"] = req.Description
		project.Description = req.Description
	}
	if req.Status != "" {
		changes["status"] = req.Status
		project.Status = req.Status
	}
	if len(req.Tags) > 0 {
		changes["tags"] = req.Tags
		project.Tags = req.Tags
	}

	err = db.Save(&project).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	// Create audit log
	createAuditLog(ctx, userID, "update", "project", project.ID, changes)

	return &project, nil
}

// DeleteProject deletes a project
func DeleteProject(ctx *context.Context, req struct{}) (*domain.MessageResponse, error) {
	db := getDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	projectID := getParam(ctx, "id")
	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	var project domain.Project
	err = db.Where("id = ? AND tenant_id = ?", projectID, tenantID).First(&project).Error
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	err = db.Delete(&project).Error
	if err != nil {
		return nil, fmt.Errorf("failed to delete project: %w", err)
	}

	// Create audit log
	createAuditLog(ctx, userID, "delete", "project", project.ID, map[string]interface{}{
		"name": project.Name,
	})

	return &domain.MessageResponse{
		Message: "Project deleted successfully",
	}, nil
}

// BatchUpdateProjects updates multiple projects at once (V2 feature)
func BatchUpdateProjects(ctx *context.Context, req domain.BatchUpdateProjectsRequest) (map[string]interface{}, error) {
	db := getDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}

	updated := 0
	for _, projectID := range req.ProjectIDs {
		var project domain.Project
		err := db.Where("id = ? AND tenant_id = ?", projectID, tenantID).First(&project).Error
		if err != nil {
			continue
		}

		// Apply updates
		if name, ok := req.Updates["name"].(string); ok {
			project.Name = name
		}
		if desc, ok := req.Updates["description"].(string); ok {
			project.Description = desc
		}
		if status, ok := req.Updates["status"].(string); ok {
			project.Status = status
		}

		db.Save(&project)
		updated++

		// Create audit log
		createAuditLog(ctx, userID, "batch_update", "project", project.ID, req.Updates)
	}

	return map[string]interface{}{
		"updated": updated,
		"total":   len(req.ProjectIDs),
	}, nil
}

// GetProjectStats returns statistics for projects
func GetProjectStats(ctx *context.Context, req struct{}) (*domain.ProjectStatsResponse, error) {
	db := getDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database not available")
	}

	tenantID, err := getTenantID(ctx)
	if err != nil {
		return nil, err
	}

	stats := &domain.ProjectStatsResponse{
		ProjectsByStatus: make(map[string]int64),
		TasksByPriority:  make(map[string]int64),
	}

	// Total projects
	db.Model(&domain.Project{}).Where("tenant_id = ?", tenantID).Count(&stats.TotalProjects)

	// Active projects
	db.Model(&domain.Project{}).Where("tenant_id = ? AND status = ?", tenantID, "active").Count(&stats.ActiveProjects)

	// Projects by status
	type StatusCount struct {
		Status string
		Count  int64
	}
	var statusCounts []StatusCount
	db.Model(&domain.Project{}).
		Where("tenant_id = ?", tenantID).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusCounts)
	
	for _, sc := range statusCounts {
		stats.ProjectsByStatus[sc.Status] = sc.Count
	}

	// Total tasks
	db.Model(&domain.Task{}).
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("projects.tenant_id = ?", tenantID).
		Count(&stats.TotalTasks)

	// Completed tasks
	db.Model(&domain.Task{}).
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("projects.tenant_id = ? AND tasks.status = ?", tenantID, "done").
		Count(&stats.CompletedTasks)

	// Tasks by priority
	type PriorityCount struct {
		Priority string
		Count    int64
	}
	var priorityCounts []PriorityCount
	db.Model(&domain.Task{}).
		Joins("JOIN projects ON tasks.project_id = projects.id").
		Where("projects.tenant_id = ?", tenantID).
		Select("tasks.priority, count(*) as count").
		Group("tasks.priority").
		Scan(&priorityCounts)
	
	for _, pc := range priorityCounts {
		stats.TasksByPriority[pc.Priority] = pc.Count
	}

	// Recent activity
	db.Where("tenant_id = ?", tenantID).
		Order("created_at DESC").
		Limit(10).
		Find(&stats.RecentActivity)

	return stats, nil
}
