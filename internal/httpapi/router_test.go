package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

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

func TestRefreshAndLogoutWithoutCookie(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	store := &authSessionStore{
		user:   domain.User{ID: 7, Email: "dian@example.com", PasswordHash: string(passwordHash), Role: domain.RoleTeamMember},
		tokens: make(map[string]bool),
	}
	router := NewRouter(config.Config{NodeEnv: "production", FrontendURL: "https://taskflowprojectmanager.vercel.app", JWTRefreshExpiresIn: 7 * 24 * time.Hour}, service.New(store, tokenManagerForHTTP()))
	request := func(path, body, accessToken string) *httptest.ResponseRecorder {
		t.Helper()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://taskflowprojectmanager.vercel.app")
		if accessToken != "" {
			req.Header.Set("Authorization", "Bearer "+accessToken)
		}
		router.ServeHTTP(rec, req)
		return rec
	}
	decode := func(rec *httptest.ResponseRecorder) authTokenResponse {
		t.Helper()
		var response authTokenResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		return response
	}

	login := request("/api/v1/auth/login", `{"email":"dian@example.com","password":"secret123"}`, "")
	if login.Code != http.StatusOK {
		t.Fatalf("login status = %d: %s", login.Code, login.Body.String())
	}
	initial := decode(login).Data
	if initial.AccessToken == "" || initial.RefreshToken == "" {
		t.Fatal("login response must include both tokens")
	}
	if cookie := login.Result().Cookies(); len(cookie) == 0 || !cookie[0].HttpOnly || !cookie[0].Secure {
		t.Fatalf("login cookie = %#v, want Secure HttpOnly", cookie)
	}

	refresh := request("/api/v1/auth/refresh", `{"refreshToken":"`+initial.RefreshToken+`"}`, "")
	if refresh.Code != http.StatusOK {
		t.Fatalf("refresh status = %d: %s", refresh.Code, refresh.Body.String())
	}
	rotated := decode(refresh).Data
	if rotated.AccessToken == "" || rotated.RefreshToken == "" || rotated.RefreshToken == initial.RefreshToken {
		t.Fatalf("invalid rotated tokens: %+v", rotated)
	}
	if store.tokens[auth.HashRefreshToken(initial.RefreshToken)] {
		t.Fatal("initial refresh token remains valid")
	}

	logout := request("/api/v1/auth/logout", `{"refreshToken":"`+rotated.RefreshToken+`"}`, rotated.AccessToken)
	if logout.Code != http.StatusNoContent {
		t.Fatalf("logout status = %d: %s", logout.Code, logout.Body.String())
	}
	if store.tokens[auth.HashRefreshToken(rotated.RefreshToken)] {
		t.Fatal("rotated refresh token remains valid after logout")
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

func TestPatchTaskRouteAcceptsPositionZero(t *testing.T) {
	user := domain.User{ID: 7, Email: "manager@example.com", Role: domain.RoleProjectManager}
	store := &taskPatchHTTPStore{authHTTPStore: authHTTPStore{user: user}}
	tokens := tokenManagerForHTTP()
	router := NewRouter(config.Config{NodeEnv: "test", FrontendURL: "http://localhost:5173"}, service.New(store, tokens))
	token, err := tokens.GenerateAccessToken(user)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/10", strings.NewReader(`{"title":"Edited task","position":0}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH status = %d: %s", rec.Code, rec.Body.String())
	}
	if store.patch.Position == nil || *store.patch.Position != 0 || store.patch.Title == nil || *store.patch.Title != "Edited task" {
		t.Fatalf("stored patch = %#v, want title and position 0", store.patch)
	}
}

func TestTaskBoardAndListResponses(t *testing.T) {
	user := domain.User{ID: 7, Email: "manager@example.com", Role: domain.RoleProjectManager}
	store := &taskManagementHTTPStore{authHTTPStore: authHTTPStore{user: user}}
	tokens := tokenManagerForHTTP()
	router := NewRouter(config.Config{NodeEnv: "test", FrontendURL: "http://localhost:5173"}, service.New(store, tokens))
	token, err := tokens.GenerateAccessToken(user)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name           string
		path           string
		wantPagination bool
	}{
		{"board", "/api/v1/tasks?project_id=9&dueAfter=2026-10-01", false},
		{"list", "/api/v1/tasks?projectId=9&page=1&limit=20", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			router.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d: %s", rec.Code, rec.Body.String())
			}
			var body map[string]json.RawMessage
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			_, hasPagination := body["pagination"]
			if hasPagination != tt.wantPagination {
				t.Fatalf("pagination present = %t, want %t: %s", hasPagination, tt.wantPagination, rec.Body.String())
			}
			if store.filter.ProjectID == nil || *store.filter.ProjectID != 9 || store.filter.Board == tt.wantPagination {
				t.Fatalf("filter = %#v", store.filter)
			}
		})
	}
}

