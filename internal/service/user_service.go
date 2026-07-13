package service

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"project-manager-go/internal/domain"
)

func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (domain.User, error) {
	if err := validateCreateUser(input); err != nil {
		return domain.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return domain.User{}, fmt.Errorf("hashing password: %w", err)
	}

	input.Password = string(hash)
	return s.store.CreateUser(ctx, input)
}

func (s *Service) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	user, err := s.store.GetUserByID(ctx, id)
	user.PasswordHash = ""
	return user, err
}

func (s *Service) ListUsers(ctx context.Context, filter ListUsersFilter) (domain.Paginated[domain.User], error) {
	filter.PageFilter = NormalizePage(filter.PageFilter)
	return s.store.ListUsers(ctx, filter)
}

func (s *Service) UpdateUser(ctx context.Context, actor domain.AuthUser, id int, input UpdateUserInput) (domain.User, error) {
	if actor.ID != id && !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.User{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	if input.Password != nil && *input.Password != "" {
		if len(*input.Password) < 6 {
			return domain.User{}, domain.NewError(domain.ErrValidation, "Password must be at least 6 characters")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), 10)
		if err != nil {
			return domain.User{}, fmt.Errorf("hashing password: %w", err)
		}
		hashed := string(hash)
		input.Password = &hashed
	}
	if input.Role != nil && actor.Role != domain.RoleAdmin {
		return domain.User{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.UpdateUser(ctx, id, input)
}

func (s *Service) DeleteUser(ctx context.Context, actor domain.AuthUser, id int) (domain.User, error) {
	if actor.Role != domain.RoleAdmin {
		return domain.User{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.DeleteUser(ctx, id)
}

func (s *Service) ListUserTasks(ctx context.Context, userID int, filter PageFilter) (domain.Paginated[domain.Task], error) {
	return s.store.ListUserTasks(ctx, userID, NormalizePage(filter))
}

func validateCreateUser(input CreateUserInput) error {
	if len(strings.TrimSpace(input.Name)) < 2 {
		return domain.NewError(domain.ErrValidation, "Name must be at least 2 characters")
	}
	if !strings.Contains(input.Email, "@") {
		return domain.NewError(domain.ErrValidation, "Please enter a valid email address")
	}
	if len(input.Password) < 6 {
		return domain.NewError(domain.ErrValidation, "Password must be at least 6 characters")
	}
	return nil
}
