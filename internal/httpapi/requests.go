package httpapi

import "project-manager-go/internal/service"

type createTeamRequest struct {
	Name        string  `json:"name" binding:"required,min=2"`
	Description *string `json:"description"`
}

func (r createTeamRequest) input() service.TeamInput {
	return service.TeamInput{
		Name:        &r.Name,
		Description: r.Description,
	}
}

type createProjectRequest struct {
	Name        string  `json:"name" binding:"required,min=2"`
	Description *string `json:"description"`
	TeamID      int     `json:"teamId" binding:"required,gt=0"`
}

func (r createProjectRequest) input() service.ProjectInput {
	return service.ProjectInput{
		Name:        &r.Name,
		Description: r.Description,
		TeamID:      &r.TeamID,
	}
}
