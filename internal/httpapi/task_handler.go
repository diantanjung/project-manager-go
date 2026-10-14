package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (h *Handler) listTasks(c *gin.Context) {
	board := c.Query("page") == "" && c.Query("limit") == ""
	filter := service.PageFilter{}
	if !board {
		var ok bool
		filter, ok = pageFilter(c)
		if !ok {
			return
		}
	}
	projectID, ok := queryIntPtr(c, "projectId")
	if !ok {
		return
	}
	legacyProjectID, ok := queryIntPtr(c, "project_id")
	if !ok {
		return
	}
	if projectID == nil {
		projectID = legacyProjectID
	} else if legacyProjectID != nil && *projectID != *legacyProjectID {
		respondBadRequest(c, "projectId", "projectId and project_id must match.")
		return
	}
	assigneeID, ok := queryIntPtr(c, "assigneeId")
	if !ok {
		return
	}
	result, err := h.service.ListTasks(c.Request.Context(), currentUser(c), service.TaskFilter{
		PageFilter: filter,
		Search:     c.Query("search"),
		ProjectID:  projectID,
		Status:     parseTaskStatus(c.Query("status")),
		Priority:   parseTaskPriority(c.Query("priority")),
		AssigneeID: assigneeID,
		DueBefore:  queryStringPtr(c, "dueBefore"),
		DueAfter:   queryStringPtr(c, "dueAfter"),
		Board:      board,
		SortBy:     c.DefaultQuery("sortBy", "createdAt"),
		Order:      c.DefaultQuery("order", "desc"),
	})
	if board {
		if err != nil {
			respondError(c, err)
			return
		}
		respond(c, http.StatusOK, result.Data, nil)
		return
	}
	respond(c, http.StatusOK, result, err)
}

func queryStringPtr(c *gin.Context, key string) *string {
	value := c.Query(key)
	if value == "" {
		return nil
	}
	return &value
}

func (h *Handler) reorderTasks(c *gin.Context) {
	var input service.TaskReorderInput
	if !bindJSON(c, &input) {
		return
	}
	tasks, err := h.service.ReorderTasks(c.Request.Context(), currentUser(c), input)
	respond(c, http.StatusOK, gin.H{"tasks": tasks}, err)
}

func (h *Handler) listTaskActivity(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	activity, err := h.service.ListTaskActivity(c.Request.Context(), currentUser(c), id)
	respond(c, http.StatusOK, activity, err)
}

func (h *Handler) assignUserToTaskByPath(c *gin.Context) {
	taskID, ok := pathInt(c, "id")
	if !ok {
		return
	}
	var body struct {
		UserID int `json:"userId" binding:"required,gt=0"`
	}
	if !bindJSON(c, &body) {
		return
	}
	assignment, exists, err := h.service.AssignUserToTask(c.Request.Context(), currentUser(c), service.TaskAssignmentInput{
		TaskID: taskID,
		UserID: body.UserID,
	})
	if exists && err == nil {
		err = domain.NewError(domain.ErrDuplicate, "User is already assigned to this task")
	}
	respond(c, http.StatusCreated, assignment, err)
}

func (h *Handler) removeTaskAssignmentByUser(c *gin.Context) {
	taskID, ok := pathInt(c, "id")
	if !ok {
		return
	}
	userID, ok := pathInt(c, "userId")
	if !ok {
		return
	}
	_, err := h.service.RemoveTaskAssignmentByUser(c.Request.Context(), currentUser(c), taskID, userID)
	respond(c, http.StatusOK, gin.H{"message": "Assignment removed successfully"}, err)
}

func (h *Handler) createTask(c *gin.Context) {
	var input service.TaskInput
	if !bindTaskPayload(c, &input, true) {
		return
	}
	task, err := h.service.CreateTask(c.Request.Context(), currentUser(c), input)
	respond(c, http.StatusCreated, task, err)
}

func (h *Handler) getTask(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	task, err := h.service.GetTaskByID(c.Request.Context(), currentUser(c), id)
	respond(c, http.StatusOK, task, err)
}

