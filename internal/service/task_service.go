package service

import (
	"context"
	"time"

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
	if input.EstimateMinutes != nil && *input.EstimateMinutes < 0 {
		return domain.Task{}, validationError("Estimate minutes must be non-negative")
	}
	if err := validatePositiveID(input.ProjectID, "Project ID"); err != nil {
		return domain.Task{}, err
	}
	if input.AssigneeID < 0 {
		return domain.Task{}, validationError("Assignee ID must be positive")
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
	for _, date := range []*string{filter.DueBefore, filter.DueAfter} {
		if date != nil {
			if _, err := time.Parse("2006-01-02", *date); err != nil {
				return domain.Paginated[domain.Task]{}, validationError("Due date filter must use YYYY-MM-DD")
			}
		}
	}
	if filter.DueBefore != nil && filter.DueAfter != nil && *filter.DueAfter > *filter.DueBefore {
		return domain.Paginated[domain.Task]{}, validationError("dueAfter must not be later than dueBefore")
	}
	if !filter.Board {
		filter.PageFilter = NormalizePage(filter.PageFilter)
	}
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
	task, err := s.store.GetTaskByID(ctx, actor, id)
	if err != nil {
		return domain.Task{}, err
	}
	if input.ProjectID != nil {
		return domain.Task{}, validationError("Project ID cannot be changed through task edit")
	}
	if input.ClearAssignee && input.AssigneeID != nil {
		return domain.Task{}, validationError("Assignee cannot be set and cleared in the same request")
	}
	if input.Title == nil && input.Description == nil && input.Status == nil && input.Priority == nil && input.AssigneeID == nil && !input.ClearAssignee && input.DueDate == nil && input.Position == nil && input.EstimateMinutes == nil {
		return domain.Task{}, validationError("At least one task field is required")
	}
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		if input.Status == nil || input.Title != nil || input.Description != nil || input.Priority != nil || input.AssigneeID != nil || input.ClearAssignee || input.DueDate != nil || input.Position != nil || input.EstimateMinutes != nil {
			return domain.Task{}, domain.NewError(domain.ErrForbidden, "Access denied. Only task status can be updated.")
		}
		if task.AssigneeID == nil || *task.AssigneeID != actor.ID {
			return domain.Task{}, domain.NewError(domain.ErrForbidden, "Access denied. Task is not assigned to this user.")
		}
	}
	if err := validateTaskStatus(input.Status); err != nil {
		return domain.Task{}, err
	}
	if err := validateTaskPriority(input.Priority); err != nil {
		return domain.Task{}, err
	}
	if input.EstimateMinutes != nil && *input.EstimateMinutes < 0 {
		return domain.Task{}, validationError("Estimate minutes must be non-negative")
	}
	if input.AssigneeID != nil {
		if err := validateOptionalPositiveID(input.AssigneeID, "Assignee ID"); err != nil {
			return domain.Task{}, err
		}
	}
	return s.store.UpdateTask(ctx, actor, id, input)
}

func (s *Service) ReorderTasks(ctx context.Context, actor domain.AuthUser, input TaskReorderInput) ([]domain.Task, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return nil, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	if err := validatePositiveID(input.ProjectID, "Project ID"); err != nil {
		return nil, err
	}
	if len(input.Moves) == 0 {
		return nil, validationError("At least one task move is required")
	}
	if _, err := s.store.GetProjectByID(ctx, actor, input.ProjectID); err != nil {
		return nil, err
	}
	seen := make(map[int]bool, len(input.Moves))
	for _, move := range input.Moves {
		if err := validatePositiveID(move.TaskID, "Task ID"); err != nil {
			return nil, err
		}
		if move.Position < 0 {
			return nil, validationError("Task position must be non-negative")
		}
		if err := validateTaskStatus(&move.Status); err != nil {
			return nil, err
		}
		if seen[move.TaskID] {
			return nil, validationError("Duplicate task in moves")
		}
		seen[move.TaskID] = true
	}
	return s.store.ReorderTasks(ctx, actor, input)
}

func (s *Service) ListTaskActivity(ctx context.Context, actor domain.AuthUser, id int) ([]domain.ActivityLog, error) {
	if err := validatePositiveID(id, "Task ID"); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTaskByID(ctx, actor, id); err != nil {
		return nil, err
	}
	return s.store.ListTaskActivity(ctx, actor, id)
}

func (s *Service) RemoveTaskAssignmentByUser(ctx context.Context, actor domain.AuthUser, taskID, userID int) (domain.TaskAssignment, error) {
	if !HasRole(actor.Role, domain.RoleProjectManager) {
		return domain.TaskAssignment{}, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions.")
	}
	if err := validatePositiveID(taskID, "Task ID"); err != nil {
		return domain.TaskAssignment{}, err
	}
	if err := validatePositiveID(userID, "User ID"); err != nil {
		return domain.TaskAssignment{}, err
	}
	if _, err := s.store.GetTaskByID(ctx, actor, taskID); err != nil {
		return domain.TaskAssignment{}, err
	}
	return s.store.RemoveTaskAssignmentByUser(ctx, taskID, userID)
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
