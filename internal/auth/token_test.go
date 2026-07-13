package auth

import (
	"testing"
	"time"

	"project-manager-go/internal/domain"
)

func TestTokenManagerGenerateAndVerifyAccessToken(t *testing.T) {
	now := time.Date(2026, 7, 13, 8, 0, 0, 0, time.UTC)
	manager := TokenManager{
		AccessSecret:  []byte("access"),
		RefreshSecret: []byte("refresh"),
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
		Now:           func() time.Time { return now },
	}

	user := domain.User{
		ID:    42,
		Email: "dian@example.com",
		Role:  domain.RoleProjectManager,
	}

	token, err := manager.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	got, err := manager.VerifyAccessToken(token)
	if err != nil {
		t.Fatalf("VerifyAccessToken() error = %v", err)
	}
	if got.ID != user.ID || got.Email != user.Email || got.Role != user.Role {
		t.Fatalf("got %#v, want user id/email/role from token", got)
	}
}

func TestHashRefreshToken(t *testing.T) {
	got := HashRefreshToken("token")
	if got == "token" {
		t.Fatal("hash should not return raw token")
	}
	if len(got) != 64 {
		t.Fatalf("got hash length %d, want 64", len(got))
	}
}
