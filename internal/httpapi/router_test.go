package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"project-manager-go/internal/auth"
	"project-manager-go/internal/config"
	"project-manager-go/internal/domain"
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

func TestAuthMeRoutesReturnCurrentUser(t *testing.T) {
	tokenManager := tokenManagerForHTTP()
	user := domain.User{
		ID:           7,
		Name:         "Dian",
		Email:        "dian@example.com",
		PasswordHash: "stored-hash",
		Role:         domain.RoleTeamMember,
	}
	router := NewRouter(
		config.Config{NodeEnv: "test", FrontendURL: "http://localhost:5173"},
		service.New(&authHTTPStore{user: user}, tokenManager),
	)
	accessToken, err := tokenManager.GenerateAccessToken(user)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
	}{
		{name: "api alias", path: "/api/auth/me"},
		{name: "versioned api", path: "/api/v1/auth/me"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}

			var got struct {
				Data domain.User `json:"data"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decoding response: %v", err)
			}
			if got.Data.ID != user.ID || got.Data.Email != user.Email || got.Data.Role != user.Role {
				t.Fatalf("user = %+v, want id/email/role from %+v", got.Data, user)
			}
			if got.Data.PasswordHash != "" {
				t.Fatalf("password hash leaked in response")
			}
		})
	}
}

func TestDashboardRouteReturnsDirectFrontendPayload(t *testing.T) {
	tokenManager := tokenManagerForHTTP()
	user := domain.User{
		ID:           7,
		Name:         "Dian",
		Email:        "dian@example.com",
		PasswordHash: "stored-hash",
		Role:         domain.RoleTeamMember,
	}
	router := NewRouter(
		config.Config{NodeEnv: "test", FrontendURL: "http://localhost:5173"},
		service.New(&dashboardHTTPStore{
			authHTTPStore: authHTTPStore{user: user},
			summary: domain.DashboardSummary{
				TotalActiveProjects: 1,
				TaskCountPerStatus: map[domain.TaskStatus]int{
					domain.TaskStatusTodo: 1,
				},
				ActiveProgress: domain.ActiveProgress{
					Doing:      1,
					Total:      1,
					Ratio:      1,
					Percentage: 100,
				},
				RecentTasks:       []domain.Task{},
				UpcomingDeadlines: []domain.Task{},
				HighPriorityTasks: []domain.Task{},
				LatestUpdates:     []domain.ActivityLog{},
			},
		}, tokenManager),
	)
	accessToken, err := tokenManager.GenerateAccessToken(user)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if _, ok := got["data"]; ok {
		t.Fatalf("dashboard response was enveloped: %s", rec.Body.String())
	}
	if got["totalActiveProjects"] != float64(1) {
		t.Fatalf("totalActiveProjects = %v, want 1", got["totalActiveProjects"])
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

type authHTTPStore struct {
	service.Store
	user domain.User
}

func (s *authHTTPStore) GetUserByID(_ context.Context, id int) (domain.User, error) {
	if id != s.user.ID {
		return domain.User{}, domain.ErrNotFound
	}
	return s.user, nil
}

type dashboardHTTPStore struct {
	authHTTPStore
	summary domain.DashboardSummary
}

func (s *dashboardHTTPStore) DashboardSummary(context.Context, domain.AuthUser) (domain.DashboardSummary, error) {
	return s.summary, nil
}
