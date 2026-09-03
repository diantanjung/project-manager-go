package postgres

import (
	"context"
	"fmt"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) CreateTeam(ctx context.Context, input service.TeamInput) (domain.Team, error) {
	return getOne[domain.Team](s, ctx, `
		INSERT INTO teams (name, description) VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`, value(input.Name), input.Description)
}

func (s *Store) ListTeams(
	ctx context.Context,
	user domain.AuthUser,
	filter service.PageFilter,
) (domain.Paginated[domain.Team], error) {
	where := teamWhere(user)
	total, err := s.count(ctx, "teams t", where.clause(), where.values...)
	if err != nil {
		return domain.Paginated[domain.Team]{}, err
	}
	args := append(where.values, filter.Limit, offset(filter))
	query := fmt.Sprintf(`
		SELECT id, name, description, created_at, updated_at
		FROM teams t `+where.clause()+`
		LIMIT $%d OFFSET $%d
	`, len(args)-1, len(args))
	return selectPage[domain.Team](s, ctx, query, filter, total, args...)
}

func (s *Store) GetTeamByID(ctx context.Context, user domain.AuthUser, id int) (domain.Team, error) {
	where := teamWhere(user)
	where.add("t.id = $%d", id)
	query := `
		SELECT id, name, description, created_at, updated_at FROM teams t ` + where.clause()
	return getOne[domain.Team](s, ctx, query, where.values...)
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
	return getOne[domain.Team](s, ctx, query, sets.values...)
}

func (s *Store) DeleteTeam(ctx context.Context, id int) (domain.Team, error) {
	return getOne[domain.Team](s, ctx, `
		DELETE FROM teams WHERE id = $1
		RETURNING id, name, description, created_at, updated_at
	`, id)
}

func (s *Store) ListTeamMembers(ctx context.Context, user domain.AuthUser, teamID int) ([]domain.TeamMember, error) {
	if _, err := s.GetTeamByID(ctx, user, teamID); err != nil {
		return nil, err
	}
	return selectAll[domain.TeamMember](s, ctx, `
		SELECT tm.id, tm.team_id, tm.user_id, u.name AS user_name, u.email AS user_email, tm.role, tm.joined_at
		FROM team_members tm
		LEFT JOIN users u ON u.id = tm.user_id
		WHERE tm.team_id = $1
	`, teamID)
}

func (s *Store) AddTeamMember(
	ctx context.Context,
	teamID int,
	input service.TeamMemberInput,
) (domain.TeamMember, error) {
	var existing int
	err := s.db.GetContext(ctx, &existing, `
		SELECT id FROM team_members WHERE team_id = $1 AND user_id = $2
	`, teamID, input.UserID)
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

	item, err := getOne[domain.TeamMember](s, ctx, `
		INSERT INTO team_members (team_id, user_id, role)
		VALUES ($1, $2, $3)
		RETURNING id, team_id, user_id, role, joined_at
	`, teamID, input.UserID, role)
	return item, err
}

func (s *Store) RemoveTeamMember(ctx context.Context, teamID int, userID int) (domain.TeamMember, error) {
	return getOne[domain.TeamMember](s, ctx, `
		DELETE FROM team_members WHERE team_id = $1 AND user_id = $2
		RETURNING id, team_id, user_id, role, joined_at
	`, teamID, userID)
}
