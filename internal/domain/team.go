package domain

import "time"

type Team struct {
	ID          int        `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description" db:"description"`
	CreatedAt   *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   *time.Time `json:"updatedAt" db:"updated_at"`
}

type TeamMember struct {
	ID        int             `json:"id" db:"id"`
	TeamID    int             `json:"teamId,omitempty" db:"team_id"`
	UserID    int             `json:"userId" db:"user_id"`
	UserName  *string         `json:"userName,omitempty" db:"user_name"`
	UserEmail *string         `json:"userEmail,omitempty" db:"user_email"`
	Role      *TeamMemberRole `json:"role" db:"role"`
	JoinedAt  *time.Time      `json:"joinedAt" db:"joined_at"`
}
