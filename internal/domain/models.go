package domain

import "time"

type UserRole string

const (
	RoleAdmin          UserRole = "admin"
	RoleProductOwner   UserRole = "productOwner"
	RoleProjectManager UserRole = "projectManager"
	RoleTeamMember     UserRole = "teamMember"
)

type TeamMemberRole string

const (
	TeamRoleOwner  TeamMemberRole = "owner"
	TeamRoleAdmin  TeamMemberRole = "admin"
	TeamRoleMember TeamMemberRole = "member"
)

type TaskStatus string

const (
	TaskStatusBacklog    TaskStatus = "backlog"
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusReview     TaskStatus = "review"
	TaskStatusDone       TaskStatus = "done"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
	TaskPriorityUrgent TaskPriority = "urgent"
)

type NotificationType string

const (
	NotificationTaskAssigned NotificationType = "task_assigned"
	NotificationMention      NotificationType = "mention"
	NotificationSystemAlert  NotificationType = "system_alert"
)

type AuthUser struct {
	ID    int      `json:"id"`
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
}

type User struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	AvatarURL    *string    `json:"avatarUrl,omitempty"`
	Role         UserRole   `json:"role"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type Team struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type TeamMember struct {
	ID        int             `json:"id"`
	TeamID    int             `json:"teamId,omitempty"`
	UserID    int             `json:"userId"`
	UserName  *string         `json:"userName,omitempty"`
	UserEmail *string         `json:"userEmail,omitempty"`
	Role      *TeamMemberRole `json:"role"`
	JoinedAt  *time.Time      `json:"joinedAt"`
}

type Project struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	TeamID      int        `json:"teamId"`
	TeamName    *string    `json:"teamName,omitempty"`
	OwnerID     int        `json:"ownerId"`
	OwnerName   *string    `json:"ownerName,omitempty"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type ProjectTeam struct {
	ID              int        `json:"id"`
	ProjectID       int        `json:"projectId"`
	TeamID          int        `json:"teamId"`
	TeamName        *string    `json:"teamName,omitempty"`
	TeamDescription *string    `json:"teamDescription,omitempty"`
	AssignedAt      *time.Time `json:"assignedAt"`
}

type Task struct {
	ID                int           `json:"id"`
	Title             string        `json:"title"`
	Description       *string       `json:"description"`
	Status            TaskStatus    `json:"status"`
	Priority          *TaskPriority `json:"priority"`
	ProjectID         int           `json:"projectId"`
	ProjectName       *string       `json:"projectName,omitempty"`
	CreatorID         int           `json:"creatorId"`
	AssigneeID        int           `json:"assigneeId"`
	AssigneeName      *string       `json:"assigneeName,omitempty"`
	AssigneeAvatarURL *string       `json:"assigneeAvatarUrl,omitempty"`
	DueDate           *string       `json:"dueDate"`
	Position          *int          `json:"position"`
	CreatedAt         *time.Time    `json:"createdAt"`
	UpdatedAt         *time.Time    `json:"updatedAt"`
	Comments          []CommentView `json:"comments,omitempty"`
	Attachments       []Attachment  `json:"attachments,omitempty"`
}

type TaskAssignment struct {
	ID            int        `json:"id"`
	TaskID        int        `json:"taskId"`
	UserID        int        `json:"userId"`
	UserName      *string    `json:"userName,omitempty"`
	UserEmail     *string    `json:"userEmail,omitempty"`
	UserAvatarURL *string    `json:"userAvatarUrl,omitempty"`
	AssignedAt    *time.Time `json:"assignedAt"`
}

type Comment struct {
	ID        int        `json:"id"`
	Content   string     `json:"content"`
	TaskID    int        `json:"taskId"`
	AuthorID  int        `json:"authorId"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type CommentView struct {
	ID              int        `json:"id"`
	Content         string     `json:"content"`
	TaskID          int        `json:"taskId,omitempty"`
	AuthorID        int        `json:"authorId"`
	AuthorName      *string    `json:"authorName,omitempty"`
	AuthorAvatarURL *string    `json:"authorAvatarUrl,omitempty"`
	CreatedAt       *time.Time `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt"`
}

type Attachment struct {
	ID         int        `json:"id"`
	FileName   string     `json:"fileName"`
	FileURL    string     `json:"fileUrl"`
	FileSize   *int       `json:"fileSize"`
	MimeType   *string    `json:"mimeType"`
	TaskID     int        `json:"taskId"`
	UploaderID int        `json:"uploaderId"`
	CreatedAt  *time.Time `json:"createdAt"`
}

type Notification struct {
	ID             int              `json:"id"`
	UserID         int              `json:"userId"`
	ActorID        *int             `json:"actorId"`
	ActorName      *string          `json:"actorName,omitempty"`
	ActorAvatarURL *string          `json:"actorAvatarUrl,omitempty"`
	Type           NotificationType `json:"type"`
	TaskID         *int             `json:"taskId"`
	IsRead         bool             `json:"isRead"`
	CreatedAt      *time.Time       `json:"createdAt"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type Paginated[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}
