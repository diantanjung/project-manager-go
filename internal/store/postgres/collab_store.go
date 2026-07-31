package postgres

import (
	"context"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) CreateComment(ctx context.Context, input service.CommentInput) (domain.Comment, error) {
	var comment domain.Comment
	err := s.db.QueryRow(ctx, `
		INSERT INTO comments (content, task_id, author_id)
		VALUES ($1, $2, $3)
		RETURNING id, content, task_id, author_id, created_at, updated_at
	`, input.Content, input.TaskID, input.AuthorID).Scan(
		&comment.ID, &comment.Content, &comment.TaskID, &comment.AuthorID, &comment.CreatedAt, &comment.UpdatedAt,
	)
	return comment, notFound(err)
}

func (s *Store) ListComments(ctx context.Context, taskID int) ([]domain.CommentView, error) {
	rows, err := s.db.Query(ctx, `
		SELECT c.id, c.content, c.task_id, c.author_id, u.name, NULL::text, c.created_at, c.updated_at
		FROM comments c
		LEFT JOIN users u ON u.id = c.author_id
		WHERE c.task_id = $1
		ORDER BY c.created_at DESC
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.CommentView{}
	for rows.Next() {
		var item domain.CommentView
		if err := rows.Scan(&item.ID, &item.Content, &item.TaskID, &item.AuthorID, &item.AuthorName, &item.AuthorAvatarURL, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) GetCommentByID(ctx context.Context, id int) (domain.Comment, error) {
	var comment domain.Comment
	err := s.db.QueryRow(ctx, `
		SELECT id, content, task_id, author_id, created_at, updated_at
		FROM comments WHERE id = $1
	`, id).Scan(&comment.ID, &comment.Content, &comment.TaskID, &comment.AuthorID, &comment.CreatedAt, &comment.UpdatedAt)
	return comment, notFound(err)
}

func (s *Store) UpdateComment(ctx context.Context, id int, content string) (domain.Comment, error) {
	var comment domain.Comment
	err := s.db.QueryRow(ctx, `
		UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2
		RETURNING id, content, task_id, author_id, created_at, updated_at
	`, content, id).Scan(&comment.ID, &comment.Content, &comment.TaskID, &comment.AuthorID, &comment.CreatedAt, &comment.UpdatedAt)
	return comment, notFound(err)
}

func (s *Store) DeleteComment(ctx context.Context, id int) (domain.Comment, error) {
	var comment domain.Comment
	err := s.db.QueryRow(ctx, `
		DELETE FROM comments WHERE id = $1
		RETURNING id, content, task_id, author_id, created_at, updated_at
	`, id).Scan(&comment.ID, &comment.Content, &comment.TaskID, &comment.AuthorID, &comment.CreatedAt, &comment.UpdatedAt)
	return comment, notFound(err)
}

func (s *Store) FindUserIDByExactName(ctx context.Context, name string) (int, bool, error) {
	var id int
	err := s.db.QueryRow(ctx, `SELECT id FROM users WHERE name = $1 LIMIT 1`, name).Scan(&id)
	ok, err := exists(err)
	return id, ok, err
}

func (s *Store) CreateAttachment(ctx context.Context, input service.AttachmentInput) (domain.Attachment, error) {
	originalName := input.OriginalName
	if originalName == "" {
		originalName = input.FileName
	}
	storageKey := input.StorageKey
	if storageKey == "" {
		storageKey = input.FileURL
	}

	var attachment domain.Attachment
	err := s.db.QueryRow(ctx, `
		INSERT INTO attachments (file_name, original_name, storage_key, file_size, mime_type, task_id, uploader_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, file_name, original_name, storage_key, file_size, mime_type, task_id, uploader_id, created_at
	`, input.FileName, originalName, storageKey, input.FileSize, input.MimeType, input.TaskID, input.UploaderID).Scan(
		&attachment.ID, &attachment.FileName, &attachment.OriginalName, &attachment.StorageKey, &attachment.FileSize,
		&attachment.MimeType, &attachment.TaskID, &attachment.UploaderID, &attachment.CreatedAt,
	)
	return attachment, notFound(err)
}

func (s *Store) ListAttachments(ctx context.Context, taskID int) ([]domain.Attachment, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, file_name, original_name, storage_key, file_size, mime_type, task_id, uploader_id, created_at
		FROM attachments WHERE task_id = $1
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Attachment{}
	for rows.Next() {
		var item domain.Attachment
		if err := rows.Scan(&item.ID, &item.FileName, &item.OriginalName, &item.StorageKey, &item.FileSize, &item.MimeType, &item.TaskID, &item.UploaderID, &item.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) GetAttachmentByID(ctx context.Context, id int) (domain.Attachment, error) {
	var item domain.Attachment
	err := s.db.QueryRow(ctx, `
		SELECT id, file_name, original_name, storage_key, file_size, mime_type, task_id, uploader_id, created_at
		FROM attachments WHERE id = $1
	`, id).Scan(&item.ID, &item.FileName, &item.OriginalName, &item.StorageKey, &item.FileSize, &item.MimeType, &item.TaskID, &item.UploaderID, &item.CreatedAt)
	return item, notFound(err)
}

func (s *Store) DeleteAttachment(ctx context.Context, id int) (domain.Attachment, error) {
	var item domain.Attachment
	err := s.db.QueryRow(ctx, `
		DELETE FROM attachments WHERE id = $1
		RETURNING id, file_name, original_name, storage_key, file_size, mime_type, task_id, uploader_id, created_at
	`, id).Scan(&item.ID, &item.FileName, &item.OriginalName, &item.StorageKey, &item.FileSize, &item.MimeType, &item.TaskID, &item.UploaderID, &item.CreatedAt)
	return item, notFound(err)
}

func (s *Store) CreateNotification(ctx context.Context, input service.NotificationInput) (domain.Notification, error) {
	var item domain.Notification
	err := s.db.QueryRow(ctx, `
		INSERT INTO notifications (user_id, actor_id, type, task_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, actor_id, type, task_id, is_read, created_at
	`, input.UserID, input.ActorID, input.Type, input.TaskID).Scan(&item.ID, &item.UserID, &item.ActorID, &item.Type, &item.TaskID, &item.IsRead, &item.CreatedAt)
	return item, notFound(err)
}

func (s *Store) ListNotifications(ctx context.Context, userID int) ([]domain.Notification, error) {
	rows, err := s.db.Query(ctx, `
		SELECT n.id, n.user_id, n.actor_id, u.name, NULL::text, n.type, n.task_id, n.is_read, n.created_at
		FROM notifications n
		LEFT JOIN users u ON u.id = n.actor_id
		WHERE n.user_id = $1
		ORDER BY n.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Notification{}
	for rows.Next() {
		var item domain.Notification
		if err := rows.Scan(&item.ID, &item.UserID, &item.ActorID, &item.ActorName, &item.ActorAvatarURL, &item.Type, &item.TaskID, &item.IsRead, &item.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) MarkNotificationRead(ctx context.Context, id int, userID int) (domain.Notification, error) {
	var item domain.Notification
	err := s.db.QueryRow(ctx, `
		UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, actor_id, type, task_id, is_read, created_at
	`, id, userID).Scan(&item.ID, &item.UserID, &item.ActorID, &item.Type, &item.TaskID, &item.IsRead, &item.CreatedAt)
	return item, notFound(err)
}

func (s *Store) MarkAllNotificationsRead(ctx context.Context, userID int) ([]domain.Notification, error) {
	rows, err := s.db.Query(ctx, `
		UPDATE notifications SET is_read = true WHERE user_id = $1
		RETURNING id, user_id, actor_id, type, task_id, is_read, created_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Notification{}
	for rows.Next() {
		var item domain.Notification
		if err := rows.Scan(&item.ID, &item.UserID, &item.ActorID, &item.Type, &item.TaskID, &item.IsRead, &item.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
