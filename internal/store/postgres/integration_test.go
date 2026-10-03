//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"project-manager-go/internal/auth"
	"project-manager-go/internal/config"
	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func TestStoreIntegrationRepositoryRoundTrip(t *testing.T) {
	ctx := t.Context()
	store := newIntegrationStore(t)
	suffix := uniqueSuffix()

	adminRole := domain.RoleAdmin
	managerRole := domain.RoleProjectManager
	memberRole := domain.RoleTeamMember
	admin, err := store.CreateUser(ctx, service.CreateUserInput{
		Name:     "Integration Admin " + suffix,
		Email:    "integration-admin-" + suffix + "@example.com",
		Password: "stored-password",
		Role:     &adminRole,
	})
	if err != nil {
		t.Fatalf("CreateUser(admin) error = %v", err)
	}
	manager, err := store.CreateUser(ctx, service.CreateUserInput{
		Name:     "Integration Manager " + suffix,
		Email:    "integration-manager-" + suffix + "@example.com",
		Password: "stored-password",
		Role:     &managerRole,
	})
	if err != nil {
		t.Fatalf("CreateUser(manager) error = %v", err)
	}
	member, err := store.CreateUser(ctx, service.CreateUserInput{
		Name:     "Integration Member " + suffix,
		Email:    "integration-member-" + suffix + "@example.com",
		Password: "stored-password",
		Role:     &memberRole,
	})
	if err != nil {
		t.Fatalf("CreateUser(member) error = %v", err)
	}

	adminActor := domain.AuthUser{ID: admin.ID, Email: admin.Email, Role: admin.Role}
	managerActor := domain.AuthUser{ID: manager.ID, Email: manager.Email, Role: manager.Role}
	description := "Integration repository team"
	teamName := "Integration Team " + suffix
	team, err := store.CreateTeam(ctx, service.TeamInput{
		Name:        &teamName,
		Description: &description,
	})
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}
	ownerRole := domain.TeamRoleOwner
	if _, err := store.AddTeamMember(ctx, team.ID, service.TeamMemberInput{
		UserID: manager.ID,
		Role:   &ownerRole,
	}); err != nil {
		t.Fatalf("AddTeamMember(manager) error = %v", err)
	}
	if _, err := store.AddTeamMember(ctx, team.ID, service.TeamMemberInput{UserID: member.ID}); err != nil {
		t.Fatalf("AddTeamMember(member) error = %v", err)
	}

	canCreate, err := store.CanCreateProjectForTeam(ctx, managerActor, team.ID)
	if err != nil {
		t.Fatalf("CanCreateProjectForTeam() error = %v", err)
	}
	if !canCreate {
		t.Fatal("CanCreateProjectForTeam() = false, want true for team owner")
	}

	projectName := "Integration Project " + suffix
	project, err := store.CreateProject(ctx, service.ProjectInput{
		Name:        &projectName,
		Description: &description,
		TeamID:      &team.ID,
		OwnerID:     manager.ID,
	})
	if err != nil {
		t.Fatalf("CreateProject() error = %v", err)
	}
	if _, _, err := store.AssignTeamToProject(ctx, adminActor, service.ProjectTeamInput{
		ProjectID: project.ID,
		TeamID:    team.ID,
	}); err != nil {
		t.Fatalf("AssignTeamToProject() error = %v", err)
	}

	status := domain.TaskStatusTodo
	priority := domain.TaskPriorityHigh
	position := 1
	dueDate := "2026-09-05"
	task, err := store.CreateTask(ctx, service.TaskInput{
		Title:       "Integration Task " + suffix,
		Description: &description,
		Status:      &status,
		Priority:    &priority,
		ProjectID:   project.ID,
		CreatorID:   manager.ID,
		AssigneeID:  member.ID,
		DueDate:     &dueDate,
		Position:    &position,
	})
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}
	if task.Status != domain.TaskStatusTodo || task.ProjectID != project.ID {
		t.Fatalf("task = %#v, want todo task in project %d", task, project.ID)
	}

	if _, exists, err := store.AssignUserToTask(ctx, service.TaskAssignmentInput{
		TaskID: task.ID,
		UserID: member.ID,
	}); err != nil || exists {
		t.Fatalf("AssignUserToTask() exists=%t error=%v, want new assignment", exists, err)
	}
	if _, exists, err := store.AssignUserToTask(ctx, service.TaskAssignmentInput{
		TaskID: task.ID,
		UserID: member.ID,
	}); err != nil || !exists {
		t.Fatalf("AssignUserToTask(duplicate) exists=%t error=%v, want existing assignment", exists, err)
	}

	comment, err := store.CreateComment(ctx, service.CommentInput{
		Content:  "Repository comment " + suffix,
		TaskID:   task.ID,
		AuthorID: manager.ID,
	})
	if err != nil {
		t.Fatalf("CreateComment() error = %v", err)
	}
	fileSize := 1234
	mimeType := "application/pdf"
	attachment, err := store.CreateAttachment(ctx, service.AttachmentInput{
		FileName:     "repo-" + suffix + ".pdf",
		OriginalName: "repo.pdf",
		StorageKey:   "attachments/repo-" + suffix + ".pdf",
		FileSize:     &fileSize,
		MimeType:     &mimeType,
		TaskID:       task.ID,
		UploaderID:   manager.ID,
	})
	if err != nil {
		t.Fatalf("CreateAttachment() error = %v", err)
	}
	notification, err := store.CreateNotification(ctx, service.NotificationInput{
		UserID: member.ID,
		Type:   domain.NotificationTaskAssigned,
		TaskID: &task.ID,
	})
	if err != nil {
		t.Fatalf("CreateNotification() error = %v", err)
	}

	gotTask, err := store.GetTaskByID(ctx, memberActor(member), task.ID)
	if err != nil {
		t.Fatalf("GetTaskByID() error = %v", err)
	}
	if gotTask.ID != task.ID || len(gotTask.Comments) != 1 || len(gotTask.Attachments) != 1 {
		t.Fatalf("GetTaskByID() = %#v, want task with one comment and one attachment", gotTask)
	}
	if gotTask.Comments[0].ID != comment.ID || gotTask.Attachments[0].ID != attachment.ID {
		t.Fatalf("GetTaskByID() nested resources mismatch: %#v", gotTask)
	}

	notifications, err := store.ListNotifications(ctx, member.ID)
	if err != nil {
		t.Fatalf("ListNotifications() error = %v", err)
	}
	if len(notifications) != 1 || notifications[0].ID != notification.ID {
		t.Fatalf("notifications = %#v, want created notification", notifications)
	}
}

