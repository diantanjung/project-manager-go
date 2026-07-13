package postgres

import (
	"context"
	"fmt"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) CanCreateProjectForTeam(ctx context.Context, user domain.AuthUser, teamID int) (bool, error) {
	if service.HasRole(user.Role, domain.RoleProductOwner) {
		return true, nil
	}
	var id int
	err := s.db.QueryRow(ctx, `
		SELECT id FROM team_members
		WHERE team_id = $1 AND user_id = $2 AND role IN ('owner', 'admin')
	`, teamID, user.ID).Scan(&id)
	return exists(err)
}

func (s *Store) CreateProject(ctx context.Context, input service.ProjectInput) (domain.Project, error) {
	row := s.db.QueryRow(ctx, `
		INSERT INTO projects (name, description, team_id, owner_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, team_id, owner_id, created_at, updated_at
	`, value(input.Name), input.Description, value(input.TeamID), input.OwnerID)
	return scanProject(row, false)
}

func (s *Store) ListProjects(
	ctx context.Context,
	user domain.AuthUser,
	filter service.ProjectFilter,
) (domain.Paginated[domain.Project], error) {
	where := projectWhere(user)
	if filter.Search != "" {
		where.add("(p.name ILIKE $%d OR p.description ILIKE $%d)", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.TeamID != nil {
		where.add("p.team_id = $%d", *filter.TeamID)
	}
	totalQuery := "SELECT count(*) FROM projects p " + where.clause()
	var total int
	if err := s.db.QueryRow(ctx, totalQuery, where.values...).Scan(&total); err != nil {
		return domain.Paginated[domain.Project]{}, err
	}

	sortBy := allow(filter.SortBy, "p.created_at", map[string]string{
		"name":      "p.name",
		"createdAt": "p.created_at",
		"updatedAt": "p.updated_at",
	})
	args := append(where.values, filter.Limit, offset(filter.PageFilter))
	query := fmt.Sprintf(`
		SELECT p.id, p.name, p.description, p.team_id, t.name, p.owner_id, u.name, p.created_at, p.updated_at
		FROM projects p
		LEFT JOIN teams t ON t.id = p.team_id
		LEFT JOIN users u ON u.id = p.owner_id
		%s
		ORDER BY %s %s LIMIT $%d OFFSET $%d
	`, where.clause(), sortBy, order(filter.Order), len(args)-1, len(args))
	rows, err := s.db.Query(ctx, query, args...)
	return scanProjectPage(rows, filter.PageFilter, total, err)
}

func (s *Store) GetProjectByID(ctx context.Context, user domain.AuthUser, id int) (domain.Project, error) {
	where := projectWhere(user)
	where.add("p.id = $%d", id)
	query := `
		SELECT p.id, p.name, p.description, p.team_id, t.name, p.owner_id, u.name, p.created_at, p.updated_at
		FROM projects p
		LEFT JOIN teams t ON t.id = p.team_id
		LEFT JOIN users u ON u.id = p.owner_id
		` + where.clause()
	return scanProject(s.db.QueryRow(ctx, query, where.values...), true)
}

func (s *Store) UpdateProject(ctx context.Context, id int, input service.ProjectInput) (domain.Project, error) {
	sets := updateBuilder{}
	sets.add("name", input.Name)
	sets.add("description", input.Description)
	sets.add("team_id", input.TeamID)
	sets.touchUpdatedAt()
	query := fmt.Sprintf(`
		UPDATE projects SET %s WHERE id = $%d
		RETURNING id, name, description, team_id, owner_id, created_at, updated_at
	`, sets.clause(), len(sets.values)+1)
	sets.values = append(sets.values, id)
	return scanProject(s.db.QueryRow(ctx, query, sets.values...), false)
}

func (s *Store) DeleteProject(ctx context.Context, id int) (domain.Project, error) {
	return scanProject(s.db.QueryRow(ctx, `
		DELETE FROM projects WHERE id = $1
		RETURNING id, name, description, team_id, owner_id, created_at, updated_at
	`, id), false)
}

func (s *Store) ListProjectTasks(
	ctx context.Context,
	user domain.AuthUser,
	projectID int,
	filter service.PageFilter,
) (domain.Paginated[domain.Task], error) {
	where := taskWhere(user)
	where.add("ta.project_id = $%d", projectID)
	totalQuery := "SELECT count(*) FROM tasks ta " + where.clause()
	var total int
	if err := s.db.QueryRow(ctx, totalQuery, where.values...).Scan(&total); err != nil {
		return domain.Paginated[domain.Task]{}, err
	}
	args := append(where.values, filter.Limit, offset(filter))
	query := fmt.Sprintf(`
		SELECT ta.id, ta.title, ta.description, ta.status, ta.priority, ta.project_id, p.name,
			ta.creator_id, ta.assignee_id, u.name, u.avatar_url, ta.due_date::text, ta.position,
			ta.created_at, ta.updated_at
		FROM tasks ta
		LEFT JOIN projects p ON p.id = ta.project_id
		LEFT JOIN users u ON u.id = ta.assignee_id
		%s
		ORDER BY ta.created_at DESC LIMIT $%d OFFSET $%d
	`, where.clause(), len(args)-1, len(args))
	rows, err := s.db.Query(ctx, query, args...)
	return s.scanTaskPage(rows, filter, total, err)
}
