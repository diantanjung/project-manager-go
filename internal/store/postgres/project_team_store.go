package postgres

import (
	"context"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) AssignTeamToProject(
	ctx context.Context,
	user domain.AuthUser,
	input service.ProjectTeamInput,
) (domain.ProjectTeam, bool, error) {
	if _, err := s.GetProjectByID(ctx, user, input.ProjectID); err != nil {
		return domain.ProjectTeam{}, false, err
	}

	var id int
	err := s.db.GetContext(ctx, &id, `
		SELECT id FROM project_teams WHERE project_id = $1 AND team_id = $2
	`, input.ProjectID, input.TeamID)
	ok, err := exists(err)
	if err != nil {
		return domain.ProjectTeam{}, false, err
	}
	if ok {
		item, err := s.GetProjectTeamByID(ctx, id)
		return item, true, err
	}

	item, err := getOne[domain.ProjectTeam](s, ctx, `
		INSERT INTO project_teams (project_id, team_id)
		VALUES ($1, $2)
		RETURNING id, project_id, team_id, assigned_at
	`, input.ProjectID, input.TeamID)
	return item, false, err
}

func (s *Store) ListProjectTeams(ctx context.Context, projectID int) ([]domain.ProjectTeam, error) {
	return selectAll[domain.ProjectTeam](s, ctx, `
		SELECT pt.id, pt.project_id, pt.team_id, t.name AS team_name, t.description AS team_description, pt.assigned_at
		FROM project_teams pt
		LEFT JOIN teams t ON t.id = pt.team_id
		WHERE pt.project_id = $1
	`, projectID)
}

func (s *Store) GetProjectTeamByID(ctx context.Context, id int) (domain.ProjectTeam, error) {
	return getOne[domain.ProjectTeam](s, ctx, `
		SELECT id, project_id, team_id, assigned_at FROM project_teams WHERE id = $1
	`, id)
}

func (s *Store) RemoveTeamFromProject(ctx context.Context, user domain.AuthUser, id int) (domain.ProjectTeam, error) {
	item, err := s.GetProjectTeamByID(ctx, id)
	if err != nil {
		return domain.ProjectTeam{}, err
	}
	if _, err := s.GetProjectByID(ctx, user, item.ProjectID); err != nil {
		return domain.ProjectTeam{}, err
	}
	return getOne[domain.ProjectTeam](s, ctx, `
		DELETE FROM project_teams WHERE id = $1
		RETURNING id, project_id, team_id, assigned_at
	`, id)
}