func TestStoreIntegrationDashboardSummary(t *testing.T) {
	ctx := t.Context()
	store := newIntegrationStore(t)
	suffix := uniqueSuffix()

	managerRole := domain.RoleProjectManager
	memberRole := domain.RoleTeamMember
	manager, err := store.CreateUser(ctx, service.CreateUserInput{
		Name:     "Dashboard Manager " + suffix,
		Email:    "dashboard-manager-" + suffix + "@example.com",
		Password: "stored-password",
		Role:     &managerRole,
	})
	if err != nil {
		t.Fatalf("CreateUser(manager) error = %v", err)
	}
	member, err := store.CreateUser(ctx, service.CreateUserInput{
		Name:     "Dashboard Member " + suffix,
		Email:    "dashboard-member-" + suffix + "@example.com",
		Password: "stored-password",
		Role:     &memberRole,
	})
	if err != nil {
		t.Fatalf("CreateUser(member) error = %v", err)
	}

	teamName := "Dashboard Team " + suffix
	team, err := store.CreateTeam(ctx, service.TeamInput{Name: &teamName})
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}
	if _, err := store.AddTeamMember(ctx, team.ID, service.TeamMemberInput{UserID: member.ID}); err != nil {
		t.Fatalf("AddTeamMember(member) error = %v", err)
	}

	projectName := "Dashboard Project " + suffix
	project, err := store.CreateProject(ctx, service.ProjectInput{
		Name:    &projectName,
		TeamID:  &team.ID,
		OwnerID: manager.ID,
	})
	if err != nil {
		t.Fatalf("CreateProject() error = %v", err)
	}
	if _, err := store.db.ExecContext(ctx, "UPDATE projects SET status = $1 WHERE id = $2", domain.ProjectStatusActive, project.ID); err != nil {
		t.Fatalf("setting project active: %v", err)
	}
	planningProjectName := "Dashboard Planning Project " + suffix
	if _, err := store.CreateProject(ctx, service.ProjectInput{
		Name:    &planningProjectName,
		TeamID:  &team.ID,
		OwnerID: manager.ID,
	}); err != nil {
		t.Fatalf("CreateProject(planning) error = %v", err)
	}

	todo := domain.TaskStatusTodo
	inProgress := domain.TaskStatusInProgress
	review := domain.TaskStatusReview
	done := domain.TaskStatusDone
	high := domain.TaskPriorityHigh
	urgent := domain.TaskPriorityUrgent
	medium := domain.TaskPriorityMedium

	taskInputs := []service.TaskInput{
		{Title: "Dashboard Todo " + suffix, Status: &todo, Priority: &high, ProjectID: project.ID, CreatorID: manager.ID, AssigneeID: member.ID},
		{Title: "Dashboard Doing " + suffix, Status: &inProgress, Priority: &medium, ProjectID: project.ID, CreatorID: manager.ID, AssigneeID: member.ID},
		{Title: "Dashboard Review " + suffix, Status: &review, Priority: &urgent, ProjectID: project.ID, CreatorID: manager.ID, AssigneeID: member.ID},
		{Title: "Dashboard Done " + suffix, Status: &done, Priority: &high, ProjectID: project.ID, CreatorID: manager.ID, AssigneeID: member.ID},
	}
	tasks := make([]domain.Task, 0, len(taskInputs))
	for _, input := range taskInputs {
		task, err := store.CreateTask(ctx, input)
		if err != nil {
			t.Fatalf("CreateTask(%s) error = %v", input.Title, err)
		}
		tasks = append(tasks, task)
	}

	dateUpdates := []struct {
		taskID int
		expr   string
	}{
		{taskID: tasks[0].ID, expr: "CURRENT_DATE"},
		{taskID: tasks[1].ID, expr: "CURRENT_DATE + INTERVAL '1 day'"},
		{taskID: tasks[2].ID, expr: "CURRENT_DATE - INTERVAL '1 day'"},
		{taskID: tasks[3].ID, expr: "CURRENT_DATE - INTERVAL '1 day'"},
	}
	for _, update := range dateUpdates {
		query := fmt.Sprintf("UPDATE tasks SET due_date = %s WHERE id = $1", update.expr)
		if _, err := store.db.ExecContext(ctx, query, update.taskID); err != nil {
			t.Fatalf("setting due date for task %d: %v", update.taskID, err)
		}
	}

	summary, err := store.DashboardSummary(ctx, memberActor(member))
	if err != nil {
		t.Fatalf("DashboardSummary() error = %v", err)
	}

	if summary.TotalActiveProjects != 1 {
		t.Fatalf("TotalActiveProjects = %d, want 1", summary.TotalActiveProjects)
	}
	if summary.TaskCountPerStatus[domain.TaskStatusTodo] != 1 ||
		summary.TaskCountPerStatus[domain.TaskStatusInProgress] != 1 ||
		summary.TaskCountPerStatus[domain.TaskStatusReview] != 1 ||
		summary.TaskCountPerStatus[domain.TaskStatusDone] != 1 {
		t.Fatalf("TaskCountPerStatus = %#v, want one task per created status", summary.TaskCountPerStatus)
	}
	if summary.ActiveProgress.Doing != 1 || summary.ActiveProgress.Todo != 1 ||
		summary.ActiveProgress.Total != 2 || summary.ActiveProgress.Ratio != 0.5 ||
		summary.ActiveProgress.Percentage != 50 {
		t.Fatalf("ActiveProgress = %#v, want 1/2 progress", summary.ActiveProgress)
	}
	if summary.InReview != 1 || summary.DueSoon != 2 || summary.Overdue != 1 || summary.OverdueTaskCount != 1 {
		t.Fatalf("dashboard counts = review:%d dueSoon:%d overdue:%d alias:%d, want 1/2/1/1",
			summary.InReview, summary.DueSoon, summary.Overdue, summary.OverdueTaskCount)
	}
	if summary.WorkloadPerMember[fmt.Sprintf("%d", member.ID)] != 3 {
		t.Fatalf("WorkloadPerMember = %#v, want three open tasks for member", summary.WorkloadPerMember)
	}
	if len(summary.RecentTasks) == 0 || len(summary.UpcomingDeadlines) != 2 || len(summary.HighPriorityTasks) != 2 {
		t.Fatalf("dashboard task lists lengths = recent:%d upcoming:%d high:%d, want nonzero/2/2",
			len(summary.RecentTasks), len(summary.UpcomingDeadlines), len(summary.HighPriorityTasks))
	}
	if len(summary.LatestUpdates) != 0 {
		t.Fatalf("LatestUpdates = %#v, want empty until activity service is implemented", summary.LatestUpdates)
	}
}

