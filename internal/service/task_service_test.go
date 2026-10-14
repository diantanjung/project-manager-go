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
		{"member cannot reorder", domain.AuthUser{ID: 2, Role: domain.RoleTeamMember}, TaskPatchInput{Position: &position}, domain.ErrForbidden},
		{"project cannot be changed", domain.AuthUser{ID: 1, Role: domain.RoleProjectManager}, TaskPatchInput{ProjectID: &projectID}, domain.ErrValidation},
		{"empty patch rejected", domain.AuthUser{ID: 1, Role: domain.RoleProjectManager}, TaskPatchInput{}, domain.ErrValidation},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMemoryStore(t)
			store.task = domain.Task{ID: 10, AssigneeID: 2}
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