func TestTaskManagementRoutes(t *testing.T) {
	user := domain.User{ID: 7, Email: "manager@example.com", Role: domain.RoleProjectManager}
	store := &taskManagementHTTPStore{authHTTPStore: authHTTPStore{user: user}}
	tokens := tokenManagerForHTTP()
	router := NewRouter(config.Config{NodeEnv: "test", FrontendURL: "http://localhost:5173"}, service.New(store, tokens))
	token, err := tokens.GenerateAccessToken(user)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		method string
		path   string
		body   string
		want   int
	}{
		{http.MethodPost, "/api/v1/tasks/10/assignments", `{"userId":8}`, http.StatusCreated},
		{http.MethodDelete, "/api/v1/tasks/10/assignments/8", "", http.StatusOK},
		{http.MethodPost, "/api/v1/tasks/reorder", `{"projectId":9,"moves":[{"taskId":10,"status":"todo","position":0}]}`, http.StatusOK},
		{http.MethodGet, "/api/v1/tasks/10/activity", "", http.StatusOK},
	} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(rec, req)
		if rec.Code != tt.want {
			t.Fatalf("%s %s status = %d, want %d: %s", tt.method, tt.path, rec.Code, tt.want, rec.Body.String())
		}
	}
}

func TestCreateTaskAcceptsFrontendPayload(t *testing.T) {
	user := domain.User{ID: 7, Email: "manager@example.com", Role: domain.RoleProjectManager}
	store := &taskManagementHTTPStore{authHTTPStore: authHTTPStore{user: user}}
	tokens := tokenManagerForHTTP()
	router := NewRouter(config.Config{NodeEnv: "test", FrontendURL: "http://localhost:5173"}, service.New(store, tokens))
	token, err := tokens.GenerateAccessToken(user)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", strings.NewReader(`{
		"title":"Backlog task", "description":"", "status":"backlog", "priority":"medium",
		"project_id":9, "creator_id":7, "assignee_id":0,
		"due_date":"2026-10-20T00:00:00.000Z"
	}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /tasks status = %d, want 201: %s", rec.Code, rec.Body.String())
	}
	if store.createdTask.ProjectID != 9 || store.createdTask.CreatorID != 7 || store.createdTask.AssigneeID != 0 || store.createdTask.DueDate == nil || *store.createdTask.DueDate != "2026-10-20" {
		t.Fatalf("created task input = %#v", store.createdTask)
	}
}

func TestUpdateTaskAcceptsFrontendPayload(t *testing.T) {
	user := domain.User{ID: 7, Email: "manager@example.com", Role: domain.RoleProjectManager}
	store := &taskManagementHTTPStore{authHTTPStore: authHTTPStore{user: user}}
	tokens := tokenManagerForHTTP()
	router := NewRouter(config.Config{NodeEnv: "test", FrontendURL: "http://localhost:5173"}, service.New(store, tokens))
	token, err := tokens.GenerateAccessToken(user)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/10", strings.NewReader(`{
		"title":"Edited backlog task", "project_id":9, "creator_id":7,
		"assignee_id":0, "due_date":"2026-10-20T00:00:00.000Z"
	}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH /tasks/10 status = %d, want 200: %s", rec.Code, rec.Body.String())
	}
	if store.taskPatch.Title == nil || *store.taskPatch.Title != "Edited backlog task" || store.taskPatch.DueDate == nil || *store.taskPatch.DueDate != "2026-10-20" || store.taskPatch.ProjectID != nil || store.taskPatch.AssigneeID != nil || !store.taskPatch.ClearAssignee {
		t.Fatalf("task patch = %#v", store.taskPatch)
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

type authSessionStore struct {
	service.Store
	user   domain.User
	tokens map[string]bool
}

type authTokenResponse struct {
	Data struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	} `json:"data"`
}

func (s *authSessionStore) GetUserByEmail(_ context.Context, email string) (domain.User, error) {
	if email != s.user.Email {
		return domain.User{}, domain.ErrNotFound
	}
	return s.user, nil
}