func TestStoreIntegrationServiceWorkflow(t *testing.T) {
	ctx := t.Context()
	store := newIntegrationStore(t)
	svc := service.New(store, integrationTokenManager())
	suffix := uniqueSuffix()

	adminRole := domain.RoleAdmin
	admin, err := svc.Register(ctx, service.CreateUserInput{
		Name:     "WorkflowAdmin" + suffix,
		Email:    "workflow-admin-" + suffix + "@example.com",
		Password: "password123",
		Role:     &adminRole,
	})
	if err != nil {
		t.Fatalf("Register(admin) error = %v", err)
	}
	login, err := svc.Login(ctx, service.LoginInput{
		Email:    admin.Email,
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if login.AccessToken == "" || login.RefreshToken == "" {
		t.Fatalf("Login() tokens are empty: %#v", login)
	}
	authUser, err := svc.Authenticate(login.AccessToken)
	if err != nil {
		t.Fatalf("Authenticate() error = %v", err)
	}
	if authUser.ID != admin.ID || authUser.Role != domain.RoleAdmin {
		t.Fatalf("Authenticate() = %#v, want admin user %d", authUser, admin.ID)
	}
	refreshed, err := svc.RefreshAccessToken(ctx, login.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshAccessToken() error = %v", err)
	}
	if refreshed.AccessToken == "" || refreshed.RefreshToken == "" {
		t.Fatalf("RefreshAccessToken() returned empty tokens: %#v", refreshed)
	}

	managerRole := domain.RoleProjectManager
	memberRole := domain.RoleTeamMember
	manager, err := svc.CreateUser(ctx, service.CreateUserInput{
		Name:     "WorkflowManager" + suffix,
		Email:    "workflow-manager-" + suffix + "@example.com",
		Password: "password123",
		Role:     &managerRole,
	})
	if err != nil {
		t.Fatalf("CreateUser(manager) error = %v", err)
	}
	member, err := svc.CreateUser(ctx, service.CreateUserInput{
		Name:     "MentionTarget" + suffix,
		Email:    "workflow-member-" + suffix + "@example.com",
		Password: "password123",
		Role:     &memberRole,
	})
	if err != nil {
		t.Fatalf("CreateUser(member) error = %v", err)
	}

	teamName := "Workflow Team " + suffix
	description := "Integration service workflow"
	team, err := svc.CreateTeam(ctx, authUser, service.TeamInput{
		Name:        &teamName,
		Description: &description,
	})
	if err != nil {
		t.Fatalf("CreateTeam() error = %v", err)
	}
	ownerRole := domain.TeamRoleOwner
	if _, err := svc.AddTeamMember(ctx, authUser, team.ID, service.TeamMemberInput{
		UserID: manager.ID,
		Role:   &ownerRole,
	}); err != nil {
		t.Fatalf("AddTeamMember(manager) error = %v", err)
	}
	if _, err := svc.AddTeamMember(ctx, authUser, team.ID, service.TeamMemberInput{UserID: member.ID}); err != nil {
		t.Fatalf("AddTeamMember(member) error = %v", err)
	}

	projectName := "Workflow Project " + suffix
	project, err := svc.CreateProject(ctx, authUser, service.ProjectInput{
		Name:        &projectName,
		Description: &description,
		TeamID:      &team.ID,
	})
	if err != nil {
		t.Fatalf("CreateProject() error = %v", err)
	}
	projects, err := svc.ListProjects(ctx, memberActor(member), service.ProjectFilter{PageFilter: service.PageFilter{Page: 1, Limit: 10}})
	if err != nil {
		t.Fatalf("ListProjects(member) error = %v", err)
	}
	if projects.Pagination.TotalItems != 1 || projects.Data[0].ID != project.ID {
		t.Fatalf("ListProjects(member) = %#v, want created project", projects)
	}

	task, err := svc.CreateTask(ctx, authUser, service.TaskInput{
		Title:      "Workflow Task " + suffix,
		ProjectID:  project.ID,
		AssigneeID: member.ID,
	})
	if err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}
	assignment, exists, err := svc.AssignUserToTask(ctx, authUser, service.TaskAssignmentInput{
		TaskID: task.ID,
		UserID: member.ID,
	})
	if err != nil || exists {
		t.Fatalf("AssignUserToTask() assignment=%#v exists=%t error=%v, want new assignment", assignment, exists, err)
	}
	comment, err := svc.CreateComment(ctx, authUser, service.CommentInput{
		Content: "Please review @" + member.Name,
		TaskID:  task.ID,
	})
	if err != nil {
		t.Fatalf("CreateComment() error = %v", err)
	}
	fileSize := 2048
	mimeType := "image/png"
	attachment, err := svc.CreateAttachment(ctx, authUser, task.ID, service.AttachmentInput{
		FileName:     "workflow-" + suffix + ".png",
		OriginalName: "workflow.png",
		StorageKey:   "attachments/workflow-" + suffix + ".png",
		FileSize:     &fileSize,
		MimeType:     &mimeType,
	})
	if err != nil {
		t.Fatalf("CreateAttachment() error = %v", err)
	}

	gotTask, err := svc.GetTaskByID(ctx, memberActor(member), task.ID)
	if err != nil {
		t.Fatalf("GetTaskByID(member) error = %v", err)
	}
	if gotTask.ID != task.ID || len(gotTask.Comments) != 1 || len(gotTask.Attachments) != 1 {
		t.Fatalf("GetTaskByID(member) = %#v, want task with one comment and one attachment", gotTask)
	}
	if gotTask.Comments[0].ID != comment.ID || gotTask.Attachments[0].ID != attachment.ID {
		t.Fatalf("GetTaskByID(member) nested resources mismatch: %#v", gotTask)
	}

	notifications, err := svc.ListNotifications(ctx, member.ID)
	if err != nil {
		t.Fatalf("ListNotifications() error = %v", err)
	}
	if len(notifications) != 2 {
		t.Fatalf("notifications = %#v, want assignment and mention notifications", notifications)
	}
	read, err := svc.MarkAllNotificationsRead(ctx, memberActor(member))
	if err != nil {
		t.Fatalf("MarkAllNotificationsRead() error = %v", err)
	}
	if len(read) != 2 {
		t.Fatalf("MarkAllNotificationsRead() count = %d, want 2", len(read))
	}
	for _, notification := range read {
		if !notification.IsRead {
			t.Fatalf("notification was not marked read: %#v", notification)
		}
	}

	if err := svc.Logout(ctx, authUser, refreshed.RefreshToken); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
}

func newIntegrationStore(t *testing.T) *Store {
	t.Helper()

	db := integrationDB(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		t.Fatalf("BeginTxx() error = %v", err)
	}
	t.Cleanup(func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			t.Logf("rolling back integration transaction: %v", err)
		}
	})

	return &Store{db: tx}
}

