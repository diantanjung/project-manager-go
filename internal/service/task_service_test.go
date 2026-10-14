package service

import (
	"context"
	"errors"
	"testing"

	"project-manager-go/internal/domain"
)

func TestUpdateTaskFieldPermissions(t *testing.T) {
	title := "Updated task"
	status := domain.TaskStatusInProgress
	position := 0
	projectID := 20
	tests := []struct {
		name      string
		actor     domain.AuthUser
		input     TaskPatchInput
		wantError error
	}{
		{"manager edits metadata", domain.AuthUser{ID: 1, Role: domain.RoleProjectManager}, TaskPatchInput{Title: &title}, nil},
		{"assignee changes status", domain.AuthUser{ID: 2, Role: domain.RoleTeamMember}, TaskPatchInput{Status: &status}, nil},
		{"other member cannot change status", domain.AuthUser{ID: 3, Role: domain.RoleTeamMember}, TaskPatchInput{Status: &status}, domain.ErrForbidden},
		{"member cannot edit title", domain.AuthUser{ID: 2, Role: domain.RoleTeamMember}, TaskPatchInput{Title: &title}, domain.ErrForbidden},
		{"member cannot clear assignee", domain.AuthUser{ID: 2, Role: domain.RoleTeamMember}, TaskPatchInput{ClearAssignee: true}, domain.ErrForbidden},
		{"manager clears assignee", domain.AuthUser{ID: 1, Role: domain.RoleProjectManager}, TaskPatchInput{ClearAssignee: true}, nil},
		{"member cannot reorder", domain.AuthUser{ID: 2, Role: domain.RoleTeamMember}, TaskPatchInput{Position: &position}, domain.ErrForbidden},
		{"project cannot be changed", domain.AuthUser{ID: 1, Role: domain.RoleProjectManager}, TaskPatchInput{ProjectID: &projectID}, domain.ErrValidation},
		{"empty patch rejected", domain.AuthUser{ID: 1, Role: domain.RoleProjectManager}, TaskPatchInput{}, domain.ErrValidation},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMemoryStore(t)
			assigneeID := 2
			store.task = domain.Task{ID: 10, AssigneeID: &assigneeID}
			_, err := New(store, testTokenManager()).UpdateTask(context.Background(), tt.actor, 10, tt.input)
			if !errors.Is(err, tt.wantError) {
				t.Fatalf("UpdateTask() error = %v, want %v", err, tt.wantError)
			}
			wantUpdates := 1
			if tt.wantError != nil {
				wantUpdates = 0
			}
			if len(store.taskUpdates) != wantUpdates {
				t.Fatalf("store updates = %d, want %d", len(store.taskUpdates), wantUpdates)
			}
		})
	}
}

func TestCreateTaskAppliesDefaults(t *testing.T) {
	store := newMemoryStore(t)
	svc := New(store, testTokenManager())

	task, err := svc.CreateTask(context.Background(), domain.AuthUser{
		ID:   1,
		Role: domain.RoleProjectManager,
	}, TaskInput{
		Title:      "Task",
		ProjectID:  10,
		AssigneeID: 2,
	})
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}
	if task.Status != domain.TaskStatusBacklog {
		t.Fatalf("status = %q, want backlog", task.Status)
	}
	if task.Priority == nil || *task.Priority != domain.TaskPriorityMedium {
		t.Fatalf("priority = %#v, want medium", task.Priority)
	}
}

func TestCreateTaskAllowsUnassignedTask(t *testing.T) {
	store := newMemoryStore(t)
	task, err := New(store, testTokenManager()).CreateTask(context.Background(), domain.AuthUser{
		ID: 1, Role: domain.RoleProjectManager,
	}, TaskInput{Title: "Unassigned task", ProjectID: 10})
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}
	if task.AssigneeID != nil {
		t.Fatalf("assigneeId = %v, want nil", task.AssigneeID)
	}
}

func TestAssignUserToTaskCreatesNotification(t *testing.T) {
	store := newMemoryStore(t)
	svc := New(store, testTokenManager())

	_, exists, err := svc.AssignUserToTask(context.Background(), domain.AuthUser{
		ID:   1,
		Role: domain.RoleProjectManager,
	}, TaskAssignmentInput{
		TaskID: 10,
		UserID: 2,
	})
	if err != nil {
		t.Fatalf("AssignUserToTask() error = %v", err)
	}
	if exists {
		t.Fatal("assignment should be new")
	}
	if len(store.notifications) != 1 || store.notifications[0].Type != domain.NotificationTaskAssigned {
		t.Fatalf("notifications = %#v, want task_assigned", store.notifications)
	}
}

func TestReorderTasksValidatesBatch(t *testing.T) {
	store := newMemoryStore(t)
	svc := New(store, testTokenManager())
	actor := domain.AuthUser{ID: 1, Role: domain.RoleProjectManager}
	status := domain.TaskStatusTodo
	tests := []struct {
		name  string
		input TaskReorderInput
		want  error
	}{
		{"valid batch", TaskReorderInput{ProjectID: 10, Moves: []TaskMove{{TaskID: 1, Status: status, Position: 0}}}, nil},
		{"empty moves", TaskReorderInput{ProjectID: 10}, domain.ErrValidation},
		{"duplicate task", TaskReorderInput{ProjectID: 10, Moves: []TaskMove{{TaskID: 1, Status: status}, {TaskID: 1, Status: status}}}, domain.ErrValidation},
		{"negative position", TaskReorderInput{ProjectID: 10, Moves: []TaskMove{{TaskID: 1, Status: status, Position: -1}}}, domain.ErrValidation},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := len(store.taskReorders)
			_, err := svc.ReorderTasks(context.Background(), actor, tt.input)
			if !errors.Is(err, tt.want) {
				t.Fatalf("ReorderTasks() error = %v, want %v", err, tt.want)
			}
			if got := len(store.taskReorders) - before; got != boolInt(tt.want == nil) {
				t.Fatalf("store reorder calls = %d", got)
			}
		})
	}
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func TestListTasksRejectsInvalidDate(t *testing.T) {
	store := newMemoryStore(t)
	bad := "2026-99-99"
	_, err := New(store, testTokenManager()).ListTasks(context.Background(), domain.AuthUser{ID: 1}, TaskFilter{DueBefore: &bad})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("ListTasks() error = %v, want validation error", err)
	}
}
