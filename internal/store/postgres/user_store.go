package postgres

import (
	"context"
	"fmt"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) CreateUser(ctx context.Context, input service.CreateUserInput) (domain.User, error) {
	role := domain.RoleTeamMember
	if input.Role != nil {
		role = *input.Role
	}
	return getOne[domain.User](s, ctx, `
		INSERT INTO users (name, email, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, avatar_storage_key, role, created_at, updated_at
	`, input.Name, input.Email, input.Password, role)
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return getOne[domain.User](s, ctx, `
		SELECT id, name, email, password, avatar_storage_key, role, created_at, updated_at
		FROM users WHERE email = $1
	`, email)
}

func (s *Store) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	return getOne[domain.User](s, ctx, `
		SELECT id, name, email, avatar_storage_key, role, created_at, updated_at
		FROM users WHERE id = $1
	`, id)
}

func (s *Store) ListUsers(ctx context.Context, filter service.ListUsersFilter) (domain.Paginated[domain.User], error) {
	sortBy := allow(filter.SortBy, "created_at", map[string]string{
		"name":      "name",
		"email":     "email",
		"role":      "role",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	})
	order := order(filter.Order)
	where := whereBuilder{}
	if filter.Search != "" {
		where.add("(name ILIKE $%d OR email ILIKE $%d)", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.Role != nil {
		where.add("role = $%d", *filter.Role)
	}

	total, err := s.count(ctx, "users", where.clause(), where.values...)
	if err != nil {
		return domain.Paginated[domain.User]{}, err
	}

	args := append(where.values, filter.Limit, offset(filter.PageFilter))
	query := fmt.Sprintf(`
		SELECT id, name, email, avatar_storage_key, role, created_at, updated_at
		FROM users %s
		ORDER BY %s %s LIMIT $%d OFFSET $%d
	`, where.clause(), sortBy, order, len(args)-1, len(args))
	return selectPage[domain.User](s, ctx, query, filter.PageFilter, total, args...)
}

func (s *Store) UpdateUser(ctx context.Context, id int, input service.UpdateUserInput) (domain.User, error) {
	sets := updateBuilder{}
	sets.add("name", input.Name)
	sets.add("email", input.Email)
	sets.add("password", input.Password)
	if input.AvatarStorageKey != nil {
		sets.add("avatar_storage_key", input.AvatarStorageKey)
	} else {
		sets.add("avatar_storage_key", input.LegacyAvatarURL)
	}
	sets.add("role", input.Role)
	sets.touchUpdatedAt()
	query := fmt.Sprintf(`
		UPDATE users SET %s WHERE id = $%d
		RETURNING id, name, email, avatar_storage_key, role, created_at, updated_at
	`, sets.clause(), len(sets.values)+1)
	sets.values = append(sets.values, id)
	return getOne[domain.User](s, ctx, query, sets.values...)
}

func (s *Store) DeleteUser(ctx context.Context, id int) (domain.User, error) {
	return getOne[domain.User](s, ctx, `
		DELETE FROM users WHERE id = $1
		RETURNING id, name, email, avatar_storage_key, role, created_at, updated_at
	`, id)
}

func (s *Store) ListUserTasks(
	ctx context.Context,
	userID int,
	filter service.PageFilter,
) (domain.Paginated[domain.Task], error) {
	total, err := s.count(ctx, "tasks", "WHERE creator_id = $1 OR assignee_id = $1", userID)
	if err != nil {
		return domain.Paginated[domain.Task]{}, err
	}
	query := `
		SELECT t.id, t.title, t.description, t.status, t.priority, t.project_id, p.name AS project_name,
			t.creator_id, t.assignee_id, NULL::text AS assignee_name, NULL::text AS assignee_avatar_url, t.due_date::text AS due_date, t.position,
			t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN projects p ON p.id = t.project_id
		WHERE t.creator_id = $1 OR t.assignee_id = $1
		ORDER BY t.created_at DESC LIMIT $2 OFFSET $3
	`
	return selectPage[domain.Task](s, ctx, query, filter, total, userID, filter.Limit, offset(filter))
}
