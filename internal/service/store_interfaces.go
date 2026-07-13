package service

import (
	"context"
	"time"

	"project-manager-go/internal/domain"
)

type Store interface {
	UserStore
	RefreshTokenStore
	TeamStore
	ProjectStore
	ProjectTeamStore
	TaskStore
	TaskAssignmentStore
	CommentStore
	AttachmentStore
	NotificationStore
}

type UserStore interface {
	CreateUser(ctx context.Context, input CreateUserInput) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id int) (domain.User, error)
	ListUsers(ctx context.Context, filter ListUsersFilter) (domain.Paginated[domain.User], error)
	UpdateUser(ctx context.Context, id int, input UpdateUserInput) (domain.User, error)
	DeleteUser(ctx context.Context, id int) (domain.User, error)
	ListUserTasks(ctx context.Context, userID int, filter PageFilter) (domain.Paginated[domain.Task], error)
}

type RefreshTokenStore interface {
	SaveRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
	FindValidRefreshToken(ctx context.Context, userID int, tokenHash string, now time.Time) (bool, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RefreshTokenBelongsToUser(ctx context.Context, userID int, tokenHash string) (bool, error)
}

type TeamStore interface {
	CreateTeam(ctx context.Context, input TeamInput) (domain.Team, error)
	ListTeams(ctx context.Context, filter PageFilter) (domain.Paginated[domain.Team], error)
	GetTeamByID(ctx context.Context, id int) (domain.Team, error)
	UpdateTeam(ctx context.Context, id int, input TeamInput) (domain.Team, error)
	DeleteTeam(ctx context.Context, id int) (domain.Team, error)
	ListTeamMembers(ctx context.Context, teamID int) ([]domain.TeamMember, error)
	AddTeamMember(ctx context.Context, teamID int, input TeamMemberInput) (domain.TeamMember, error)
	RemoveTeamMember(ctx context.Context, teamID int, userID int) (domain.TeamMember, error)
}

type ProjectStore interface {
	CanCreateProjectForTeam(ctx context.Context, user domain.AuthUser, teamID int) (bool, error)
	CreateProject(ctx context.Context, input ProjectInput) (domain.Project, error)
	ListProjects(ctx context.Context, user domain.AuthUser, filter ProjectFilter) (domain.Paginated[domain.Project], error)
	GetProjectByID(ctx context.Context, user domain.AuthUser, id int) (domain.Project, error)
	UpdateProject(ctx context.Context, id int, input ProjectInput) (domain.Project, error)
	DeleteProject(ctx context.Context, id int) (domain.Project, error)
	ListProjectTasks(
		ctx context.Context,
		user domain.AuthUser,
		projectID int,
		filter PageFilter,
	) (domain.Paginated[domain.Task], error)
}

type ProjectTeamStore interface {
	AssignTeamToProject(ctx context.Context, input ProjectTeamInput) (domain.ProjectTeam, bool, error)
	ListProjectTeams(ctx context.Context, projectID int) ([]domain.ProjectTeam, error)
	GetProjectTeamByID(ctx context.Context, id int) (domain.ProjectTeam, error)
	RemoveTeamFromProject(ctx context.Context, id int) (domain.ProjectTeam, error)
}

type TaskStore interface {
	CreateTask(ctx context.Context, input TaskInput) (domain.Task, error)
	ListTasks(ctx context.Context, user domain.AuthUser, filter TaskFilter) (domain.Paginated[domain.Task], error)
	GetTaskByID(ctx context.Context, user domain.AuthUser, id int) (domain.Task, error)
	UpdateTask(ctx context.Context, id int, input TaskPatchInput) (domain.Task, error)
	DeleteTask(ctx context.Context, id int) (domain.Task, error)
}

type TaskAssignmentStore interface {
	AssignUserToTask(ctx context.Context, input TaskAssignmentInput) (domain.TaskAssignment, bool, error)
	ListTaskAssignments(ctx context.Context, taskID int) ([]domain.TaskAssignment, error)
	GetTaskAssignmentByID(ctx context.Context, id int) (domain.TaskAssignment, error)
	RemoveTaskAssignment(ctx context.Context, id int) (domain.TaskAssignment, error)
}

type CommentStore interface {
	CreateComment(ctx context.Context, input CommentInput) (domain.Comment, error)
	ListComments(ctx context.Context, taskID int) ([]domain.CommentView, error)
	GetCommentByID(ctx context.Context, id int) (domain.Comment, error)
	UpdateComment(ctx context.Context, id int, content string) (domain.Comment, error)
	DeleteComment(ctx context.Context, id int) (domain.Comment, error)
	FindUserIDByExactName(ctx context.Context, name string) (int, bool, error)
}

type AttachmentStore interface {
	CreateAttachment(ctx context.Context, input AttachmentInput) (domain.Attachment, error)
	ListAttachments(ctx context.Context, taskID int) ([]domain.Attachment, error)
	GetAttachmentByID(ctx context.Context, id int) (domain.Attachment, error)
	DeleteAttachment(ctx context.Context, id int) (domain.Attachment, error)
}

type NotificationStore interface {
	CreateNotification(ctx context.Context, input NotificationInput) (domain.Notification, error)
	ListNotifications(ctx context.Context, userID int) ([]domain.Notification, error)
	MarkNotificationRead(ctx context.Context, id int, userID int) (domain.Notification, error)
	MarkAllNotificationsRead(ctx context.Context, userID int) ([]domain.Notification, error)
}
