package service

import (
	"context"
	"testing"

	"project-manager-go/internal/domain"
)

func TestCreateProjectRequiresAllowedTeam(t *testing.T) {
	store := newMemoryStore(t)
	store.canCreateProject = false
	svc := New(store, testTokenManager())

	teamID := 99
	_, err := svc.CreateProject(context.Background(), domain.AuthUser{
		ID:   1,
		Role: domain.RoleProjectManager,
	}, ProjectInput{
		Name:   stringPtr("Project"),
		TeamID: &teamID,
	})
	if err == nil {
		t.Fatal("expected error")
	}
}
