package service

import "project-manager-go/internal/domain"

type CreateUserInput struct {
	Name     string           `json:"name"`
	Email    string           `json:"email"`
	Password string           `json:"password,omitempty"`
	Role     *domain.UserRole `json:"role,omitempty"`
}

type UpdateUserInput struct {
	Name      *string          `json:"name"`
	Email     *string          `json:"email"`
	Password  *string          `json:"password,omitempty"`
	AvatarURL *string          `json:"avatarUrl"`
	Role      *domain.UserRole `json:"role,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResult struct {
	User         domain.User `json:"user"`
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"-"`
}

type PageFilter struct {
	Page  int
	Limit int
}

type ListUsersFilter struct {
	PageFilter
	Search string
	Role   *domain.UserRole
	SortBy string
	Order  string
}

type TeamInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type TeamMemberInput struct {
	UserID int                    `json:"userId"`
	Role   *domain.TeamMemberRole `json:"role"`
}

type ProjectInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	TeamID      *int    `json:"teamId"`
	OwnerID     int     `json:"-"`
}

type ProjectFilter struct {
	PageFilter
	Search string
	TeamID *int
	SortBy string
	Order  string
}

type ProjectTeamInput struct {
	ProjectID int `json:"projectId"`
	TeamID    int `json:"teamId"`
}

type TaskInput struct {
	Title       string               `json:"title"`
	Description *string              `json:"description"`
	Status      *domain.TaskStatus   `json:"status"`
	Priority    *domain.TaskPriority `json:"priority"`
	ProjectID   int                  `json:"projectId"`
	CreatorID   int                  `json:"-"`
	AssigneeID  int                  `json:"assigneeId"`
	DueDate     *string              `json:"dueDate"`
	Position    *int                 `json:"position"`
}

type TaskPatchInput struct {
	Title       *string              `json:"title"`
	Description *string              `json:"description"`
	Status      *domain.TaskStatus   `json:"status"`
	Priority    *domain.TaskPriority `json:"priority"`
	ProjectID   *int                 `json:"projectId"`
	AssigneeID  *int                 `json:"assigneeId"`
	DueDate     *string              `json:"dueDate"`
	Position    *int                 `json:"position"`
}

type TaskFilter struct {
	PageFilter
	Search     string
	ProjectID  *int
	Status     *domain.TaskStatus
	Priority   *domain.TaskPriority
	AssigneeID *int
	SortBy     string
	Order      string
}

type TaskAssignmentInput struct {
	TaskID int `json:"taskId"`
	UserID int `json:"userId"`
}

type CommentInput struct {
	Content  string `json:"content"`
	TaskID   int    `json:"-"`
	AuthorID int    `json:"-"`
}

type AttachmentInput struct {
	FileName   string  `json:"fileName"`
	FileURL    string  `json:"fileUrl"`
	FileSize   *int    `json:"fileSize"`
	MimeType   *string `json:"mimeType"`
	TaskID     int     `json:"-"`
	UploaderID int     `json:"-"`
}

type NotificationInput struct {
	UserID  int                     `json:"userId"`
	ActorID *int                    `json:"actorId"`
	Type    domain.NotificationType `json:"type"`
	TaskID  *int                    `json:"taskId"`
}
