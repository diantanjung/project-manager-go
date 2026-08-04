package domain

import "time"

type AuthUser struct {
	ID    int      `json:"id"`
	Email string   `json:"email"`
	Role  UserRole `json:"role"`
}

type User struct {
	ID               int        `json:"id" db:"id"`
	Name             string     `json:"name" db:"name"`
	Email            string     `json:"email" db:"email"`
	PasswordHash     string     `json:"-" db:"password"`
	AvatarStorageKey *string    `json:"avatarStorageKey,omitempty" db:"avatar_storage_key"`
	AvatarURL        *string    `json:"avatarUrl,omitempty" db:"avatar_url"`
	Role             UserRole   `json:"role" db:"role"`
	CreatedAt        *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt        *time.Time `json:"updatedAt" db:"updated_at"`
}
