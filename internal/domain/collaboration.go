package domain

import "time"

type Comment struct {
	ID        int        `json:"id" db:"id"`
	Content   string     `json:"content" db:"content"`
	TaskID    int        `json:"taskId" db:"task_id"`
	AuthorID  int        `json:"authorId" db:"author_id"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

type CommentView struct {
	ID              int        `json:"id" db:"id"`
	Content         string     `json:"content" db:"content"`
	TaskID          int        `json:"taskId,omitempty" db:"task_id"`
	AuthorID        int        `json:"authorId" db:"author_id"`
	AuthorName      *string    `json:"authorName,omitempty" db:"author_name"`
	AuthorAvatarURL *string    `json:"authorAvatarUrl,omitempty" db:"author_avatar_url"`
	CreatedAt       *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       *time.Time `json:"updatedAt" db:"updated_at"`
}

type Attachment struct {
	ID           int        `json:"id" db:"id"`
	FileName     string     `json:"fileName" db:"file_name"`
	OriginalName string     `json:"originalName" db:"original_name"`
	StorageKey   string     `json:"storageKey" db:"storage_key"`
	FileSize     *int       `json:"fileSize" db:"file_size"`
	MimeType     *string    `json:"mimeType" db:"mime_type"`
	TaskID       int        `json:"taskId" db:"task_id"`
	UploaderID   int        `json:"uploaderId" db:"uploader_id"`
	CreatedAt    *time.Time `json:"createdAt" db:"created_at"`
}
