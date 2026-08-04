package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) count(ctx context.Context, table string, clause string, args ...any) (int, error) {
	var total int
	query := fmt.Sprintf("SELECT count(*) FROM %s %s", table, clause)
	err := s.db.GetContext(ctx, &total, query, args...)
	return total, err
}

type whereBuilder struct {
	parts  []string
	values []any
}

func (w *whereBuilder) add(format string, values ...any) {
	next := len(w.values) + 1
	part := format
	for strings.Contains(part, "$%d") {
		part = strings.Replace(part, "$%d", fmt.Sprintf("$%d", next), 1)
		next++
	}
	w.parts = append(w.parts, part)
	w.values = append(w.values, values...)
}

func (w whereBuilder) clause() string {
	if len(w.parts) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(w.parts, " AND ")
}

type updateBuilder struct {
	parts  []string
	values []any
}

func (u *updateBuilder) add(column string, value any) {
	if value == nil {
		return
	}
	u.values = append(u.values, value)
	u.parts = append(u.parts, fmt.Sprintf("%s = $%d", column, len(u.values)))
}

func (u *updateBuilder) touchUpdatedAt() {
	u.parts = append(u.parts, "updated_at = NOW()")
}

func (u updateBuilder) clause() string {
	return strings.Join(u.parts, ", ")
}

func allow(value string, fallback string, allowed map[string]string) string {
	if got, ok := allowed[value]; ok {
		return got
	}
	return fallback
}

func order(value string) string {
	if value == "asc" {
		return "ASC"
	}
	return "DESC"
}

func offset(filter service.PageFilter) int {
	return (filter.Page - 1) * filter.Limit
}

func paginated[T any](data []T, filter service.PageFilter, total int) domain.Paginated[T] {
	return domain.Paginated[T]{
		Data: data,
		Pagination: domain.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalItems: total,
			TotalPages: service.TotalPages(total, filter.Limit),
		},
	}
}

func value[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

func notFound(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func exists(err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, err
}

func projectWhere(user domain.AuthUser) whereBuilder {
	where := whereBuilder{}
	if !service.HasRole(user.Role, domain.RoleProductOwner) {
		where.add(`(
			p.owner_id = $%d OR
			EXISTS (
				SELECT 1
				FROM team_members primary_tm
				WHERE primary_tm.team_id = p.team_id AND primary_tm.user_id = $%d
			) OR
			EXISTS (
				SELECT 1
				FROM project_teams access_pt
				JOIN team_members additional_tm ON additional_tm.team_id = access_pt.team_id
				WHERE access_pt.project_id = p.id AND additional_tm.user_id = $%d
			)
		)`, user.ID, user.ID, user.ID)
	}
	return where
}

func taskWhere(user domain.AuthUser) whereBuilder {
	where := whereBuilder{}
	if !service.HasRole(user.Role, domain.RoleProductOwner) {
		where.add(`(
			ta.creator_id = $%d OR
			ta.assignee_id = $%d OR
			EXISTS (
				SELECT 1
				FROM task_assignments access_ta
				WHERE access_ta.task_id = ta.id AND access_ta.user_id = $%d
			) OR
			EXISTS (
				SELECT 1
				FROM projects access_p
				WHERE access_p.id = ta.project_id AND access_p.owner_id = $%d
			) OR
			EXISTS (
				SELECT 1
				FROM projects access_p
				JOIN team_members primary_tm ON primary_tm.team_id = access_p.team_id
				WHERE access_p.id = ta.project_id AND primary_tm.user_id = $%d
			) OR
			EXISTS (
				SELECT 1
				FROM projects access_p
				JOIN project_teams access_pt ON access_pt.project_id = access_p.id
				JOIN team_members additional_tm ON additional_tm.team_id = access_pt.team_id
				WHERE access_p.id = ta.project_id AND additional_tm.user_id = $%d
			)
		)`, user.ID, user.ID, user.ID, user.ID, user.ID, user.ID)
	}
	return where
}
