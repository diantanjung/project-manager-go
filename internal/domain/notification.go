package domain

import "time"

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
