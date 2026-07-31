package domain

import "time"

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
