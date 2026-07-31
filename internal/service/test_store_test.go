package service

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"project-manager-go/internal/auth"
	"project-manager-go/internal/domain"
)

type memoryStore struct {
	users            map[int]domain.User
	usersByEmail     map[string]int
	comments         map[int]domain.Comment
	attachments      map[int]domain.Attachment
	notifications    []NotificationInput
	canCreateProject bool
}

func newMemoryStore(t *testing.T) *memoryStore {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte("secret123"), 10)
	if err != nil {
		t.Fatal(err)
	}
	return &memoryStore{
		users: map[int]domain.User{
			1: {
				ID:           1,
				Name:         "Author",
				Email:        "author@example.com",
				PasswordHash: string(hash),
				Role:         domain.RoleTeamMember,
			},
			2: {
				ID:           2,
				Name:         "Dian",
				Email:        "dian@example.com",
				PasswordHash: string(hash),
				Role:         domain.RoleTeamMember,
			},
		},
		usersByEmail: map[string]int{
			"author@example.com": 1,
			"dian@example.com":   2,
		},
		comments: map[int]domain.Comment{
			1: {ID: 1, TaskID: 10, AuthorID: 1, Content: "original"},
		},
		attachments: map[int]domain.Attachment{
			1: {
				ID:           1,
				TaskID:       10,
				UploaderID:   1,
				FileName:     "spec.pdf",
				OriginalName: "spec.pdf",
				StorageKey:   "attachments/tasks/10/spec.pdf",
			},
		},
		notifications:    []NotificationInput{},
		canCreateProject: true,
	}
}

func testTokenManager() auth.TokenManager {
	now := time.Date(2026, 7, 13, 8, 0, 0, 0, time.UTC)
	return auth.TokenManager{
		AccessSecret:  []byte("access"),
		RefreshSecret: []byte("refresh"),
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
		Now:           func() time.Time { return now },
	}
}

func (m *memoryStore) CreateUser(context.Context, CreateUserInput) (domain.User, error) {
	return domain.User{}, nil
}

func (m *memoryStore) GetUserByEmail(_ context.Context, email string) (domain.User, error) {
	id, ok := m.usersByEmail[email]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return m.users[id], nil
}

