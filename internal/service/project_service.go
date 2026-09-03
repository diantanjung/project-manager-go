package service

import (
	"context"

	"project-manager-go/internal/domain"
)

func (s *Service) CreateProject(ctx context.Context, actor domain.AuthUser, input ProjectInput) (domain.Project, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.Project{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	if input.Name == nil {
		return domain.Project{}, validationError("Name is required")
	}
	if input.TeamID == nil {
		return domain.Project{}, domain.NewError(domain.ErrValidation, "Team ID is required")
	}
	if err := validatePositiveID(*input.TeamID, "Team ID"); err != nil {
		return domain.Project{}, err
	}
	ok, err := s.store.CanCreateProjectForTeam(ctx, actor, *input.TeamID)
	if err != nil {
		return domain.Project{}, err
	}
	if !ok {
		return domain.Project{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	input.OwnerID = actor.ID
	return s.store.CreateProject(ctx, input)
}

func (s *Service) ListProjects(
	ctx context.Context,
	actor domain.AuthUser,
	filter ProjectFilter,
) (domain.Paginated[domain.Project], error) {
	if err := validateOptionalPositiveID(filter.TeamID, "Team ID"); err != nil {
		return domain.Paginated[domain.Project]{}, err
	}
	filter.PageFilter = NormalizePage(filter.PageFilter)
	return s.store.ListProjects(ctx, actor, filter)
}

func (s *Service) ListSidebarProjects(ctx context.Context, actor domain.AuthUser) ([]domain.SidebarProject, error) {
	return s.store.ListSidebarProjects(ctx, actor)
}

func (s *Service) GetProjectByID(ctx context.Context, actor domain.AuthUser, id int) (domain.Project, error) {
	if err := validatePositiveID(id, "Project ID"); err != nil {
		return domain.Project{}, err
	}
	return s.store.GetProjectByID(ctx, actor, id)
}

func (s *Service) UpdateProject(ctx context.Context, actor domain.AuthUser, id int, input ProjectInput) (domain.Project, error) {
	if err := validatePositiveID(id, "Project ID"); err != nil {
		return domain.Project{}, err
	}
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.Project{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	if input.TeamID != nil {
		if err := validatePositiveID(*input.TeamID, "Team ID"); err != nil {
			return domain.Project{}, err
		}
		ok, err := s.store.CanCreateProjectForTeam(ctx, actor, *input.TeamID)
		if err != nil {
			return domain.Project{}, err
		}
		if !ok {
			return domain.Project{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
		}
	}
	return s.store.UpdateProject(ctx, actor, id, input)
}

func (s *Service) DeleteProject(ctx context.Context, actor domain.AuthUser, id int) (domain.Project, error) {
	if err := validatePositiveID(id, "Project ID"); err != nil {
		return domain.Project{}, err
	}
	if !HasRole(actor.Role, domain.RoleProductOwner) {
		return domain.Project{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.DeleteProject(ctx, actor, id)
}

func (s *Service) ListProjectTasks(
	ctx context.Context,
	actor domain.AuthUser,
	projectID int,
	filter PageFilter,
) (domain.Paginated[domain.Task], error) {
	if err := validatePositiveID(projectID, "Project ID"); err != nil {
		return domain.Paginated[domain.Task]{}, err
	}
	return s.store.ListProjectTasks(ctx, actor, projectID, NormalizePage(filter))
}

func (s *Service) AssignTeamToProject(
	ctx context.Context,
	actor domain.AuthUser,
	input ProjectTeamInput,
) (domain.ProjectTeam, bool, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.ProjectTeam{}, false, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	if err := validatePositiveID(input.ProjectID, "Project ID"); err != nil {
		return domain.ProjectTeam{}, false, err
	}
	if err := validatePositiveID(input.TeamID, "Team ID"); err != nil {
		return domain.ProjectTeam{}, false, err
	}
	if _, err := s.store.GetProjectByID(ctx, actor, input.ProjectID); err != nil {
		return domain.ProjectTeam{}, false, err
	}
	ok, err := s.store.CanCreateProjectForTeam(ctx, actor, input.TeamID)
	if err != nil {
		return domain.ProjectTeam{}, false, err
	}
	if !ok {
		return domain.ProjectTeam{}, false, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.AssignTeamToProject(ctx, actor, input)
}

func (s *Service) ListProjectTeams(ctx context.Context, actor domain.AuthUser, projectID int) ([]domain.ProjectTeam, error) {
	if err := validatePositiveID(projectID, "Project ID"); err != nil {
		return nil, err
	}
	if _, err := s.store.GetProjectByID(ctx, actor, projectID); err != nil {
		return nil, err
	}
	return s.store.ListProjectTeams(ctx, projectID)
}

func (s *Service) RemoveTeamFromProject(ctx context.Context, actor domain.AuthUser, id int) (domain.ProjectTeam, error) {
	if err := validatePositiveID(id, "Project team ID"); err != nil {
		return domain.ProjectTeam{}, err
	}
	if !HasRole(actor.Role, domain.RoleProductOwner) {
		return domain.ProjectTeam{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.RemoveTeamFromProject(ctx, actor, id)
}
