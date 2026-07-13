package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (h *Handler) listTasks(c *gin.Context) {
	result, err := h.service.ListTasks(c.Request.Context(), currentUser(c), service.TaskFilter{
		PageFilter: pageFilter(c),
		Search:     c.Query("search"),
		ProjectID:  queryIntPtr(c, "projectId"),
		Status:     parseTaskStatus(c.Query("status")),
		Priority:   parseTaskPriority(c.Query("priority")),
		AssigneeID: queryIntPtr(c, "assigneeId"),
		SortBy:     c.DefaultQuery("sortBy", "createdAt"),
		Order:      c.DefaultQuery("order", "desc"),
	})
	respond(c, http.StatusOK, result, err)
}

func (h *Handler) createTask(c *gin.Context) {
	var input service.TaskInput
	if !bindJSON(c, &input) {
		return
	}
	task, err := h.service.CreateTask(c.Request.Context(), currentUser(c), input)
	respond(c, http.StatusCreated, task, err)
}

func (h *Handler) getTask(c *gin.Context) {
	task, err := h.service.GetTaskByID(c.Request.Context(), currentUser(c), pathInt(c, "id"))
	respond(c, http.StatusOK, task, err)
}

func (h *Handler) updateTask(c *gin.Context) {
	var input service.TaskPatchInput
	if !bindJSON(c, &input) {
		return
	}
	task, err := h.service.UpdateTask(c.Request.Context(), currentUser(c), pathInt(c, "id"), input)
	respond(c, http.StatusOK, task, err)
}

func (h *Handler) updateTaskStatus(c *gin.Context) {
	var input struct {
		Status   domain.TaskStatus `json:"status"`
		Position *int              `json:"position"`
	}
	if !bindJSON(c, &input) {
		return
	}
	task, err := h.service.UpdateTask(
		c.Request.Context(),
		currentUser(c),
		pathInt(c, "id"),
		service.TaskPatchInput{
			Status:   &input.Status,
			Position: input.Position,
		},
	)
	respond(c, http.StatusOK, task, err)
}

func (h *Handler) deleteTask(c *gin.Context) {
	_, err := h.service.DeleteTask(c.Request.Context(), currentUser(c), pathInt(c, "id"))
	respond(c, http.StatusOK, gin.H{"message": "Task deleted successfully"}, err)
}

func (h *Handler) listTaskAssignments(c *gin.Context) {
	assignments, err := h.service.ListTaskAssignments(c.Request.Context(), pathInt(c, "taskId"))
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
	_, err := h.service.RemoveTaskAssignment(c.Request.Context(), currentUser(c), pathInt(c, "id"))
	respond(c, http.StatusOK, gin.H{"message": "Assignment removed successfully"}, err)
}
