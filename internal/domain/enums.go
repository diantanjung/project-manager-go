package domain

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

type ProjectStatus string

const (
	ProjectStatusPlanning  ProjectStatus = "planning"
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusPaused    ProjectStatus = "paused"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusArchived  ProjectStatus = "archived"
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
	NotificationTaskAssigned  NotificationType = "task_assigned"
	NotificationMention       NotificationType = "mention"
	NotificationTaskDue       NotificationType = "task_due"
	NotificationProjectUpdate NotificationType = "project_update"
	NotificationSystemAlert   NotificationType = "system_alert"
)

type WebhookEvent string

const (
	WebhookEventTaskCreated    WebhookEvent = "task.created"
	WebhookEventTaskUpdated    WebhookEvent = "task.updated"
	WebhookEventTaskCompleted  WebhookEvent = "task.completed"
	WebhookEventCommentCreated WebhookEvent = "comment.created"
	WebhookEventProjectUpdated WebhookEvent = "project.updated"
)

type WebhookStatus string

const (
	WebhookStatusPending   WebhookStatus = "pending"
	WebhookStatusDelivered WebhookStatus = "delivered"
	WebhookStatusFailed    WebhookStatus = "failed"
)

type ExportStatus string

const (
	ExportStatusPending    ExportStatus = "pending"
	ExportStatusProcessing ExportStatus = "processing"
	ExportStatusCompleted  ExportStatus = "completed"
	ExportStatusFailed     ExportStatus = "failed"
)
