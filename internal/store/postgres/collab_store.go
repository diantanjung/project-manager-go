package postgres

import (
	"context"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (s *Store) CreateComment(ctx context.Context, input service.CommentInput) (domain.Comment, error) {
	return getOne[domain.Comment](s, ctx, `
		INSERT INTO comments (content, task_id, author_id)
		VALUES ($1, $2, $3)
		RETURNING id, content, task_id, author_id, created_at, updated_at
	`, input.Content, input.TaskID, input.AuthorID)
}

func (s *Store) ListComments(ctx context.Context, taskID int) ([]domain.CommentView, error) {
	return selectAll[domain.CommentView](s, ctx, `
		SELECT c.id, c.content, c.task_id, c.author_id, u.name AS author_name, NULL::text AS author_avatar_url, c.created_at, c.updated_at
		FROM comments c
		LEFT JOIN users u ON u.id = c.author_id
		WHERE c.task_id = $1
		ORDER BY c.created_at DESC
	`, taskID)
}

func (s *Store) GetCommentByID(ctx context.Context, id int) (domain.Comment, error) {
	return getOne[domain.Comment](s, ctx, `
		SELECT id, content, task_id, author_id, created_at, updated_at
		FROM comments WHERE id = $1
	`, id)
}

func (s *Store) UpdateComment(ctx context.Context, id int, content string) (domain.Comment, error) {
	return getOne[domain.Comment](s, ctx, `
		UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2
		RETURNING id, content, task_id, author_id, created_at, updated_at
	`, content, id)
}

func (s *Store) DeleteComment(ctx context.Context, id int) (domain.Comment, error) {
	return getOne[domain.Comment](s, ctx, `
		DELETE FROM comments WHERE id = $1
		RETURNING id, content, task_id, author_id, created_at, updated_at
	`, id)
}

func (s *Store) FindUserIDByExactName(ctx context.Context, name string) (int, bool, error) {
	var id int
	err := s.db.GetContext(ctx, &id, `SELECT id FROM users WHERE name = $1 LIMIT 1`, name)
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

	return getOne[domain.Attachment](s, ctx, `
		INSERT INTO attachments (file_name, original_name, storage_key, file_size, mime_type, task_id, uploader_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, file_name, original_name, storage_key, file_size, mime_type, task_id, uploader_id, created_at
	`, input.FileName, originalName, storageKey, input.FileSize, input.MimeType, input.TaskID, input.UploaderID)
}

func (s *Store) ListAttachments(ctx context.Context, taskID int) ([]domain.Attachment, error) {
	return selectAll[domain.Attachment](s, ctx, `
		SELECT id, file_name, original_name, storage_key, file_size, mime_type, task_id, uploader_id, created_at
		FROM attachments WHERE task_id = $1
	`, taskID)
}

func (s *Store) GetAttachmentByID(ctx context.Context, id int) (domain.Attachment, error) {
	return getOne[domain.Attachment](s, ctx, `
		SELECT id, file_name, original_name, storage_key, file_size, mime_type, task_id, uploader_id, created_at
		FROM attachments WHERE id = $1
	`, id)
}

func (s *Store) DeleteAttachment(ctx context.Context, id int) (domain.Attachment, error) {
	return getOne[domain.Attachment](s, ctx, `
		DELETE FROM attachments WHERE id = $1
		RETURNING id, file_name, original_name, storage_key, file_size, mime_type, task_id, uploader_id, created_at
	`, id)
}

func (s *Store) CreateNotification(ctx context.Context, input service.NotificationInput) (domain.Notification, error) {
	return getOne[domain.Notification](s, ctx, `
		INSERT INTO notifications (user_id, actor_id, type, task_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, actor_id, type, task_id, is_read, created_at
	`, input.UserID, input.ActorID, input.Type, input.TaskID)
}

func (s *Store) ListNotifications(ctx context.Context, userID int) ([]domain.Notification, error) {
	return selectAll[domain.Notification](s, ctx, `
		SELECT n.id, n.user_id, n.actor_id, u.name AS actor_name, NULL::text AS actor_avatar_url, n.type, n.task_id, n.is_read, n.created_at
		FROM notifications n
		LEFT JOIN users u ON u.id = n.actor_id
		WHERE n.user_id = $1
		ORDER BY n.created_at DESC
	`, userID)
}

func (s *Store) MarkNotificationRead(ctx context.Context, id int, userID int) (domain.Notification, error) {
	return getOne[domain.Notification](s, ctx, `
		UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, actor_id, type, task_id, is_read, created_at
	`, id, userID)
}

func (s *Store) MarkAllNotificationsRead(ctx context.Context, userID int) ([]domain.Notification, error) {
	return selectAll[domain.Notification](s, ctx, `
		UPDATE notifications SET is_read = true WHERE user_id = $1
		RETURNING id, user_id, actor_id, type, task_id, is_read, created_at
	`, userID)
}
