package postgres

import (
	"context"
	"fmt"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) CreateTask(ctx context.Context, input service.TaskInput) (domain.Task, error) {
	row := s.db.QueryRow(ctx, `
		INSERT INTO tasks (title, description, status, priority, project_id, creator_id, assignee_id, due_date, position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, title, description, status, priority, project_id, NULL::text,
			creator_id, assignee_id, NULL::text, NULL::text, due_date::text, position, created_at, updated_at
	`, input.Title, input.Description, value(input.Status), value(input.Priority), input.ProjectID, input.CreatorID, input.AssigneeID, input.DueDate, input.Position)
	return scanTask(row)
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
	if err := s.db.QueryRow(ctx, totalQuery, where.values...).Scan(&total); err != nil {
		return domain.Paginated[domain.Task]{}, err
	}

	sortBy := allow(filter.SortBy, "ta.created_at", map[string]string{
		"title": "ta.title", "createdAt": "ta.created_at", "updatedAt": "ta.updated_at",
		"dueDate": "ta.due_date", "priority": "ta.priority",
	})
	args := append(where.values, filter.Limit, offset(filter.PageFilter))
	query := fmt.Sprintf(`
		SELECT ta.id, ta.title, ta.description, ta.status, ta.priority, ta.project_id, p.name,
			ta.creator_id, ta.assignee_id, u.name, NULL::text, ta.due_date::text, ta.position,
			ta.created_at, ta.updated_at
		FROM tasks ta
		LEFT JOIN projects p ON p.id = ta.project_id
		LEFT JOIN users u ON u.id = ta.assignee_id
		%s
		ORDER BY %s %s LIMIT $%d OFFSET $%d
	`, where.clause(), sortBy, order(filter.Order), len(args)-1, len(args))
	rows, err := s.db.Query(ctx, query, args...)
	return s.scanTaskPage(rows, filter.PageFilter, total, err)
}

func (s *Store) GetTaskByID(ctx context.Context, user domain.AuthUser, id int) (domain.Task, error) {
	where := taskWhere(user)
	where.add("ta.id = $%d", id)
	query := `
		SELECT ta.id, ta.title, ta.description, ta.status, ta.priority, ta.project_id, p.name,
			ta.creator_id, ta.assignee_id, u.name, NULL::text, ta.due_date::text, ta.position,
			ta.created_at, ta.updated_at
		FROM tasks ta
		LEFT JOIN projects p ON p.id = ta.project_id
		LEFT JOIN users u ON u.id = ta.assignee_id
		` + where.clause()
	task, err := scanTask(s.db.QueryRow(ctx, query, where.values...))
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

func (s *Store) UpdateTask(ctx context.Context, id int, input service.TaskPatchInput) (domain.Task, error) {
	sets := updateBuilder{}
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
		UPDATE tasks SET %s WHERE id = $%d
		RETURNING id, title, description, status, priority, project_id, NULL::text,
			creator_id, assignee_id, NULL::text, NULL::text, due_date::text, position, created_at, updated_at
	`, sets.clause(), len(sets.values)+1)
	sets.values = append(sets.values, id)
	return scanTask(s.db.QueryRow(ctx, query, sets.values...))
}

func (s *Store) DeleteTask(ctx context.Context, id int) (domain.Task, error) {
	return scanTask(s.db.QueryRow(ctx, `
		DELETE FROM tasks WHERE id = $1
		RETURNING id, title, description, status, priority, project_id, NULL::text,
			creator_id, assignee_id, NULL::text, NULL::text, due_date::text, position, created_at, updated_at
	`, id))
}

func (s *Store) AssignUserToTask(ctx context.Context, input service.TaskAssignmentInput) (domain.TaskAssignment, bool, error) {
	var id int
	err := s.db.QueryRow(ctx, `
		SELECT id FROM task_assignments WHERE task_id = $1 AND user_id = $2
	`, input.TaskID, input.UserID).Scan(&id)
	ok, err := exists(err)
	if err != nil {
		return domain.TaskAssignment{}, false, err
	}
	if ok {
		item, err := s.GetTaskAssignmentByID(ctx, id)
		return item, true, err
	}

	var item domain.TaskAssignment
	err = s.db.QueryRow(ctx, `
		INSERT INTO task_assignments (task_id, user_id)
		VALUES ($1, $2)
		RETURNING id, task_id, user_id, assigned_at
	`, input.TaskID, input.UserID).Scan(&item.ID, &item.TaskID, &item.UserID, &item.AssignedAt)
	return item, false, notFound(err)
}

func (s *Store) ListTaskAssignments(ctx context.Context, taskID int) ([]domain.TaskAssignment, error) {
	rows, err := s.db.Query(ctx, `
		SELECT ta.id, ta.task_id, ta.user_id, u.name, u.email, NULL::text, ta.assigned_at
		FROM task_assignments ta
		LEFT JOIN users u ON u.id = ta.user_id
		WHERE ta.task_id = $1
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.TaskAssignment{}
	for rows.Next() {
		var item domain.TaskAssignment
		if err := rows.Scan(&item.ID, &item.TaskID, &item.UserID, &item.UserName, &item.UserEmail, &item.UserAvatarURL, &item.AssignedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) GetTaskAssignmentByID(ctx context.Context, id int) (domain.TaskAssignment, error) {
	var item domain.TaskAssignment
	err := s.db.QueryRow(ctx, `
		SELECT id, task_id, user_id, assigned_at FROM task_assignments WHERE id = $1
	`, id).Scan(&item.ID, &item.TaskID, &item.UserID, &item.AssignedAt)
	return item, notFound(err)
}

func (s *Store) RemoveTaskAssignment(ctx context.Context, id int) (domain.TaskAssignment, error) {
	var item domain.TaskAssignment
	err := s.db.QueryRow(ctx, `
		DELETE FROM task_assignments WHERE id = $1
		RETURNING id, task_id, user_id, assigned_at
	`, id).Scan(&item.ID, &item.TaskID, &item.UserID, &item.AssignedAt)
	return item, notFound(err)
}
