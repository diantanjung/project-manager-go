package postgres

import (
	"context"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func getOne[T any](s *Store, ctx context.Context, query string, args ...any) (T, error) {
	var item T
	err := s.db.GetContext(ctx, &item, query, args...)
	return item, notFound(err)
}

func selectAll[T any](s *Store, ctx context.Context, query string, args ...any) ([]T, error) {
	out := []T{}
	if err := s.db.SelectContext(ctx, &out, query, args...); err != nil {
		return nil, err
	}
	return out, nil
}

func selectPage[T any](
	s *Store,
	ctx context.Context,
	query string,
	filter service.PageFilter,
	total int,
	args ...any,
) (domain.Paginated[T], error) {
	out, err := selectAll[T](s, ctx, query, args...)
	if err != nil {
		return domain.Paginated[T]{}, err
	}
	return paginated(out, filter, total), nil
}
