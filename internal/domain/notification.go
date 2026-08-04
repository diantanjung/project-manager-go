package domain

import "time"

type Notification struct {
	ID             int              `json:"id" db:"id"`
	UserID         int              `json:"userId" db:"user_id"`
	ActorID        *int             `json:"actorId" db:"actor_id"`
	ActorName      *string          `json:"actorName,omitempty" db:"actor_name"`
	ActorAvatarURL *string          `json:"actorAvatarUrl,omitempty" db:"actor_avatar_url"`
	Type           NotificationType `json:"type" db:"type"`
	TaskID         *int             `json:"taskId" db:"task_id"`
	IsRead         bool             `json:"isRead" db:"is_read"`
	CreatedAt      *time.Time       `json:"createdAt" db:"created_at"`
}
