package postgres

import (
	"context"
	"fmt"

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
		"dueDate": "ta.due_date", "priority": "ta.priority",
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
		ORDER BY %s %s LIMIT $%d OFFSET $%d
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
