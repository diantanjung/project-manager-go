package service

import (
	"context"
	"errors"
	"testing"

	"project-manager-go/internal/domain"
)

func TestCreateTeamValidatesName(t *testing.T) {
	svc := New(newMemoryStore(t), testTokenManager())

	_, err := svc.CreateTeam(context.Background(), domain.AuthUser{ID: 1, Role: domain.RoleProductOwner}, TeamInput{})

	assertValidationError(t, err)
}

func TestCreateProjectValidatesTeamID(t *testing.T) {
	svc := New(newMemoryStore(t), testTokenManager())
	teamID := 0

	_, err := svc.CreateProject(context.Background(), domain.AuthUser{ID: 1, Role: domain.RoleProjectManager}, ProjectInput{
		Name:   stringPtr("Roadmap"),
		TeamID: &teamID,
	})

	assertValidationError(t, err)
}

func TestCreateTaskValidatesEnumValues(t *testing.T) {
	svc := New(newMemoryStore(t), testTokenManager())
	status := domain.TaskStatus("shipped")

	_, err := svc.CreateTask(context.Background(), domain.AuthUser{ID: 1, Role: domain.RoleProjectManager}, TaskInput{
		Title:      "Task",
		Status:     &status,
		ProjectID:  10,
		AssigneeID: 2,
	})

	assertValidationError(t, err)
}

func TestUpdateCommentValidatesContent(t *testing.T) {
	svc := New(newMemoryStore(t), testTokenManager())

	_, err := svc.UpdateComment(context.Background(), domain.AuthUser{ID: 1}, 1, " ")

	assertValidationError(t, err)
}

func assertValidationError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected validation error")
	}
	var appErr *domain.AppError
	if !errors.As(err, &appErr) || !errors.Is(appErr.Err, domain.ErrValidation) {
		t.Fatalf("got %v, want validation app error", err)
	}
}
