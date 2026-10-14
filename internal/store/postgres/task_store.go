package postgres

import (
	"context"
	"fmt"
	"slices"

	"github.com/jmoiron/sqlx"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) CreateTask(ctx context.Context, input service.TaskInput) (domain.Task, error) {
	return getOne[domain.Task](s, ctx, `
		INSERT INTO tasks (title, description, status, priority, project_id, creator_id, assignee_id, due_date, position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, title, description, status, priority, project_id, NULL::text AS project_name,
			creator_id, assignee_id, NULL::text AS assignee_name, NULL::text AS assignee_avatar_url, due_date::text AS due_date, position, created_at, updated_at
	`, input.Title, input.Description, value(input.Status), value(input.Priority), input.ProjectID, input.CreatorID, input.AssigneeID, input.DueDate, input.Position)
}

func (s *Store) ListTasks(ctx context.Context, user domain.AuthUser, filter service.TaskFilter) (domain.Paginated[domain.Task], error) {
	where := taskWhere(user)
	if filter.Search != "" {
		where.add("(ta.title ILIKE $%d OR ta.description ILIKE $%d)", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.ProjectID != nil {
		where.add("ta.project_id = $%d", *filter.ProjectID)
	}
	if filter.Status != nil {
		where.add("ta.status = $%d", *filter.Status)
	}
	if filter.Priority != nil {
		where.add("ta.priority = $%d", *filter.Priority)
	}
	if filter.AssigneeID != nil {
		where.add("ta.assignee_id = $%d", *filter.AssigneeID)
	}

	totalQuery := "SELECT count(*) FROM tasks ta " + where.clause()
	var total int
	if err := s.db.GetContext(ctx, &total, totalQuery, where.values...); err != nil {
		return domain.Paginated[domain.Task]{}, err
	}

	sortBy := allow(filter.SortBy, "ta.created_at", map[string]string{
		"title": "ta.title", "createdAt": "ta.created_at", "updatedAt": "ta.updated_at",
		"dueDate": "ta.due_date", "priority": "ta.priority", "position": "ta.position",
	})
	args := append(where.values, filter.Limit, offset(filter.PageFilter))
	query := fmt.Sprintf(`
		SELECT ta.id, ta.title, ta.description, ta.status, ta.priority, ta.project_id, p.name AS project_name,
			ta.creator_id, ta.assignee_id, u.name AS assignee_name, NULL::text AS assignee_avatar_url, ta.due_date::text AS due_date, ta.position,
			ta.created_at, ta.updated_at
		FROM tasks ta
		LEFT JOIN projects p ON p.id = ta.project_id
		LEFT JOIN users u ON u.id = ta.assignee_id
		%s
		ORDER BY %s %s NULLS LAST, ta.id ASC LIMIT $%d OFFSET $%d
	`, where.clause(), sortBy, order(filter.Order), len(args)-1, len(args))
	return selectPage[domain.Task](s, ctx, query, filter.PageFilter, total, args...)
}

func (s *Store) GetTaskByID(ctx context.Context, user domain.AuthUser, id int) (domain.Task, error) {
	where := taskWhere(user)
	where.add("ta.id = $%d", id)
	query := `
		SELECT ta.id, ta.title, ta.description, ta.status, ta.priority, ta.project_id, p.name AS project_name,
			ta.creator_id, ta.assignee_id, u.name AS assignee_name, NULL::text AS assignee_avatar_url, ta.due_date::text AS due_date, ta.position,
			ta.created_at, ta.updated_at
		FROM tasks ta
		LEFT JOIN projects p ON p.id = ta.project_id
		LEFT JOIN users u ON u.id = ta.assignee_id
		` + where.clause()
	task, err := getOne[domain.Task](s, ctx, query, where.values...)
	if err != nil {
		return domain.Task{}, err
	}

	comments, err := s.ListComments(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}
	attachments, err := s.ListAttachments(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}
	task.Comments = comments
	task.Attachments = attachments
	return task, nil
}

func (s *Store) UpdateTask(ctx context.Context, user domain.AuthUser, id int, input service.TaskPatchInput) (domain.Task, error) {
	if input.Position != nil || input.Status != nil {
		return s.reorderTask(ctx, user, id, input)
	}
	return s.updateTaskFields(ctx, user, id, input)
}

func (s *Store) updateTaskFields(ctx context.Context, user domain.AuthUser, id int, input service.TaskPatchInput) (domain.Task, error) {
	where := taskWhere(user)
	where.add("ta.id = $%d", id)
	sets := updateBuilder{values: append([]any{}, where.values...)}
	sets.add("title", input.Title)
	sets.add("description", input.Description)
	sets.add("status", input.Status)
	sets.add("priority", input.Priority)
	sets.add("project_id", input.ProjectID)
	sets.add("assignee_id", input.AssigneeID)
	sets.add("due_date", input.DueDate)
	sets.add("position", input.Position)
	sets.touchUpdatedAt()
	query := fmt.Sprintf(`
		UPDATE tasks ta SET %s FROM (
			SELECT ta.id FROM tasks ta %s
		) allowed
		WHERE ta.id = allowed.id
		RETURNING ta.id, ta.title, ta.description, ta.status, ta.priority, ta.project_id, NULL::text AS project_name,
			ta.creator_id, ta.assignee_id, NULL::text AS assignee_name, NULL::text AS assignee_avatar_url, ta.due_date::text AS due_date, ta.position, ta.created_at, ta.updated_at
	`, sets.clause(), where.clause())
	return getOne[domain.Task](s, ctx, query, sets.values...)
}

func (s *Store) reorderTask(ctx context.Context, user domain.AuthUser, id int, input service.TaskPatchInput) (domain.Task, error) {
	switch db := s.db.(type) {
	case *sqlx.DB:
		tx, err := db.BeginTxx(ctx, nil)
		if err != nil {
			return domain.Task{}, err
		}
		defer tx.Rollback()
		task, err := (&Store{db: tx}).reorderTaskInTransaction(ctx, user, id, input)
		if err != nil {
			return domain.Task{}, err
		}
		if err := tx.Commit(); err != nil {
			return domain.Task{}, err
		}
		return task, nil
	case *sqlx.Tx:
		// Integration tests already run inside a transaction.
		if _, err := db.ExecContext(ctx, "SAVEPOINT task_reorder"); err != nil {
			return domain.Task{}, err
		}
		task, err := s.reorderTaskInTransaction(ctx, user, id, input)
		if err != nil {
			_, _ = db.ExecContext(ctx, "ROLLBACK TO SAVEPOINT task_reorder")
			return domain.Task{}, err
		}
		_, err = db.ExecContext(ctx, "RELEASE SAVEPOINT task_reorder")
		return task, err
	default:
		return domain.Task{}, fmt.Errorf("task reorder requires a SQL transaction")
	}
}

func (s *Store) reorderTaskInTransaction(ctx context.Context, user domain.AuthUser, id int, input service.TaskPatchInput) (domain.Task, error) {
	current, err := s.GetTaskByID(ctx, user, id)
	if err != nil {
		return domain.Task{}, err
	}
	// Project locks serialize concurrent moves within the same project.
	var lockedID int
	if err := s.db.GetContext(ctx, &lockedID, "SELECT id FROM projects WHERE id = $1 FOR UPDATE", current.ProjectID); err != nil {
		return domain.Task{}, notFound(err)
	}
	current, err = s.GetTaskByID(ctx, user, id)
	if err != nil {
		return domain.Task{}, err
	}
	targetStatus := current.Status
	if input.Status != nil {
		targetStatus = *input.Status
	}
	if input.Position == nil && targetStatus == current.Status {
		return s.updateTaskFields(ctx, user, id, input)
	}
	source, err := s.taskIDsForOrder(ctx, current.ProjectID, current.Status)
	if err != nil {
		return domain.Task{}, err
	}
	source = removeTaskID(source, id)
	if current.Status == targetStatus {
		return s.saveTaskOrder(ctx, user, id, input, source)
	}
	if err := s.setTaskPositions(ctx, source); err != nil {
		return domain.Task{}, err
	}
	destination, err := s.taskIDsForOrder(ctx, current.ProjectID, targetStatus)
	if err != nil {
		return domain.Task{}, err
	}
	return s.saveTaskOrder(ctx, user, id, input, removeTaskID(destination, id))
}

func (s *Store) taskIDsForOrder(ctx context.Context, projectID int, status domain.TaskStatus) ([]int, error) {
	ids := []int{}
	err := s.db.SelectContext(ctx, &ids, `
		SELECT id FROM tasks WHERE project_id = $1 AND status = $2
		ORDER BY position ASC NULLS LAST, id ASC
	`, projectID, status)
	return ids, err
}

func removeTaskID(ids []int, id int) []int {
	return slices.DeleteFunc(ids, func(candidate int) bool { return candidate == id })
}

func (s *Store) saveTaskOrder(ctx context.Context, user domain.AuthUser, id int, input service.TaskPatchInput, ids []int) (domain.Task, error) {
	position := len(ids)
	if input.Position != nil {
		position = min(*input.Position, len(ids))
	}
	ids = slices.Insert(ids, position, id)
	if err := s.setTaskPositions(ctx, ids); err != nil {
		return domain.Task{}, err
	}
	input.Position = &position
	return s.updateTaskFields(ctx, user, id, input)
}

func (s *Store) setTaskPositions(ctx context.Context, ids []int) error {
	for position, id := range ids {
		if _, err := s.db.ExecContext(ctx, "UPDATE tasks SET position = $1 WHERE id = $2 AND position IS DISTINCT FROM $1", position, id); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) DeleteTask(ctx context.Context, user domain.AuthUser, id int) (domain.Task, error) {
	where := taskWhere(user)
	where.add("ta.id = $%d", id)
	query := `
		DELETE FROM tasks ta ` + where.clause() + `
		RETURNING ta.id, ta.title, ta.description, ta.status, ta.priority, ta.project_id, NULL::text AS project_name,
			ta.creator_id, ta.assignee_id, NULL::text AS assignee_name, NULL::text AS assignee_avatar_url, ta.due_date::text AS due_date, ta.position, ta.created_at, ta.updated_at
	`
	return getOne[domain.Task](s, ctx, query, where.values...)
}

func (s *Store) AssignUserToTask(ctx context.Context, input service.TaskAssignmentInput) (domain.TaskAssignment, bool, error) {
	var id int
	err := s.db.GetContext(ctx, &id, `
		SELECT id FROM task_assignments WHERE task_id = $1 AND user_id = $2
	`, input.TaskID, input.UserID)
	ok, err := exists(err)
	if err != nil {
		return domain.TaskAssignment{}, false, err
	}
	if ok {
		item, err := s.GetTaskAssignmentByID(ctx, id)
		return item, true, err
	}

	item, err := getOne[domain.TaskAssignment](s, ctx, `
		INSERT INTO task_assignments (task_id, user_id)
		VALUES ($1, $2)
		RETURNING id, task_id, user_id, assigned_at
	`, input.TaskID, input.UserID)
	return item, false, err
}

func (s *Store) ListTaskAssignments(ctx context.Context, taskID int) ([]domain.TaskAssignment, error) {
	return selectAll[domain.TaskAssignment](s, ctx, `
		SELECT ta.id, ta.task_id, ta.user_id, u.name AS user_name, u.email AS user_email, NULL::text AS user_avatar_url, ta.assigned_at
		FROM task_assignments ta
		LEFT JOIN users u ON u.id = ta.user_id
		WHERE ta.task_id = $1
	`, taskID)
}

func (s *Store) GetTaskAssignmentByID(ctx context.Context, id int) (domain.TaskAssignment, error) {
	return getOne[domain.TaskAssignment](s, ctx, `
		SELECT id, task_id, user_id, assigned_at FROM task_assignments WHERE id = $1
	`, id)
}

func (s *Store) RemoveTaskAssignment(ctx context.Context, id int) (domain.TaskAssignment, error) {
	return getOne[domain.TaskAssignment](s, ctx, `
		DELETE FROM task_assignments WHERE id = $1
		RETURNING id, task_id, user_id, assigned_at
	`, id)
}