func integrationDB(t *testing.T) *sqlx.DB {
	t.Helper()

	root := integrationProjectRoot(t)
	loadRootConfig(t, root)
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		t.Skip("set TEST_DATABASE_URL or DATABASE_URL to run integration tests")
	}

	runIntegrationMigrations(t, databaseURL)

	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	t.Cleanup(cancel)
	db, err := sqlx.ConnectContext(ctx, "pgx", databaseURL)
	if err != nil {
		t.Fatalf("connecting integration database: %v", err)
	}
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Minute)
	db.SetConnMaxIdleTime(time.Minute)
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Logf("closing integration database: %v", err)
		}
	})
	return db
}

func runIntegrationMigrations(t *testing.T, databaseURL string) {
	t.Helper()

	migrationsDir := filepath.Join(integrationProjectRoot(t), "migrations")
	migrator, err := migrate.New("file://"+migrationsDir, databaseURL)
	if err != nil {
		t.Fatalf("creating migrator: %v", err)
	}
	defer func() {
		sourceErr, databaseErr := migrator.Close()
		if sourceErr != nil {
			t.Logf("closing migration source: %v", sourceErr)
		}
		if databaseErr != nil {
			t.Logf("closing migration database: %v", databaseErr)
		}
	}()
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("running migrations: %v", err)
	}
}

func integrationProjectRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting working directory: %v", err)
	}
	root, err := filepath.Abs(filepath.Join(wd, "../../.."))
	if err != nil {
		t.Fatalf("resolving project root: %v", err)
	}
	return root
}

func loadRootConfig(t *testing.T, root string) {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getting working directory: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("changing to project root: %v", err)
	}
	_, loadErr := config.Load()
	if err := os.Chdir(cwd); err != nil {
		t.Fatalf("restoring working directory: %v", err)
	}
	if loadErr != nil {
		t.Fatalf("loading config: %v", loadErr)
	}
}

func integrationTokenManager() auth.TokenManager {
	now := time.Date(2026, 9, 4, 9, 0, 0, 0, time.UTC)
	return auth.TokenManager{
		AccessSecret:  []byte("integration-access-secret"),
		RefreshSecret: []byte("integration-refresh-secret"),
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    24 * time.Hour,
		Now: func() time.Time {
			current := now
			now = now.Add(time.Second)
			return current
		},
	}
}

func memberActor(user domain.User) domain.AuthUser {
	return domain.AuthUser{ID: user.ID, Email: user.Email, Role: user.Role}
}

func uniqueSuffix() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
