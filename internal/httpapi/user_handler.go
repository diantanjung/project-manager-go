package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/service"
)

func (h *Handler) listUsers(c *gin.Context) {
	filter, ok := pageFilter(c)
	if !ok {
		return
	}
	role := parseRole(c.Query("role"))
	result, err := h.service.ListUsers(c.Request.Context(), service.ListUsersFilter{
		PageFilter: filter,
		Search:     c.Query("search"),
		Role:       role,
		SortBy:     c.DefaultQuery("sortBy", "createdAt"),
		Order:      c.DefaultQuery("order", "desc"),
	})
	respond(c, http.StatusOK, result, err)
}

func (h *Handler) createUser(c *gin.Context) {
	var input service.CreateUserInput
	if !bindJSON(c, &input) {
		return
	}
	user, err := h.service.CreateUser(c.Request.Context(), input)
	respond(c, http.StatusCreated, user, err)
}

func (h *Handler) getUser(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	user, err := h.service.GetUserByID(c.Request.Context(), id)
	respond(c, http.StatusOK, user, err)
}

func (h *Handler) updateUser(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	var input service.UpdateUserInput
	if !bindJSON(c, &input) {
		return
	}
	user, err := h.service.UpdateUser(c.Request.Context(), currentUser(c), id, input)
	respond(c, http.StatusOK, user, err)
}

func (h *Handler) deleteUser(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	_, err := h.service.DeleteUser(c.Request.Context(), currentUser(c), id)
	respond(c, http.StatusOK, gin.H{"message": "User deleted successfully"}, err)
}

func (h *Handler) listUserTasks(c *gin.Context) {
	id, ok := pathInt(c, "id")
	if !ok {
		return
	}
	filter, ok := pageFilter(c)
	if !ok {
		return
	}
	result, err := h.service.ListUserTasks(c.Request.Context(), id, filter)
	respond(c, http.StatusOK, result, err)
}
