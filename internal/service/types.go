package service

import "project-manager-go/internal/domain"

type CreateUserInput struct {
	Name     string           `json:"name" binding:"required,min=2"`
	Email    string           `json:"email" binding:"required,email"`
	Password string           `json:"password,omitempty" binding:"required,min=6"`
	Role     *domain.UserRole `json:"role,omitempty" binding:"omitempty,oneof=admin productOwner projectManager teamMember"`
}

type UpdateUserInput struct {
	Name             *string          `json:"name" binding:"omitempty,min=2"`
	Email            *string          `json:"email" binding:"omitempty,email"`
	Password         *string          `json:"password,omitempty" binding:"omitempty,min=6"`
	AvatarStorageKey *string          `json:"avatarStorageKey"`
	LegacyAvatarURL  *string          `json:"avatarUrl"`
	Role             *domain.UserRole `json:"role,omitempty" binding:"omitempty,oneof=admin productOwner projectManager teamMember"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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
	Name        *string `json:"name" binding:"omitempty,min=2"`
	Description *string `json:"description"`
}

type TeamMemberInput struct {
	UserID int                    `json:"userId" binding:"required,gt=0"`
	Role   *domain.TeamMemberRole `json:"role" binding:"omitempty,oneof=owner admin member"`
}

type ProjectInput struct {
	Name        *string `json:"name" binding:"omitempty,min=2"`
	Description *string `json:"description"`
	TeamID      *int    `json:"teamId" binding:"omitempty,gt=0"`
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
	ProjectID int `json:"projectId" binding:"required,gt=0"`
	TeamID    int `json:"teamId" binding:"required,gt=0"`
}

type TaskInput struct {
	Title       string               `json:"title" binding:"required,min=2"`
	Description *string              `json:"description"`
	Status      *domain.TaskStatus   `json:"status" binding:"omitempty,oneof=backlog todo in_progress review done"`
	Priority    *domain.TaskPriority `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	ProjectID   int                  `json:"projectId" binding:"required,gt=0"`
	CreatorID   int                  `json:"-"`
	AssigneeID  int                  `json:"assigneeId" binding:"required,gt=0"`
	DueDate     *string              `json:"dueDate" binding:"omitempty,datetime=2006-01-02"`
	Position    *int                 `json:"position" binding:"omitempty,gte=0"`
}

type TaskPatchInput struct {
	Title       *string              `json:"title" binding:"omitempty,min=2"`
	Description *string              `json:"description"`
	Status      *domain.TaskStatus   `json:"status" binding:"omitempty,oneof=backlog todo in_progress review done"`
	Priority    *domain.TaskPriority `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	ProjectID   *int                 `json:"projectId" binding:"omitempty,gt=0"`
	AssigneeID  *int                 `json:"assigneeId" binding:"omitempty,gt=0"`
	DueDate     *string              `json:"dueDate" binding:"omitempty,datetime=2006-01-02"`
	Position    *int                 `json:"position" binding:"omitempty,gte=0"`
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
	TaskID int `json:"taskId" binding:"required,gt=0"`
	UserID int `json:"userId" binding:"required,gt=0"`
}

type CommentInput struct {
	Content  string `json:"content" binding:"required"`
	TaskID   int    `json:"-"`
	AuthorID int    `json:"-"`
}

type AttachmentInput struct {
	FileName     string  `json:"fileName" binding:"required"`
	OriginalName string  `json:"originalName" binding:"required"`
	StorageKey   string  `json:"storageKey" binding:"required"`
	FileURL      string  `json:"fileUrl" binding:"required"`
	FileSize     *int    `json:"fileSize" binding:"omitempty,gte=0"`
	MimeType     *string `json:"mimeType"`
	TaskID       int     `json:"-"`
	UploaderID   int     `json:"-"`
}

type NotificationInput struct {
	UserID  int                     `json:"userId" binding:"required,gt=0"`
	ActorID *int                    `json:"actorId" binding:"omitempty,gt=0"`
	Type    domain.NotificationType `json:"type" binding:"required,oneof=task_assigned mention task_due project_update system_alert"`
	TaskID  *int                    `json:"taskId" binding:"omitempty,gt=0"`
}
