package domain

import "time"

type ActivityLog struct {
	ID         int            `json:"id" db:"id"`
	ActorID    *int           `json:"actorId" db:"actor_id"`
	EntityType string         `json:"entityType" db:"entity_type"`
	EntityID   int            `json:"entityId" db:"entity_id"`
	Action     string         `json:"action" db:"action"`
	Before     map[string]any `json:"before" db:"-"`
	After      map[string]any `json:"after" db:"-"`
	IPAddress  *string        `json:"ipAddress" db:"ip_address"`
	UserAgent  *string        `json:"userAgent" db:"user_agent"`
	CreatedAt  *time.Time     `json:"createdAt" db:"created_at"`
	Actor      *User          `json:"actor,omitempty" db:"-"`
}
