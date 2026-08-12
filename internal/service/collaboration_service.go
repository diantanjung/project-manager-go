package service

import (
	"context"
	"regexp"
	"strings"

	"project-manager-go/internal/domain"
)

var mentionPattern = regexp.MustCompile(`@(\w+)`)

func (s *Service) CreateComment(ctx context.Context, actor domain.AuthUser, input CommentInput) (domain.Comment, error) {
	if err := validatePositiveID(input.TaskID, "Task ID"); err != nil {
		return domain.Comment{}, err
	}
	if strings.TrimSpace(input.Content) == "" {
		return domain.Comment{}, domain.NewError(domain.ErrValidation, "Content is required")
	}
	input.AuthorID = actor.ID
	comment, err := s.store.CreateComment(ctx, input)
	if err != nil {
		return domain.Comment{}, err
	}

	seen := map[string]struct{}{}
	for _, match := range mentionPattern.FindAllStringSubmatch(input.Content, -1) {
		name := match[1]
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}

		userID, ok, err := s.store.FindUserIDByExactName(ctx, name)
		if err != nil {
			return domain.Comment{}, err
		}
		if !ok || userID == actor.ID {
			continue
		}

		actorID := actor.ID
		taskID := input.TaskID
		if _, err := s.store.CreateNotification(ctx, NotificationInput{
			UserID:  userID,
			ActorID: &actorID,
			Type:    domain.NotificationMention,
			TaskID:  &taskID,
		}); err != nil {
			return domain.Comment{}, err
		}
	}
	return comment, nil
}

func (s *Service) ListComments(ctx context.Context, taskID int) ([]domain.CommentView, error) {
	if err := validatePositiveID(taskID, "Task ID"); err != nil {
		return nil, err
	}
	return s.store.ListComments(ctx, taskID)
}

func (s *Service) UpdateComment(ctx context.Context, actor domain.AuthUser, id int, content string) (domain.Comment, error) {
	if err := validatePositiveID(id, "Comment ID"); err != nil {
		return domain.Comment{}, err
	}
	if strings.TrimSpace(content) == "" {
		return domain.Comment{}, domain.NewError(domain.ErrValidation, "Content is required")
	}
	comment, err := s.store.GetCommentByID(ctx, id)
	if err != nil {
		return domain.Comment{}, err
	}
	if comment.AuthorID != actor.ID {
		return domain.Comment{}, domain.NewError(domain.ErrForbidden, "You can only update your own comments")
	}
	return s.store.UpdateComment(ctx, id, content)
}

func (s *Service) DeleteComment(ctx context.Context, actor domain.AuthUser, id int) (domain.Comment, error) {
	if err := validatePositiveID(id, "Comment ID"); err != nil {
		return domain.Comment{}, err
	}
	comment, err := s.store.GetCommentByID(ctx, id)
	if err != nil {
		return domain.Comment{}, err
	}
	if comment.AuthorID != actor.ID {
		return domain.Comment{}, domain.NewError(domain.ErrForbidden, "You can only delete your own comments")
	}
	return s.store.DeleteComment(ctx, id)
}

func (s *Service) CreateAttachment(
	ctx context.Context,
	actor domain.AuthUser,
	taskID int,
	input AttachmentInput,
) (domain.Attachment, error) {
	if err := validatePositiveID(taskID, "Task ID"); err != nil {
		return domain.Attachment{}, err
	}
	input.TaskID = taskID
	input.UploaderID = actor.ID
	return s.store.CreateAttachment(ctx, input)
}

func (s *Service) ListAttachments(ctx context.Context, taskID int) ([]domain.Attachment, error) {
	if err := validatePositiveID(taskID, "Task ID"); err != nil {
		return nil, err
	}
	return s.store.ListAttachments(ctx, taskID)
}

func (s *Service) GetAttachmentByID(ctx context.Context, id int) (domain.Attachment, error) {
	if err := validatePositiveID(id, "Attachment ID"); err != nil {
		return domain.Attachment{}, err
	}
	return s.store.GetAttachmentByID(ctx, id)
}

func (s *Service) DeleteAttachment(ctx context.Context, actor domain.AuthUser, id int) (domain.Attachment, error) {
	if err := validatePositiveID(id, "Attachment ID"); err != nil {
		return domain.Attachment{}, err
	}
	attachment, err := s.store.GetAttachmentByID(ctx, id)
	if err != nil {
		return domain.Attachment{}, err
	}
	if attachment.UploaderID != actor.ID {
		return domain.Attachment{}, domain.NewError(domain.ErrForbidden, "You can only delete your own attachments")
	}
	return s.store.DeleteAttachment(ctx, id)
}

func (s *Service) ListNotifications(ctx context.Context, userID int) ([]domain.Notification, error) {
	if err := validatePositiveID(userID, "User ID"); err != nil {
		return nil, err
	}
	return s.store.ListNotifications(ctx, userID)
}

func (s *Service) MarkNotificationRead(ctx context.Context, actor domain.AuthUser, id int) (domain.Notification, error) {
	if err := validatePositiveID(id, "Notification ID"); err != nil {
		return domain.Notification{}, err
	}
	return s.store.MarkNotificationRead(ctx, id, actor.ID)
}

func (s *Service) MarkAllNotificationsRead(ctx context.Context, actor domain.AuthUser) ([]domain.Notification, error) {
	return s.store.MarkAllNotificationsRead(ctx, actor.ID)
}
