package service

import (
	"context"
	"testing"

	"project-manager-go/internal/domain"
)

func TestDashboardSummaryUsesAuthenticatedUser(t *testing.T) {
	store := newMemoryStore(t)
	store.dashboardSummary = domain.DashboardSummary{
		TotalActiveProjects: 2,
		TaskCountPerStatus: map[domain.TaskStatus]int{
			domain.TaskStatusTodo:       3,
			domain.TaskStatusInProgress: 2,
		},
		ActiveProgress: domain.ActiveProgress{
			Doing:      2,
			Todo:       3,
			Total:      5,
			Ratio:      0.4,
			Percentage: 40,
		},
	}
	svc := New(store, testTokenManager())
	actor := domain.AuthUser{ID: 7, Role: domain.RoleTeamMember}

	summary, err := svc.DashboardSummary(context.Background(), actor)
	if err != nil {
		t.Fatalf("DashboardSummary() error = %v", err)
	}
	if summary.TotalActiveProjects != 2 {
		t.Fatalf("TotalActiveProjects = %d, want 2", summary.TotalActiveProjects)
	}
	if store.dashboardUser.ID != actor.ID || store.dashboardUser.Role != actor.Role {
		t.Fatalf("dashboard user = %+v, want %+v", store.dashboardUser, actor)
	}
}
