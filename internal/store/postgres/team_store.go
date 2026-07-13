package postgres

import (
	"context"
	"fmt"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) CreateTeam(ctx context.Context, input service.TeamInput) (domain.Team, error) {
	row := s.db.QueryRow(ctx, `
		INSERT INTO teams (name, description) VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`, value(input.Name), input.Description)
	return scanTeam(row)
}

func (s *Store) ListTeams(ctx context.Context, filter service.PageFilter) (domain.Paginated[domain.Team], error) {
	total, err := s.count(ctx, "teams", "")
	if err != nil {
		return domain.Paginated[domain.Team]{}, err
	}
	rows, err := s.db.Query(ctx, `
		SELECT id, name, description, created_at, updated_at
		FROM teams LIMIT $1 OFFSET $2
	`, filter.Limit, offset(filter))
	if err != nil {
		return domain.Paginated[domain.Team]{}, err
	}
	defer rows.Close()

	teams := []domain.Team{}
	for rows.Next() {
		team, err := scanTeam(rows)
		if err != nil {
			return domain.Paginated[domain.Team]{}, err
		}
		teams = append(teams, team)
	}
	return paginated(teams, filter, total), rows.Err()
}

func (s *Store) GetTeamByID(ctx context.Context, id int) (domain.Team, error) {
	return scanTeam(s.db.QueryRow(ctx, `
		SELECT id, name, description, created_at, updated_at FROM teams WHERE id = $1
	`, id))
}

func (s *Store) UpdateTeam(ctx context.Context, id int, input service.TeamInput) (domain.Team, error) {
	sets := updateBuilder{}
	sets.add("name", input.Name)
	sets.add("description", input.Description)
	sets.touchUpdatedAt()
	query := fmt.Sprintf(`
		UPDATE teams SET %s WHERE id = $%d
		RETURNING id, name, description, created_at, updated_at
	`, sets.clause(), len(sets.values)+1)
	sets.values = append(sets.values, id)
	return scanTeam(s.db.QueryRow(ctx, query, sets.values...))
}

func (s *Store) DeleteTeam(ctx context.Context, id int) (domain.Team, error) {
	return scanTeam(s.db.QueryRow(ctx, `
		DELETE FROM teams WHERE id = $1
		RETURNING id, name, description, created_at, updated_at
	`, id))
}

func (s *Store) ListTeamMembers(ctx context.Context, teamID int) ([]domain.TeamMember, error) {
	rows, err := s.db.Query(ctx, `
		SELECT tm.id, tm.team_id, tm.user_id, u.name, u.email, tm.role, tm.joined_at
		FROM team_members tm
		LEFT JOIN users u ON u.id = tm.user_id
		WHERE tm.team_id = $1
	`, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.TeamMember{}
	for rows.Next() {
		var item domain.TeamMember
		if err := rows.Scan(
			&item.ID,
			&item.TeamID,
			&item.UserID,
			&item.UserName,
			&item.UserEmail,
			&item.Role,
			&item.JoinedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) AddTeamMember(
	ctx context.Context,
	teamID int,
	input service.TeamMemberInput,
) (domain.TeamMember, error) {
	var existing int
	err := s.db.QueryRow(ctx, `
		SELECT id FROM team_members WHERE team_id = $1 AND user_id = $2
	`, teamID, input.UserID).Scan(&existing)
	ok, err := exists(err)
	if err != nil {
		return domain.TeamMember{}, err
	}
	if ok {
		return domain.TeamMember{}, domain.NewError(domain.ErrDuplicate, "User is already a member of this team")
	}

	role := domain.TeamRoleMember
	if input.Role != nil {
		role = *input.Role
	}

	var item domain.TeamMember
	err = s.db.QueryRow(ctx, `
		INSERT INTO team_members (team_id, user_id, role)
		VALUES ($1, $2, $3)
		RETURNING id, team_id, user_id, role, joined_at
	`, teamID, input.UserID, role).Scan(&item.ID, &item.TeamID, &item.UserID, &item.Role, &item.JoinedAt)
	return item, notFound(err)
}

func (s *Store) RemoveTeamMember(ctx context.Context, teamID int, userID int) (domain.TeamMember, error) {
	var item domain.TeamMember
	err := s.db.QueryRow(ctx, `
		DELETE FROM team_members WHERE team_id = $1 AND user_id = $2
		RETURNING id, team_id, user_id, role, joined_at
	`, teamID, userID).Scan(&item.ID, &item.TeamID, &item.UserID, &item.Role, &item.JoinedAt)
	return item, notFound(err)
}
