package domain

import "time"

type AuthUser struct {
	ID    int      `json:"id"`
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
}

type User struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	Email            string     `json:"email"`
	PasswordHash     string     `json:"-"`
	AvatarStorageKey *string    `json:"avatarStorageKey,omitempty"`
	AvatarURL        *string    `json:"avatarUrl,omitempty"`
	Role             UserRole   `json:"role"`
	CreatedAt        *time.Time `json:"createdAt"`
	UpdatedAt        *time.Time `json:"updatedAt"`
}
