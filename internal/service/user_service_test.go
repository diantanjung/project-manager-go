package service

import (
	"context"
	"errors"
	"testing"

	"project-manager-go/internal/domain"
)

func TestUpdateUserRequiresAdminForOtherUsers(t *testing.T) {
	store := newMemoryStore(t)
	svc := New(store, testTokenManager())

	password := "new-secret"
	_, err := svc.UpdateUser(context.Background(), domain.AuthUser{
		ID:   1,
		Role: domain.RoleProjectManager,
	}, 2, UpdateUserInput{
		Password: &password,
	})
	if err == nil {
		t.Fatal("expected error")
	}

	var appErr *domain.AppError
	if !errors.As(err, &appErr) || !errors.Is(appErr.Err, domain.ErrForbidden) {
		t.Fatalf("got %v, want forbidden app error", err)
	}
}

func TestListUsersRequiresProjectManager(t *testing.T) {
	store := newMemoryStore(t)
	svc := New(store, testTokenManager())

	_, err := svc.ListUsers(context.Background(), domain.AuthUser{
		ID:   1,
		Role: domain.RoleTeamMember,
	}, ListUsersFilter{})
	if err == nil {
		t.Fatal("expected error")
	}

	var appErr *domain.AppError
	if !errors.As(err, &appErr) || !errors.Is(appErr.Err, domain.ErrForbidden) {
		t.Fatalf("got %v, want forbidden app error", err)
	}
}

func TestGetUserAllowsSelf(t *testing.T) {
	store := newMemoryStore(t)
	svc := New(store, testTokenManager())

	user, err := svc.GetUserByID(context.Background(), domain.AuthUser{
		ID:   1,
		Role: domain.RoleTeamMember,
	}, 1)
	if err != nil {
		t.Fatalf("GetUserByID() error = %v", err)
	}
	if user.ID != 1 {
		t.Fatalf("user ID = %d, want 1", user.ID)
	}
}
