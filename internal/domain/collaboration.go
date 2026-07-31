package domain

import "time"

type Comment struct {
	ID        int        `json:"id"`
	Content   string     `json:"content"`
	TaskID    int        `json:"taskId"`
	AuthorID  int        `json:"authorId"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type CommentView struct {
	ID              int        `json:"id"`
	Content         string     `json:"content"`
	TaskID          int        `json:"taskId,omitempty"`
	AuthorID        int        `json:"authorId"`
	AuthorName      *string    `json:"authorName,omitempty"`
	AuthorAvatarURL *string    `json:"authorAvatarUrl,omitempty"`
	CreatedAt       *time.Time `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt"`
}

type Attachment struct {
	ID           int        `json:"id"`
	FileName     string     `json:"fileName"`
	OriginalName string     `json:"originalName"`
	StorageKey   string     `json:"storageKey"`
	FileSize     *int       `json:"fileSize"`
	MimeType     *string    `json:"mimeType"`
	TaskID       int        `json:"taskId"`
	UploaderID   int        `json:"uploaderId"`
	CreatedAt    *time.Time `json:"createdAt"`
}