func (s *authSessionStore) GetUserByID(_ context.Context, id int) (domain.User, error) {
	if id != s.user.ID {
		return domain.User{}, domain.ErrNotFound
	}
	return s.user, nil
}

func (s *authSessionStore) SaveRefreshToken(_ context.Context, _ int, hash string, _ time.Time) error {
	s.tokens[hash] = true
	return nil
}

func (s *authSessionStore) FindValidRefreshToken(_ context.Context, _ int, hash string, _ time.Time) (bool, error) {
	return s.tokens[hash], nil
}

func (s *authSessionStore) RevokeRefreshToken(_ context.Context, hash string) error {
	s.tokens[hash] = false
	return nil
}

func (s *authSessionStore) RefreshTokenBelongsToUser(_ context.Context, _ int, hash string) (bool, error) {
	_, ok := s.tokens[hash]
	return ok, nil
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

type taskPatchHTTPStore struct {
	authHTTPStore
	patch service.TaskPatchInput
}

type taskManagementHTTPStore struct {
	authHTTPStore
	filter      service.TaskFilter
	createdTask service.TaskInput
	taskPatch   service.TaskPatchInput
}

func (s *taskManagementHTTPStore) CreateTask(_ context.Context, input service.TaskInput) (domain.Task, error) {
	s.createdTask = input
	return domain.Task{ID: 10, Title: input.Title, Status: *input.Status, ProjectID: input.ProjectID}, nil
}

func (s *taskManagementHTTPStore) UpdateTask(_ context.Context, _ domain.AuthUser, id int, input service.TaskPatchInput) (domain.Task, error) {
	s.taskPatch = input
	return domain.Task{ID: id, Title: *input.Title}, nil
}

func (s *taskManagementHTTPStore) ListTasks(_ context.Context, _ domain.AuthUser, filter service.TaskFilter) (domain.Paginated[domain.Task], error) {
	s.filter = filter
	return domain.Paginated[domain.Task]{Data: []domain.Task{{ID: 10}}, Pagination: domain.Pagination{Page: 1, Limit: 20, TotalItems: 1, TotalPages: 1}}, nil
}

func (s *taskManagementHTTPStore) GetProjectByID(_ context.Context, _ domain.AuthUser, id int) (domain.Project, error) {
	return domain.Project{ID: id}, nil
}

func (s *taskManagementHTTPStore) GetTaskByID(_ context.Context, _ domain.AuthUser, id int) (domain.Task, error) {
	return domain.Task{ID: id}, nil
}

func (s *taskManagementHTTPStore) AssignUserToTask(_ context.Context, input service.TaskAssignmentInput) (domain.TaskAssignment, bool, error) {
	return domain.TaskAssignment{ID: 1, TaskID: input.TaskID, UserID: input.UserID}, false, nil
}

func (s *taskManagementHTTPStore) CreateNotification(_ context.Context, _ service.NotificationInput) (domain.Notification, error) {
	return domain.Notification{}, nil
}

func (s *taskManagementHTTPStore) RemoveTaskAssignmentByUser(_ context.Context, taskID, userID int) (domain.TaskAssignment, error) {
	return domain.TaskAssignment{ID: 1, TaskID: taskID, UserID: userID}, nil
}

func (s *taskManagementHTTPStore) ReorderTasks(_ context.Context, _ domain.AuthUser, _ service.TaskReorderInput) ([]domain.Task, error) {
	return []domain.Task{{ID: 10}}, nil
}

func (s *taskManagementHTTPStore) ListTaskActivity(_ context.Context, _ domain.AuthUser, _ int) ([]domain.ActivityLog, error) {
	return []domain.ActivityLog{}, nil
}

func (s *taskPatchHTTPStore) GetTaskByID(context.Context, domain.AuthUser, int) (domain.Task, error) {
	return domain.Task{ID: 10}, nil
}

func (s *taskPatchHTTPStore) UpdateTask(_ context.Context, _ domain.AuthUser, _ int, input service.TaskPatchInput) (domain.Task, error) {
	s.patch = input
	return domain.Task{ID: 10, Title: *input.Title, Position: input.Position}, nil
}

func (s *dashboardHTTPStore) DashboardSummary(context.Context, domain.AuthUser) (domain.DashboardSummary, error) {
	return s.summary, nil
}