func (h *Handler) updateTask(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	var input service.TaskPatchInput
	if !bindTaskPayload(c, &input, false) {
		return
	}
	task, err := h.service.UpdateTask(c.Request.Context(), currentUser(c), id, input)
	respond(c, http.StatusOK, task, err)
}

func bindTaskPayload(c *gin.Context, dest any, create bool) bool {
	var body map[string]json.RawMessage
	if err := c.ShouldBindJSON(&body); err != nil {
		respondError(c, bindingError(dest, err))
		return false
	}
	if body == nil {
		respondError(c, domain.NewValidationError("The given data was invalid.", map[string][]string{
			"body": {"The request body must be a JSON object."},
		}))
		return false
	}
	if create {
		copyTaskAlias(body, "projectId", "project_id")
	}
	clearAssignee := false
	if !create {
		for _, key := range []string{"assigneeId", "assignee_id"} {
			if string(body[key]) == "0" || string(body[key]) == "null" {
				clearAssignee = true
				delete(body, key)
			}
		}
	}
	copyTaskAlias(body, "assigneeId", "assignee_id")
	copyTaskAlias(body, "dueDate", "due_date")
	if raw, ok := body["dueDate"]; ok {
		var date string
		if err := json.Unmarshal(raw, &date); err == nil && len(date) > 10 {
			if _, err := time.Parse(time.RFC3339Nano, date); err == nil {
				body["dueDate"], _ = json.Marshal(date[:10])
			}
		}
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		respondError(c, bindingError(dest, err))
		return false
	}
	if err := json.Unmarshal(encoded, dest); err != nil {
		respondError(c, bindingError(dest, err))
		return false
	}
	if err := binding.Validator.ValidateStruct(dest); err != nil {
		respondError(c, bindingError(dest, err))
		return false
	}
	if clearAssignee {
		if patch, ok := dest.(*service.TaskPatchInput); ok {
			patch.ClearAssignee = true
		}
	}
	return true
}

func copyTaskAlias(body map[string]json.RawMessage, target, alias string) {
	if _, exists := body[target]; exists {
		return
	}
	if value, exists := body[alias]; exists {
		body[target] = value
	}
}

func (h *Handler) updateTaskStatus(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	var input struct {
		Status   domain.TaskStatus `json:"status" binding:"required,oneof=backlog todo in_progress review done"`
		Position *int              `json:"position" binding:"omitempty,gte=0"`
	}
	if !bindJSON(c, &input) {
		return
	}
	task, err := h.service.UpdateTask(
		c.Request.Context(),
		currentUser(c),
		id,
		service.TaskPatchInput{
			Status:   &input.Status,
			Position: input.Position,
		},
	)
	respond(c, http.StatusOK, task, err)
}

func (h *Handler) deleteTask(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	_, err := h.service.DeleteTask(c.Request.Context(), currentUser(c), id)
	respond(c, http.StatusOK, gin.H{"message": "Task deleted successfully"}, err)
}

func (h *Handler) listTaskAssignments(c *gin.Context) {
	taskID, ok := pathInt(c, "taskId")
	if !ok {
		return
	}
	assignments, err := h.service.ListTaskAssignments(c.Request.Context(), currentUser(c), taskID)
	respond(c, http.StatusOK, assignments, err)
}

func (h *Handler) assignUserToTask(c *gin.Context) {
	var input service.TaskAssignmentInput
	if !bindJSON(c, &input) {
		return
	}
	assignment, exists, err := h.service.AssignUserToTask(c.Request.Context(), currentUser(c), input)
	if exists && err == nil {
		err = domain.NewError(domain.ErrDuplicate, "User is already assigned to this task")
	}
	respond(c, http.StatusCreated, assignment, err)
}

func (h *Handler) removeTaskAssignment(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	_, err := h.service.RemoveTaskAssignment(c.Request.Context(), currentUser(c), id)
	respond(c, http.StatusOK, gin.H{"message": "Assignment removed successfully"}, err)
}
