package domain

import "time"

type Task struct {
	ID                int           `json:"id"`
	Title             string        `json:"title"`
	Description       *string       `json:"description"`
	Status            TaskStatus    `json:"status"`
	Priority          *TaskPriority `json:"priority"`
	ProjectID         int           `json:"projectId"`
	ProjectName       *string       `json:"projectName,omitempty"`
	CreatorID         int           `json:"creatorId"`
	AssigneeID        int           `json:"assigneeId"`
	AssigneeName      *string       `json:"assigneeName,omitempty"`
	AssigneeAvatarURL *string       `json:"assigneeAvatarUrl,omitempty"`
	DueDate           *string       `json:"dueDate"`
	Position          *int          `json:"position"`
	CreatedAt         *time.Time    `json:"createdAt"`
	UpdatedAt         *time.Time    `json:"updatedAt"`
	Comments          []CommentView `json:"comments,omitempty"`
	Attachments       []Attachment  `json:"attachments,omitempty"`
}

type TaskAssignment struct {
	ID            int        `json:"id"`
	TaskID        int        `json:"taskId"`
	UserID        int        `json:"userId"`
	UserName      *string    `json:"userName,omitempty"`
	UserEmail     *string    `json:"userEmail,omitempty"`
	UserAvatarURL *string    `json:"userAvatarUrl,omitempty"`
	AssignedAt    *time.Time `json:"assignedAt"`
}
