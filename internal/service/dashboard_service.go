package service

import (
	"context"

	"project-manager-go/internal/domain"
)

func (s *Service) DashboardSummary(ctx context.Context, actor domain.AuthUser) (domain.DashboardSummary, error) {
	return s.store.DashboardSummary(ctx, actor)
}
