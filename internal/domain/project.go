package domain

import "time"

type Project struct {
	ID          int        `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description" db:"description"`
	TeamID      int        `json:"teamId" db:"team_id"`
	TeamName    *string    `json:"teamName,omitempty" db:"team_name"`
	OwnerID     int        `json:"ownerId" db:"owner_id"`
	OwnerName   *string    `json:"ownerName,omitempty" db:"owner_name"`
	CreatedAt   *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   *time.Time `json:"updatedAt" db:"updated_at"`
}

type SidebarProject struct {
	ID            int    `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	OpenTaskCount int    `json:"openTaskCount" db:"open_task_count"`
}

type ProjectTeam struct {
	ID              int        `json:"id" db:"id"`
	ProjectID       int        `json:"projectId" db:"project_id"`
	TeamID          int        `json:"teamId" db:"team_id"`
	TeamName        *string    `json:"teamName,omitempty" db:"team_name"`
	TeamDescription *string    `json:"teamDescription,omitempty" db:"team_description"`
	AssignedAt      *time.Time `json:"assignedAt" db:"assigned_at"`
}
