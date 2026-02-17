package domain

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	TenantID  string    `json:"tenant_id" gorm:"index;not null"`
	Roles     []string  `json:"roles" gorm:"serializer:json"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Project represents a project in the system
type Project struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	TenantID    string    `json:"tenant_id" gorm:"index;not null"`
	OwnerID     string    `json:"owner_id" gorm:"index;not null"`
	Status      string    `json:"status" gorm:"default:'active'"` // active, inactive, archived
	Tags        []string  `json:"tags" gorm:"serializer:json"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Task represents a task within a project
type Task struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	ProjectID   string     `json:"project_id" gorm:"index;not null"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	AssigneeID  string     `json:"assignee_id" gorm:"index"`
	Priority    string     `json:"priority" gorm:"default:'medium'"` // low, medium, high, urgent
	Status      string     `json:"status" gorm:"default:'todo'"`     // todo, in_progress, done
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Comment represents a comment on a task
type Comment struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	TaskID    string    `json:"task_id" gorm:"index;not null"`
	UserID    string    `json:"user_id" gorm:"index;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID         string                 `json:"id" gorm:"primaryKey"`
	TenantID   string                 `json:"tenant_id" gorm:"index;not null"`
	UserID     string                 `json:"user_id" gorm:"index;not null"`
	Action     string                 `json:"action" gorm:"not null"`   // create, update, delete, read
	Resource   string                 `json:"resource" gorm:"not null"` // user, project, task, comment
	ResourceID string                 `json:"resource_id" gorm:"index"`
	Changes    map[string]interface{} `json:"changes" gorm:"serializer:json"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	CreatedAt  time.Time              `json:"created_at" gorm:"index"`
}

// TableName overrides
func (User) TableName() string {
	return "users"
}

func (Project) TableName() string {
	return "projects"
}

func (Task) TableName() string {
	return "tasks"
}

func (Comment) TableName() string {
	return "comments"
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
