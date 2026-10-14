package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) CreateTask(ctx context.Context, input service.TaskInput) (domain.Task, error) {
	return getOne[domain.Task](s, ctx, `
		WITH created AS (
			INSERT INTO tasks (title, description, status, priority, project_id, creator_id, assignee_id, due_date, position, estimate_minutes, completed_at)
			VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, 0), $8, $9, $10, CASE WHEN $3 = 'done' THEN NOW() ELSE NULL END)
			RETURNING *
		), logged AS (
			INSERT INTO activity_logs (actor_id, entity_type, entity_id, action, after)
			SELECT creator_id, 'task', id, 'created', to_jsonb(created) FROM created RETURNING id
		)
		SELECT created.id, created.title, created.description, created.status, created.priority, created.project_id, NULL::text AS project_name,
			created.creator_id, created.assignee_id, NULL::text AS assignee_name, NULL::text AS assignee_avatar_url,
			created.due_date::text AS due_date, created.position, created.estimate_minutes, created.completed_at, created.created_at, created.updated_at
		FROM created CROSS JOIN logged
	`, input.Title, input.Description, value(input.Status), value(input.Priority), input.ProjectID, input.CreatorID, input.AssigneeID, input.DueDate, input.Position, input.EstimateMinutes)
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
	if filter.DueBefore != nil {
		where.add("ta.due_date <= $%d", *filter.DueBefore)
	}
	if filter.DueAfter != nil {
		where.add("ta.due_date >= $%d", *filter.DueAfter)
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
	args := append([]any{}, where.values...)
	limitClause := ""
	if !filter.Board {
		args = append(args, filter.Limit, offset(filter.PageFilter))
		limitClause = fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}
	query := fmt.Sprintf(`
		SELECT ta.id, ta.title, ta.description, ta.status, ta.priority, ta.project_id, p.name AS project_name,
			ta.creator_id, ta.assignee_id, u.name AS assignee_name, NULL::text AS assignee_avatar_url, ta.due_date::text AS due_date, ta.position,
			ta.estimate_minutes, ta.completed_at, ta.created_at, ta.updated_at
		FROM tasks ta
		LEFT JOIN projects p ON p.id = ta.project_id
		LEFT JOIN users u ON u.id = ta.assignee_id
		%s
		ORDER BY %s %s NULLS LAST, ta.id ASC%s
	`, where.clause(), sortBy, order(filter.Order), limitClause)
	return selectPage[domain.Task](s, ctx, query, filter.PageFilter, total, args...)
}

func (s *Store) GetTaskByID(ctx context.Context, user domain.AuthUser, id int) (domain.Task, error) {
	where := taskWhere(user)
	where.add("ta.id = $%d", id)
	query := `
		SELECT ta.id, ta.title, ta.description, ta.status, ta.priority, ta.project_id, p.name AS project_name,
			ta.creator_id, ta.assignee_id, u.name AS assignee_name, NULL::text AS assignee_avatar_url, ta.due_date::text AS due_date, ta.position,
			ta.estimate_minutes, ta.completed_at, ta.created_at, ta.updated_at
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
	if input.ClearAssignee {
		sets.parts = append(sets.parts, "assignee_id = NULL")
	} else {
		sets.add("assignee_id", input.AssigneeID)
	}
	sets.add("due_date", input.DueDate)
	sets.add("position", input.Position)
	sets.add("estimate_minutes", input.EstimateMinutes)
	if input.Status != nil {
		if *input.Status == domain.TaskStatusDone {
			sets.parts = append(sets.parts, "completed_at = COALESCE(ta.completed_at, NOW())")
		} else {
			sets.parts = append(sets.parts, "completed_at = NULL")
		}
	}
	sets.touchUpdatedAt()
	actorArg := len(sets.values) + 1
	query := fmt.Sprintf(`
		WITH previous AS MATERIALIZED (
			SELECT ta.* FROM tasks ta %s
		), changed AS (
			UPDATE tasks ta SET %s FROM previous
			WHERE ta.id = previous.id
			RETURNING ta.*
		), logged AS (
			INSERT INTO activity_logs (actor_id, entity_type, entity_id, action, before, after)
			SELECT $%d, 'task', changed.id, 'updated', to_jsonb(previous), to_jsonb(changed)
			FROM changed JOIN previous ON previous.id = changed.id RETURNING id
		)
		SELECT changed.id, changed.title, changed.description, changed.status, changed.priority, changed.project_id, NULL::text AS project_name,
			changed.creator_id, changed.assignee_id, NULL::text AS assignee_name, NULL::text AS assignee_avatar_url,
			changed.due_date::text AS due_date, changed.position, changed.estimate_minutes, changed.completed_at, changed.created_at, changed.updated_at
		FROM changed CROSS JOIN logged
	`, where.clause(), sets.clause(), actorArg)
	return getOne[domain.Task](s, ctx, query, append(sets.values, user.ID)...)
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
	query := fmt.Sprintf(`
		WITH removed AS (
			DELETE FROM tasks ta %s RETURNING ta.*
		), logged AS (
			INSERT INTO activity_logs (actor_id, entity_type, entity_id, action, before)
			SELECT $%d, 'task', id, 'deleted', to_jsonb(removed) FROM removed RETURNING id
		)
		SELECT removed.id, removed.title, removed.description, removed.status, removed.priority, removed.project_id, NULL::text AS project_name,
			removed.creator_id, removed.assignee_id, NULL::text AS assignee_name, NULL::text AS assignee_avatar_url,
			removed.due_date::text AS due_date, removed.position, removed.estimate_minutes, removed.completed_at, removed.created_at, removed.updated_at
		FROM removed CROSS JOIN logged
	`, where.clause(), len(where.values)+1)
	return getOne[domain.Task](s, ctx, query, append(where.values, user.ID)...)
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

func (s *Store) RemoveTaskAssignmentByUser(ctx context.Context, taskID, userID int) (domain.TaskAssignment, error) {
	return getOne[domain.TaskAssignment](s, ctx, `
		DELETE FROM task_assignments WHERE task_id = $1 AND user_id = $2
		RETURNING id, task_id, user_id, assigned_at
	`, taskID, userID)
}

func (s *Store) ReorderTasks(ctx context.Context, user domain.AuthUser, input service.TaskReorderInput) ([]domain.Task, error) {
	switch db := s.db.(type) {
	case *sqlx.DB:
		tx, err := db.BeginTxx(ctx, nil)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()
		tasks, err := (&Store{db: tx}).reorderTasksInTransaction(ctx, user, input)
		if err != nil {
			return nil, err
		}
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return tasks, nil
	case *sqlx.Tx:
		if _, err := db.ExecContext(ctx, "SAVEPOINT task_batch_reorder"); err != nil {
			return nil, err
		}
		tasks, err := s.reorderTasksInTransaction(ctx, user, input)
		if err != nil {
			if _, rollbackErr := db.ExecContext(ctx, "ROLLBACK TO SAVEPOINT task_batch_reorder"); rollbackErr != nil {
				return nil, fmt.Errorf("reorder failed: %w; rollback failed: %v", err, rollbackErr)
			}
			return nil, err
		}
		if _, err := db.ExecContext(ctx, "RELEASE SAVEPOINT task_batch_reorder"); err != nil {
			return nil, err
		}
		return tasks, nil
	default:
		return nil, fmt.Errorf("task reorder requires a SQL transaction")
	}
}

func (s *Store) reorderTasksInTransaction(ctx context.Context, user domain.AuthUser, input service.TaskReorderInput) ([]domain.Task, error) {
	var projectID int
	if err := s.db.GetContext(ctx, &projectID, "SELECT id FROM projects WHERE id = $1 FOR UPDATE", input.ProjectID); err != nil {
		return nil, notFound(err)
	}
	for _, move := range input.Moves {
		task, err := s.GetTaskByID(ctx, user, move.TaskID)
		if err != nil {
			return nil, err
		}
		if task.ProjectID != input.ProjectID {
			return nil, domain.NewError(domain.ErrValidation, "Task does not belong to the requested project")
		}
	}
	for _, move := range input.Moves {
		status, position := move.Status, move.Position
		if _, err := s.reorderTaskInTransaction(ctx, user, move.TaskID, service.TaskPatchInput{
			Status: &status, Position: &position,
		}); err != nil {
			return nil, err
		}
	}
	tasks := make([]domain.Task, 0, len(input.Moves))
	for _, move := range input.Moves {
		task, err := s.GetTaskByID(ctx, user, move.TaskID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *Store) ListTaskActivity(ctx context.Context, user domain.AuthUser, id int) ([]domain.ActivityLog, error) {
	where := taskWhere(user)
	where.add("ta.id = $%d", id)
	var rows []struct {
		ID         int        `db:"id"`
		ActorID    *int       `db:"actor_id"`
		EntityType string     `db:"entity_type"`
		EntityID   int        `db:"entity_id"`
		Action     string     `db:"action"`
		Before     []byte     `db:"before"`
		After      []byte     `db:"after"`
		IPAddress  *string    `db:"ip_address"`
		UserAgent  *string    `db:"user_agent"`
		CreatedAt  *time.Time `db:"created_at"`
	}
	query := `SELECT l.id, l.actor_id, l.entity_type, l.entity_id, l.action,
		l.before, l.after, l.ip_address::text AS ip_address, l.user_agent, l.created_at
		FROM activity_logs l WHERE l.entity_type = 'task' AND l.entity_id IN
		(SELECT ta.id FROM tasks ta ` + where.clause() + `)
		ORDER BY l.created_at DESC, l.id DESC`
	if err := s.db.SelectContext(ctx, &rows, query, where.values...); err != nil {
		return nil, err
	}
	activity := make([]domain.ActivityLog, 0, len(rows))
	for _, row := range rows {
		item := domain.ActivityLog{
			ID: row.ID, ActorID: row.ActorID, EntityType: row.EntityType,
			EntityID: row.EntityID, Action: row.Action, IPAddress: row.IPAddress,
			UserAgent: row.UserAgent, CreatedAt: row.CreatedAt,
		}
		if len(row.Before) > 0 {
			if err := json.Unmarshal(row.Before, &item.Before); err != nil {
				return nil, err
			}
		}
		if len(row.After) > 0 {
			if err := json.Unmarshal(row.After, &item.After); err != nil {
				return nil, err
			}
		}
		activity = append(activity, item)
	}
	return activity, nil
}
