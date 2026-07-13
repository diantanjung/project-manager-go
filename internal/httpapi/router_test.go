package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"project-manager-go/internal/auth"
	"project-manager-go/internal/config"
	"project-manager-go/internal/service"
)

func TestHealthRoute(t *testing.T) {
	router := NewRouter(config.Config{NodeEnv: "test", FrontendURL: "http://localhost:5173"}, service.New(newMemoryStoreForHTTP(), tokenManagerForHTTP()))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestProtectedRouteRequiresToken(t *testing.T) {
	router := NewRouter(config.Config{NodeEnv: "test", FrontendURL: "http://localhost:5173"}, service.New(newMemoryStoreForHTTP(), tokenManagerForHTTP()))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func tokenManagerForHTTP() auth.TokenManager {
	return auth.TokenManager{
		AccessSecret:  []byte("access"),
		RefreshSecret: []byte("refresh"),
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
		Now:           time.Now,
	}
}

func newMemoryStoreForHTTP() service.Store {
	return &noopStore{}
}

type noopStore struct {
	service.Store
}
