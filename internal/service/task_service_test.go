package service

import (
	"context"
	"testing"

	"project-manager-go/internal/domain"
)

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