func (m *memoryStore) GetUserByID(_ context.Context, id int) (domain.User, error) {
	user, ok := m.users[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

func (m *memoryStore) ListUsers(context.Context, ListUsersFilter) (domain.Paginated[domain.User], error) {
	return domain.Paginated[domain.User]{}, nil
}

func (m *memoryStore) UpdateUser(context.Context, int, UpdateUserInput) (domain.User, error) {
	return domain.User{}, nil
}

func (m *memoryStore) DeleteUser(context.Context, int) (domain.User, error) {
	return domain.User{}, nil
}

func (m *memoryStore) ListUserTasks(context.Context, int, PageFilter) (domain.Paginated[domain.Task], error) {
	return domain.Paginated[domain.Task]{}, nil
}

func (m *memoryStore) SaveRefreshToken(context.Context, int, string, time.Time) error {
	return nil
}

func (m *memoryStore) FindValidRefreshToken(context.Context, int, string, time.Time) (bool, error) {
	return true, nil
}

func (m *memoryStore) RevokeRefreshToken(context.Context, string) error {
	return nil
}

func (m *memoryStore) RefreshTokenBelongsToUser(context.Context, int, string) (bool, error) {
	return true, nil
}

func (m *memoryStore) CreateTeam(context.Context, TeamInput) (domain.Team, error) {
	return domain.Team{}, nil
}

func (m *memoryStore) ListTeams(context.Context, PageFilter) (domain.Paginated[domain.Team], error) {
	return domain.Paginated[domain.Team]{}, nil
}

func (m *memoryStore) GetTeamByID(context.Context, int) (domain.Team, error) {
	return domain.Team{}, nil
}

func (m *memoryStore) UpdateTeam(context.Context, int, TeamInput) (domain.Team, error) {
	return domain.Team{}, nil
}

func (m *memoryStore) DeleteTeam(context.Context, int) (domain.Team, error) {
	return domain.Team{}, nil
}

func (m *memoryStore) ListTeamMembers(context.Context, int) ([]domain.TeamMember, error) {
	return nil, nil
}

func (m *memoryStore) AddTeamMember(context.Context, int, TeamMemberInput) (domain.TeamMember, error) {
	return domain.TeamMember{}, nil
}

func (m *memoryStore) RemoveTeamMember(context.Context, int, int) (domain.TeamMember, error) {
	return domain.TeamMember{}, nil
}

func (m *memoryStore) CanCreateProjectForTeam(context.Context, domain.AuthUser, int) (bool, error) {
	return m.canCreateProject, nil
}

func (m *memoryStore) CreateProject(context.Context, ProjectInput) (domain.Project, error) {
	return domain.Project{}, nil
}

func (m *memoryStore) ListProjects(context.Context, domain.AuthUser, ProjectFilter) (domain.Paginated[domain.Project], error) {
	return domain.Paginated[domain.Project]{}, nil
}

func (m *memoryStore) GetProjectByID(context.Context, domain.AuthUser, int) (domain.Project, error) {
	return domain.Project{}, nil
}

func (m *memoryStore) UpdateProject(context.Context, int, ProjectInput) (domain.Project, error) {
	return domain.Project{}, nil
}

func (m *memoryStore) DeleteProject(context.Context, int) (domain.Project, error) {
	return domain.Project{}, nil
}

func (m *memoryStore) ListProjectTasks(
	context.Context,
	domain.AuthUser,
	int,
	PageFilter,
) (domain.Paginated[domain.Task], error) {
	return domain.Paginated[domain.Task]{}, nil
}

func (m *memoryStore) AssignTeamToProject(context.Context, ProjectTeamInput) (domain.ProjectTeam, bool, error) {
	return domain.ProjectTeam{}, false, nil
}

func (m *memoryStore) ListProjectTeams(context.Context, int) ([]domain.ProjectTeam, error) {
	return nil, nil
}

func (m *memoryStore) GetProjectTeamByID(context.Context, int) (domain.ProjectTeam, error) {
	return domain.ProjectTeam{}, nil
}

func (m *memoryStore) RemoveTeamFromProject(context.Context, int) (domain.ProjectTeam, error) {
	return domain.ProjectTeam{}, nil
}

func (m *memoryStore) CreateTask(_ context.Context, input TaskInput) (domain.Task, error) {
	return domain.Task{
		ID:         1,
		Title:      input.Title,
		Status:     *input.Status,
		Priority:   input.Priority,
		ProjectID:  input.ProjectID,
		CreatorID:  input.CreatorID,
		AssigneeID: input.AssigneeID,
	}, nil
}

func (m *memoryStore) ListTasks(context.Context, domain.AuthUser, TaskFilter) (domain.Paginated[domain.Task], error) {
	return domain.Paginated[domain.Task]{}, nil
}

func (m *memoryStore) GetTaskByID(context.Context, domain.AuthUser, int) (domain.Task, error) {
	return domain.Task{}, nil
}

func (m *memoryStore) UpdateTask(context.Context, int, TaskPatchInput) (domain.Task, error) {
	return domain.Task{}, nil
}

func (m *memoryStore) DeleteTask(context.Context, int) (domain.Task, error) {
	return domain.Task{}, nil
}

func (m *memoryStore) AssignUserToTask(context.Context, TaskAssignmentInput) (domain.TaskAssignment, bool, error) {
	return domain.TaskAssignment{ID: 1}, false, nil
}

func (m *memoryStore) ListTaskAssignments(context.Context, int) ([]domain.TaskAssignment, error) {
	return nil, nil
}

func (m *memoryStore) GetTaskAssignmentByID(context.Context, int) (domain.TaskAssignment, error) {
	return domain.TaskAssignment{}, nil
}

func (m *memoryStore) RemoveTaskAssignment(context.Context, int) (domain.TaskAssignment, error) {
	return domain.TaskAssignment{}, nil
}

func (m *memoryStore) CreateComment(_ context.Context, input CommentInput) (domain.Comment, error) {
	comment := domain.Comment{ID: 2, TaskID: input.TaskID, AuthorID: input.AuthorID, Content: input.Content}
	m.comments[comment.ID] = comment
	return comment, nil
}

func (m *memoryStore) ListComments(context.Context, int) ([]domain.CommentView, error) {
	return nil, nil
}

func (m *memoryStore) GetCommentByID(_ context.Context, id int) (domain.Comment, error) {
	comment, ok := m.comments[id]
	if !ok {
		return domain.Comment{}, domain.ErrNotFound
	}
	return comment, nil
}

func (m *memoryStore) UpdateComment(context.Context, int, string) (domain.Comment, error) {
	return domain.Comment{}, nil
}

func (m *memoryStore) DeleteComment(context.Context, int) (domain.Comment, error) {
	return domain.Comment{}, nil
}

func (m *memoryStore) FindUserIDByExactName(_ context.Context, name string) (int, bool, error) {
	for _, user := range m.users {
		if user.Name == name {
			return user.ID, true, nil
		}
	}
	return 0, false, nil
}

func (m *memoryStore) CreateAttachment(context.Context, AttachmentInput) (domain.Attachment, error) {
	return domain.Attachment{}, nil
}

func (m *memoryStore) ListAttachments(context.Context, int) ([]domain.Attachment, error) {
	return nil, nil
}

func (m *memoryStore) GetAttachmentByID(_ context.Context, id int) (domain.Attachment, error) {
	attachment, ok := m.attachments[id]
	if !ok {
		return domain.Attachment{}, domain.ErrNotFound
	}
	return attachment, nil
}

func (m *memoryStore) DeleteAttachment(context.Context, int) (domain.Attachment, error) {
	return domain.Attachment{}, nil
}

func (m *memoryStore) CreateNotification(_ context.Context, input NotificationInput) (domain.Notification, error) {
	m.notifications = append(m.notifications, input)
	return domain.Notification{UserID: input.UserID, Type: input.Type}, nil
}

func (m *memoryStore) ListNotifications(context.Context, int) ([]domain.Notification, error) {
	return nil, nil
}

func (m *memoryStore) MarkNotificationRead(context.Context, int, int) (domain.Notification, error) {
	return domain.Notification{}, nil
}

func (m *memoryStore) MarkAllNotificationsRead(context.Context, int) ([]domain.Notification, error) {
	return nil, nil
}

func stringPtr(value string) *string {
	return &value
}
