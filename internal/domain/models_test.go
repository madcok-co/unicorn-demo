package domain

import (
	"testing"
	"time"
)

func TestUserTableName(t *testing.T) {
	user := User{}
	if got := user.TableName(); got != "users" {
		t.Errorf("User.TableName() = %v, want %v", got, "users")
	}
}

func TestProjectTableName(t *testing.T) {
	project := Project{}
	if got := project.TableName(); got != "projects" {
		t.Errorf("Project.TableName() = %v, want %v", got, "projects")
	}
}

func TestTaskTableName(t *testing.T) {
	task := Task{}
	if got := task.TableName(); got != "tasks" {
		t.Errorf("Task.TableName() = %v, want %v", got, "tasks")
	}
}

func TestCommentTableName(t *testing.T) {
	comment := Comment{}
	if got := comment.TableName(); got != "comments" {
		t.Errorf("Comment.TableName() = %v, want %v", got, "comments")
	}
}

func TestAuditLogTableName(t *testing.T) {
	auditLog := AuditLog{}
	if got := auditLog.TableName(); got != "audit_logs" {
		t.Errorf("AuditLog.TableName() = %v, want %v", got, "audit_logs")
	}
}

func TestUserModel(t *testing.T) {
	user := User{
		ID:       "user-1",
		Email:    "test@example.com",
		Name:     "Test User",
		TenantID: "tenant-1",
		Roles:    []string{"admin", "developer"},
		Active:   true,
	}

	if user.ID != "user-1" {
		t.Errorf("User.ID = %v, want %v", user.ID, "user-1")
	}
	if user.Email != "test@example.com" {
		t.Errorf("User.Email = %v, want %v", user.Email, "test@example.com")
	}
	if len(user.Roles) != 2 {
		t.Errorf("User.Roles length = %v, want %v", len(user.Roles), 2)
	}
	if !user.Active {
		t.Error("User.Active = false, want true")
	}
}

func TestProjectModel(t *testing.T) {
	project := Project{
		ID:          "proj-1",
		Name:        "Test Project",
		Description: "Test Description",
		TenantID:    "tenant-1",
		OwnerID:     "user-1",
		Status:      "active",
		Tags:        []string{"demo", "test"},
	}

	if project.ID != "proj-1" {
		t.Errorf("Project.ID = %v, want %v", project.ID, "proj-1")
	}
	if project.Status != "active" {
		t.Errorf("Project.Status = %v, want %v", project.Status, "active")
	}
	if len(project.Tags) != 2 {
		t.Errorf("Project.Tags length = %v, want %v", len(project.Tags), 2)
	}
}

func TestTaskModel(t *testing.T) {
	dueDate := time.Now().Add(24 * time.Hour)
	task := Task{
		ID:          "task-1",
		ProjectID:   "proj-1",
		Title:       "Test Task",
		Description: "Test Description",
		AssigneeID:  "user-1",
		Priority:    "high",
		Status:      "in_progress",
		DueDate:     &dueDate,
	}

	if task.ID != "task-1" {
		t.Errorf("Task.ID = %v, want %v", task.ID, "task-1")
	}
	if task.Priority != "high" {
		t.Errorf("Task.Priority = %v, want %v", task.Priority, "high")
	}
	if task.Status != "in_progress" {
		t.Errorf("Task.Status = %v, want %v", task.Status, "in_progress")
	}
	if task.DueDate == nil {
		t.Error("Task.DueDate is nil, want non-nil")
	}
}

func TestCommentModel(t *testing.T) {
	comment := Comment{
		ID:      "comment-1",
		TaskID:  "task-1",
		UserID:  "user-1",
		Content: "This is a test comment",
	}

	if comment.ID != "comment-1" {
		t.Errorf("Comment.ID = %v, want %v", comment.ID, "comment-1")
	}
	if comment.Content != "This is a test comment" {
		t.Errorf("Comment.Content = %v, want %v", comment.Content, "This is a test comment")
	}
}

func TestAuditLogModel(t *testing.T) {
	auditLog := AuditLog{
		ID:         "audit-1",
		TenantID:   "tenant-1",
		UserID:     "user-1",
		Action:     "create",
		Resource:   "project",
		ResourceID: "proj-1",
		Changes: map[string]interface{}{
			"name":   "New Project",
			"status": "active",
		},
		IPAddress: "127.0.0.1",
		UserAgent: "Mozilla/5.0",
	}

	if auditLog.Action != "create" {
		t.Errorf("AuditLog.Action = %v, want %v", auditLog.Action, "create")
	}
	if auditLog.Resource != "project" {
		t.Errorf("AuditLog.Resource = %v, want %v", auditLog.Resource, "project")
	}
	if len(auditLog.Changes) != 2 {
		t.Errorf("AuditLog.Changes length = %v, want %v", len(auditLog.Changes), 2)
	}
}
