package postgres

import (
	"context"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) AssignTeamToProject(
	ctx context.Context,
	input service.ProjectTeamInput,
) (domain.ProjectTeam, bool, error) {
	var id int
	err := s.db.QueryRow(ctx, `
		SELECT id FROM project_teams WHERE project_id = $1 AND team_id = $2
	`, input.ProjectID, input.TeamID).Scan(&id)
	ok, err := exists(err)
	if err != nil {
		return domain.ProjectTeam{}, false, err
	}
	if ok {
		item, err := s.GetProjectTeamByID(ctx, id)
		return item, true, err
	}

	var item domain.ProjectTeam
	err = s.db.QueryRow(ctx, `
		INSERT INTO project_teams (project_id, team_id)
		VALUES ($1, $2)
		RETURNING id, project_id, team_id, assigned_at
	`, input.ProjectID, input.TeamID).Scan(&item.ID, &item.ProjectID, &item.TeamID, &item.AssignedAt)
	return item, false, notFound(err)
}

func (s *Store) ListProjectTeams(ctx context.Context, projectID int) ([]domain.ProjectTeam, error) {
	rows, err := s.db.Query(ctx, `
		SELECT pt.id, pt.project_id, pt.team_id, t.name, t.description, pt.assigned_at
		FROM project_teams pt
		LEFT JOIN teams t ON t.id = pt.team_id
		WHERE pt.project_id = $1
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.ProjectTeam{}
	for rows.Next() {
		var item domain.ProjectTeam
		if err := rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.TeamID,
			&item.TeamName,
			&item.TeamDescription,
			&item.AssignedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) GetProjectTeamByID(ctx context.Context, id int) (domain.ProjectTeam, error) {
	var item domain.ProjectTeam
	err := s.db.QueryRow(ctx, `
		SELECT id, project_id, team_id, assigned_at FROM project_teams WHERE id = $1
	`, id).Scan(&item.ID, &item.ProjectID, &item.TeamID, &item.AssignedAt)
	return item, notFound(err)
}

func (s *Store) RemoveTeamFromProject(ctx context.Context, id int) (domain.ProjectTeam, error) {
	var item domain.ProjectTeam
	err := s.db.QueryRow(ctx, `
		DELETE FROM project_teams WHERE id = $1
		RETURNING id, project_id, team_id, assigned_at
	`, id).Scan(&item.ID, &item.ProjectID, &item.TeamID, &item.AssignedAt)
	return item, notFound(err)
}
