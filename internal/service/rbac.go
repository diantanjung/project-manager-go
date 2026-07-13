package service

import "project-manager-go/internal/domain"

func HasRole(userRole domain.UserRole, required domain.UserRole) bool {
	return roleLevel(userRole) >= roleLevel(required)
}

func roleLevel(role domain.UserRole) int {
	switch role {
	case domain.RoleAdmin:
		return 4
	case domain.RoleProductOwner:
		return 3
	case domain.RoleProjectManager:
		return 2
	default:
		return 1
	}
}
