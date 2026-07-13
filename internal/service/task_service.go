package service

import (
	"context"
	"strings"

	"project-manager-go/internal/domain"
)

func (s *Service) CreateTask(ctx context.Context, actor domain.AuthUser, input TaskInput) (domain.Task, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.Task{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	if strings.TrimSpace(input.Title) == "" || len(input.Title) < 2 {
		return domain.Task{}, domain.NewError(domain.ErrValidation, "Task title must be at least 2 characters")
	}
	if input.Status == nil {
		status := domain.TaskStatusBacklog
		input.Status = &status
	}
	if input.Priority == nil {
		priority := domain.TaskPriorityMedium
		input.Priority = &priority
	}
	input.CreatorID = actor.ID
	return s.store.CreateTask(ctx, input)
}

func (s *Service) ListTasks(
	ctx context.Context,
	actor domain.AuthUser,
	filter TaskFilter,
) (domain.Paginated[domain.Task], error) {
	filter.PageFilter = NormalizePage(filter.PageFilter)
	return s.store.ListTasks(ctx, actor, filter)
}

func (s *Service) GetTaskByID(ctx context.Context, actor domain.AuthUser, id int) (domain.Task, error) {
	return s.store.GetTaskByID(ctx, actor, id)
}

func (s *Service) UpdateTask(ctx context.Context, _ domain.AuthUser, id int, input TaskPatchInput) (domain.Task, error) {
	return s.store.UpdateTask(ctx, id, input)
}

func (s *Service) DeleteTask(ctx context.Context, actor domain.AuthUser, id int) (domain.Task, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.Task{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.DeleteTask(ctx, id)
}

func (s *Service) AssignUserToTask(
	ctx context.Context,
	actor domain.AuthUser,
	input TaskAssignmentInput,
) (domain.TaskAssignment, bool, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.TaskAssignment{}, false, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	assignment, exists, err := s.store.AssignUserToTask(ctx, input)
	if err != nil || exists {
		return assignment, exists, err
	}
	taskID := input.TaskID
	_, err = s.store.CreateNotification(ctx, NotificationInput{
		UserID: input.UserID,
		Type:   domain.NotificationTaskAssigned,
		TaskID: &taskID,
	})
	return assignment, false, err
}

func (s *Service) ListTaskAssignments(ctx context.Context, taskID int) ([]domain.TaskAssignment, error) {
	return s.store.ListTaskAssignments(ctx, taskID)
}

func (s *Service) RemoveTaskAssignment(ctx context.Context, actor domain.AuthUser, id int) (domain.TaskAssignment, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.TaskAssignment{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.RemoveTaskAssignment(ctx, id)
}
