package service

import (
	"context"

	"project-manager-go/internal/domain"
)

func (s *Service) CreateTask(ctx context.Context, actor domain.AuthUser, input TaskInput) (domain.Task, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.Task{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	if input.Status == nil {
		status := domain.TaskStatusBacklog
		input.Status = &status
	}
	if input.Priority == nil {
		priority := domain.TaskPriorityMedium
		input.Priority = &priority
	}
	if err := validateTaskStatus(input.Status); err != nil {
		return domain.Task{}, err
	}
	if err := validateTaskPriority(input.Priority); err != nil {
		return domain.Task{}, err
	}
	if err := validatePositiveID(input.ProjectID, "Project ID"); err != nil {
		return domain.Task{}, err
	}
	if err := validatePositiveID(input.AssigneeID, "Assignee ID"); err != nil {
		return domain.Task{}, err
	}
	if _, err := s.store.GetProjectByID(ctx, actor, input.ProjectID); err != nil {
		return domain.Task{}, err
	}
	input.CreatorID = actor.ID
	return s.store.CreateTask(ctx, input)
}

func (s *Service) ListTasks(
	ctx context.Context,
	actor domain.AuthUser,
	filter TaskFilter,
) (domain.Paginated[domain.Task], error) {
	if err := validateOptionalPositiveID(filter.ProjectID, "Project ID"); err != nil {
		return domain.Paginated[domain.Task]{}, err
	}
	if err := validateTaskStatus(filter.Status); err != nil {
		return domain.Paginated[domain.Task]{}, err
	}
	if err := validateTaskPriority(filter.Priority); err != nil {
		return domain.Paginated[domain.Task]{}, err
	}
	if err := validateOptionalPositiveID(filter.AssigneeID, "Assignee ID"); err != nil {
		return domain.Paginated[domain.Task]{}, err
	}
	filter.PageFilter = NormalizePage(filter.PageFilter)
	return s.store.ListTasks(ctx, actor, filter)
}

func (s *Service) GetTaskByID(ctx context.Context, actor domain.AuthUser, id int) (domain.Task, error) {
	if err := validatePositiveID(id, "Task ID"); err != nil {
		return domain.Task{}, err
	}
	return s.store.GetTaskByID(ctx, actor, id)
}

func (s *Service) UpdateTask(ctx context.Context, actor domain.AuthUser, id int, input TaskPatchInput) (domain.Task, error) {
	if err := validatePositiveID(id, "Task ID"); err != nil {
		return domain.Task{}, err
	}
	if _, err := s.store.GetTaskByID(ctx, actor, id); err != nil {
		return domain.Task{}, err
	}
	if err := validateTaskStatus(input.Status); err != nil {
		return domain.Task{}, err
	}
	if err := validateTaskPriority(input.Priority); err != nil {
		return domain.Task{}, err
	}
	if input.ProjectID != nil {
		if err := validateOptionalPositiveID(input.ProjectID, "Project ID"); err != nil {
			return domain.Task{}, err
		}
		if _, err := s.store.GetProjectByID(ctx, actor, *input.ProjectID); err != nil {
			return domain.Task{}, err
		}
	}
	if input.AssigneeID != nil {
		if err := validateOptionalPositiveID(input.AssigneeID, "Assignee ID"); err != nil {
			return domain.Task{}, err
		}
	}
	return s.store.UpdateTask(ctx, actor, id, input)
}

func (s *Service) DeleteTask(ctx context.Context, actor domain.AuthUser, id int) (domain.Task, error) {
	if err := validatePositiveID(id, "Task ID"); err != nil {
		return domain.Task{}, err
	}
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.Task{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	return s.store.DeleteTask(ctx, actor, id)
}

func (s *Service) AssignUserToTask(
	ctx context.Context,
	actor domain.AuthUser,
	input TaskAssignmentInput,
) (domain.TaskAssignment, bool, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.TaskAssignment{}, false, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	if err := validatePositiveID(input.TaskID, "Task ID"); err != nil {
		return domain.TaskAssignment{}, false, err
	}
	if err := validatePositiveID(input.UserID, "User ID"); err != nil {
		return domain.TaskAssignment{}, false, err
	}
	if _, err := s.store.GetTaskByID(ctx, actor, input.TaskID); err != nil {
		return domain.TaskAssignment{}, false, err
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

func (s *Service) ListTaskAssignments(ctx context.Context, actor domain.AuthUser, taskID int) ([]domain.TaskAssignment, error) {
	if err := validatePositiveID(taskID, "Task ID"); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTaskByID(ctx, actor, taskID); err != nil {
		return nil, err
	}
	return s.store.ListTaskAssignments(ctx, taskID)
}

func (s *Service) RemoveTaskAssignment(ctx context.Context, actor domain.AuthUser, id int) (domain.TaskAssignment, error) {
	if err := validatePositiveID(id, "Task assignment ID"); err != nil {
		return domain.TaskAssignment{}, err
	}
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.TaskAssignment{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	assignment, err := s.store.GetTaskAssignmentByID(ctx, id)
	if err != nil {
		return domain.TaskAssignment{}, err
	}
	if _, err := s.store.GetTaskByID(ctx, actor, assignment.TaskID); err != nil {
		return domain.TaskAssignment{}, err
	}
	return s.store.RemoveTaskAssignment(ctx, id)
}
