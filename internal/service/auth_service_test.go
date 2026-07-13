package service

import (
	"context"
	"testing"
)

func TestLoginRejectsBadPassword(t *testing.T) {
	store := newMemoryStore(t)
	svc := New(store, testTokenManager())

	_, err := svc.Login(context.Background(), LoginInput{
		Email:    "dian@example.com",
		Password: "wrong-password",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}
