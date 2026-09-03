package service

import (
	"context"
	"errors"
	"testing"

	"project-manager-go/internal/domain"
)

func TestAddTeamMemberRequiresScopedTeamAccess(t *testing.T) {
	store := newMemoryStore(t)
	store.canCreateProject = false
	svc := New(store, testTokenManager())

	_, err := svc.AddTeamMember(context.Background(), domain.AuthUser{
		ID:   1,
		Role: domain.RoleProjectManager,
	}, 10, TeamMemberInput{
		UserID: 2,
	})
	if err == nil {
		t.Fatal("expected error")
	}

	var appErr *domain.AppError
	if !errors.As(err, &appErr) || !errors.Is(appErr.Err, domain.ErrForbidden) {
		t.Fatalf("got %v, want forbidden app error", err)
	}
}
