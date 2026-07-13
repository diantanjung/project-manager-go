package service

import (
	"context"

	"project-manager-go/internal/domain"
)

func (s *Service) CreateTeam(ctx context.Context, actor domain.AuthUser, input TeamInput) (domain.Team, error) {
	if !HasRole(actor.Role, domain.RoleProductOwner) {
		return domain.Team{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.CreateTeam(ctx, input)
}

func (s *Service) ListTeams(ctx context.Context, filter PageFilter) (domain.Paginated[domain.Team], error) {
	return s.store.ListTeams(ctx, NormalizePage(filter))
}

func (s *Service) GetTeamByID(ctx context.Context, id int) (domain.Team, error) {
	return s.store.GetTeamByID(ctx, id)
}

func (s *Service) UpdateTeam(ctx context.Context, actor domain.AuthUser, id int, input TeamInput) (domain.Team, error) {
	if !HasRole(actor.Role, domain.RoleProductOwner) {
		return domain.Team{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.UpdateTeam(ctx, id, input)
}

func (s *Service) DeleteTeam(ctx context.Context, actor domain.AuthUser, id int) (domain.Team, error) {
	if actor.Role != domain.RoleAdmin {
		return domain.Team{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.DeleteTeam(ctx, id)
}

func (s *Service) ListTeamMembers(ctx context.Context, teamID int) ([]domain.TeamMember, error) {
	return s.store.ListTeamMembers(ctx, teamID)
}

func (s *Service) AddTeamMember(
	ctx context.Context,
	actor domain.AuthUser,
	teamID int,
	input TeamMemberInput,
) (domain.TeamMember, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.TeamMember{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.AddTeamMember(ctx, teamID, input)
}

func (s *Service) RemoveTeamMember(
	ctx context.Context,
	actor domain.AuthUser,
	teamID int,
	userID int,
) (domain.TeamMember, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.TeamMember{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.RemoveTeamMember(ctx, teamID, userID)
}
