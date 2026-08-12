package domain

type UserRole string

const (
	RoleAdmin          UserRole = "admin"
	RoleProductOwner   UserRole = "productOwner"
	RoleProjectManager UserRole = "projectManager"
	RoleTeamMember     UserRole = "teamMember"
)

func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleProductOwner, RoleProjectManager, RoleTeamMember:
		return true
	default:
		return false
	}
}

type TeamMemberRole string

const (
	TeamRoleOwner  TeamMemberRole = "owner"
	TeamRoleAdmin  TeamMemberRole = "admin"
	TeamRoleMember TeamMemberRole = "member"
)

func (r TeamMemberRole) IsValid() bool {
	switch r {
	case TeamRoleOwner, TeamRoleAdmin, TeamRoleMember:
		return true
	default:
		return false
	}
}

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

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusBacklog, TaskStatusTodo, TaskStatusInProgress, TaskStatusReview, TaskStatusDone:
		return true
	default:
		return false
	}
}

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
	TaskPriorityUrgent TaskPriority = "urgent"
)

func (p TaskPriority) IsValid() bool {
	switch p {
	case TaskPriorityLow, TaskPriorityMedium, TaskPriorityHigh, TaskPriorityUrgent:
		return true
	default:
		return false
	}
}

type NotificationType string

const (
	NotificationTaskAssigned  NotificationType = "task_assigned"
	NotificationMention       NotificationType = "mention"
	NotificationTaskDue       NotificationType = "task_due"
	NotificationProjectUpdate NotificationType = "project_update"
	NotificationSystemAlert   NotificationType = "system_alert"
)

func (t NotificationType) IsValid() bool {
	switch t {
	case NotificationTaskAssigned, NotificationMention, NotificationTaskDue, NotificationProjectUpdate, NotificationSystemAlert:
		return true
	default:
		return false
	}
}

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
