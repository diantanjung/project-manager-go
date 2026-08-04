package domain

import "time"

type Task struct {
	ID                int           `json:"id" db:"id"`
	Title             string        `json:"title" db:"title"`
	Description       *string       `json:"description" db:"description"`
	Status            TaskStatus    `json:"status" db:"status"`
	Priority          *TaskPriority `json:"priority" db:"priority"`
	ProjectID         int           `json:"projectId" db:"project_id"`
	ProjectName       *string       `json:"projectName,omitempty" db:"project_name"`
	CreatorID         int           `json:"creatorId" db:"creator_id"`
	AssigneeID        int           `json:"assigneeId" db:"assignee_id"`
	AssigneeName      *string       `json:"assigneeName,omitempty" db:"assignee_name"`
	AssigneeAvatarURL *string       `json:"assigneeAvatarUrl,omitempty" db:"assignee_avatar_url"`
	DueDate           *string       `json:"dueDate" db:"due_date"`
	Position          *int          `json:"position" db:"position"`
	CreatedAt         *time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt         *time.Time    `json:"updatedAt" db:"updated_at"`
	Comments          []CommentView `json:"comments,omitempty" db:"-"`
	Attachments       []Attachment  `json:"attachments,omitempty" db:"-"`
}

type TaskAssignment struct {
	ID            int        `json:"id" db:"id"`
	TaskID        int        `json:"taskId" db:"task_id"`
	UserID        int        `json:"userId" db:"user_id"`
	UserName      *string    `json:"userName,omitempty" db:"user_name"`
	UserEmail     *string    `json:"userEmail,omitempty" db:"user_email"`
	UserAvatarURL *string    `json:"userAvatarUrl,omitempty" db:"user_avatar_url"`
	AssignedAt    *time.Time `json:"assignedAt" db:"assigned_at"`
}
