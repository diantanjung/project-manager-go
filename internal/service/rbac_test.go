package service

import (
	"testing"

	"project-manager-go/internal/domain"
)

func TestHasRole(t *testing.T) {
	tests := []struct {
		name     string
		userRole domain.UserRole
		required domain.UserRole
		want     bool
	}{
		{name: "admin has product owner", userRole: domain.RoleAdmin, required: domain.RoleProductOwner, want: true},
		{name: "team member lacks project manager", userRole: domain.RoleTeamMember, required: domain.RoleProjectManager, want: false},
		{name: "exact role", userRole: domain.RoleProjectManager, required: domain.RoleProjectManager, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasRole(tt.userRole, tt.required); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
