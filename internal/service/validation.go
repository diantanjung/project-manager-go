package service

import (
	"project-manager-go/internal/domain"
)

func validationError(message string) error {
	return domain.NewError(domain.ErrValidation, message)
}

func validatePositiveID(value int, field string) error {
	if value < 1 {
		return validationError(field + " is required")
	}
	return nil
}

func validateOptionalPositiveID(value *int, field string) error {
	if value == nil {
		return nil
	}
	return validatePositiveID(*value, field)
}

func validateUserRole(value *domain.UserRole) error {
	if value != nil && !value.IsValid() {
		return validationError("Role is invalid")
	}
	return nil
}

func validateTeamMemberRole(value *domain.TeamMemberRole) error {
	if value != nil && !value.IsValid() {
		return validationError("Member role is invalid")
	}
	return nil
}

func validateTaskStatus(value *domain.TaskStatus) error {
	if value != nil && !value.IsValid() {
		return validationError("Task status is invalid")
	}
	return nil
}

func validateTaskPriority(value *domain.TaskPriority) error {
	if value != nil && !value.IsValid() {
		return validationError("Task priority is invalid")
	}
	return nil
}
