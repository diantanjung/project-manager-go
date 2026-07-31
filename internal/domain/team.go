package domain

import "time"

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
