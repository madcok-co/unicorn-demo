package domain

import "time"

// ============================================================================
// Authentication DTOs
// ============================================================================

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         User      `json:"user"`
}

type OAuth2CallbackRequest struct {
	Code  string `json:"code" query:"code" validate:"required"`
	State string `json:"state" query:"state" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ============================================================================
// User DTOs
// ============================================================================

type CreateUserRequest struct {
	Email string   `json:"email" validate:"required,email"`
	Name  string   `json:"name" validate:"required,min=2,max=100"`
	Roles []string `json:"roles" validate:"required,min=1"`
}

type UpdateUserRequest struct {
	Name   string   `json:"name" validate:"omitempty,min=2,max=100"`
	Roles  []string `json:"roles" validate:"omitempty,min=1"`
	Active *bool    `json:"active" validate:"omitempty"`
}

type ListUsersRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort   string `query:"sort" validate:"omitempty,oneof=created_at updated_at name email"`
	Order  string `query:"order" validate:"omitempty,oneof=asc desc"`
	Role   string `query:"role"`
	Active *bool  `query:"active"`
}

// ============================================================================
// Project DTOs
// ============================================================================

type CreateProjectRequest struct {
	Name        string   `json:"name" validate:"required,min=3,max=100"`
	Description string   `json:"description" validate:"max=500"`
	Tags        []string `json:"tags" validate:"max=10"`
}

type UpdateProjectRequest struct {
	Name        string   `json:"name" validate:"omitempty,min=3,max=100"`
	Description string   `json:"description" validate:"omitempty,max=500"`
	Status      string   `json:"status" validate:"omitempty,oneof=active inactive archived"`
	Tags        []string `json:"tags" validate:"omitempty,max=10"`
}

type ListProjectsRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort   string `query:"sort" validate:"omitempty,oneof=created_at updated_at name"`
	Order  string `query:"order" validate:"omitempty,oneof=asc desc"`
	Status string `query:"status" validate:"omitempty,oneof=active inactive archived"`
	Cursor string `query:"cursor"`
}

type BatchUpdateProjectsRequest struct {
	ProjectIDs []string               `json:"project_ids" validate:"required,min=1,max=100"`
	Updates    map[string]interface{} `json:"updates" validate:"required"`
}

// ============================================================================
// Task DTOs
// ============================================================================

type CreateTaskRequest struct {
	ProjectID   string     `json:"project_id" validate:"required"`
	Title       string     `json:"title" validate:"required,min=3,max=200"`
	Description string     `json:"description" validate:"max=1000"`
	AssigneeID  string     `json:"assignee_id"`
	Priority    string     `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title" validate:"omitempty,min=3,max=200"`
	Description string     `json:"description" validate:"omitempty,max=1000"`
	AssigneeID  string     `json:"assignee_id"`
	Priority    string     `json:"priority" validate:"omitempty,oneof=low medium high urgent"`
	Status      string     `json:"status" validate:"omitempty,oneof=todo in_progress done"`
	DueDate     *time.Time `json:"due_date"`
}

type ListTasksRequest struct {
	Page       int    `query:"page" validate:"omitempty,min=1"`
	Limit      int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort       string `query:"sort" validate:"omitempty,oneof=created_at updated_at due_date priority"`
	Order      string `query:"order" validate:"omitempty,oneof=asc desc"`
	ProjectID  string `query:"project_id"`
	AssigneeID string `query:"assignee_id"`
	Status     string `query:"status" validate:"omitempty,oneof=todo in_progress done"`
	Priority   string `query:"priority" validate:"omitempty,oneof=low medium high urgent"`
}

// ============================================================================
// Comment DTOs
// ============================================================================

type CreateCommentRequest struct {
	TaskID  string `json:"task_id" validate:"required"`
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

type ListCommentsRequest struct {
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Sort   string `query:"sort" validate:"omitempty,oneof=created_at updated_at"`
	Order  string `query:"order" validate:"omitempty,oneof=asc desc"`
	TaskID string `query:"task_id" validate:"required"`
}

// ============================================================================
// Analytics DTOs
// ============================================================================

type ProjectStatsResponse struct {
	TotalProjects    int64            `json:"total_projects"`
	ActiveProjects   int64            `json:"active_projects"`
	TotalTasks       int64            `json:"total_tasks"`
	CompletedTasks   int64            `json:"completed_tasks"`
	ProjectsByStatus map[string]int64 `json:"projects_by_status"`
	TasksByPriority  map[string]int64 `json:"tasks_by_priority"`
	RecentActivity   []AuditLog       `json:"recent_activity"`
}

type UserStatsResponse struct {
	TotalUsers    int64            `json:"total_users"`
	ActiveUsers   int64            `json:"active_users"`
	UsersByRole   map[string]int64 `json:"users_by_role"`
	RecentSignups []User           `json:"recent_signups"`
}

// ============================================================================
// Generic Response DTOs
// ============================================================================

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string                 `json:"error"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}
