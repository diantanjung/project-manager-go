package service

import (
	"context"
	"errors"
	"testing"

	"project-manager-go/internal/domain"
)

func TestCreateCommentCreatesMentionNotification(t *testing.T) {
	store := newMemoryStore(t)
	svc := New(store, testTokenManager())
	actor := domain.AuthUser{ID: 1, Role: domain.RoleTeamMember}

	_, err := svc.CreateComment(context.Background(), actor, CommentInput{
		TaskID:  10,
		Content: "please review @Dian @Dian @Nobody",
	})
	if err != nil {
		t.Fatalf("CreateComment() error = %v", err)
	}

	if len(store.notifications) != 1 {
		t.Fatalf("notifications = %d, want 1", len(store.notifications))
	}
	if store.notifications[0].UserID != 2 || store.notifications[0].Type != domain.NotificationMention {
		t.Fatalf("unexpected notification: %#v", store.notifications[0])
	}
}

func TestCreateCommentRequiresTaskAccess(t *testing.T) {
	store := newMemoryStore(t)
	store.canAccessTask = false
	svc := New(store, testTokenManager())

	_, err := svc.CreateComment(context.Background(), domain.AuthUser{ID: 1, Role: domain.RoleTeamMember}, CommentInput{
		TaskID:  10,
		Content: "blocked",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUpdateCommentRequiresAuthor(t *testing.T) {
	store := newMemoryStore(t)
	svc := New(store, testTokenManager())

	_, err := svc.UpdateComment(context.Background(), domain.AuthUser{ID: 2}, 1, "updated")
	if err == nil {
		t.Fatal("expected error")
	}
	var appErr *domain.AppError
	if !errors.As(err, &appErr) || !errors.Is(appErr.Err, domain.ErrForbidden) {
		t.Fatalf("got %v, want forbidden app error", err)
	}
}

func TestDeleteAttachmentRequiresUploader(t *testing.T) {
	store := newMemoryStore(t)
	svc := New(store, testTokenManager())

	_, err := svc.DeleteAttachment(context.Background(), domain.AuthUser{ID: 2}, 1)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetAttachmentRequiresTaskAccess(t *testing.T) {
	store := newMemoryStore(t)
	store.canAccessTask = false
	svc := New(store, testTokenManager())

	_, err := svc.GetAttachmentByID(context.Background(), domain.AuthUser{ID: 1, Role: domain.RoleTeamMember}, 1)
	if err == nil {
		t.Fatal("expected error")
	}
}
