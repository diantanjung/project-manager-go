package httpapi

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func currentUser(c *gin.Context) domain.AuthUser {
	user, _ := c.Get("user")
	authUser, _ := user.(domain.AuthUser)
	return authUser
}

func pageFilter(c *gin.Context) (service.PageFilter, bool) {
	page, ok := queryInt(c, "page", 1)
	if !ok {
		return service.PageFilter{}, false
	}
	limit, ok := queryInt(c, "limit", 10)
	if !ok {
		return service.PageFilter{}, false
	}
	return service.NormalizePage(service.PageFilter{
		Page:  page,
		Limit: limit,
	}), true
}

func pathInt(c *gin.Context, key string) (int, bool) {
	value, err := strconv.Atoi(c.Param(key))
	if err != nil || value < 1 {
		respondBadRequest(c, key, "The "+key+" parameter must be a positive integer.")
		return 0, false
	}
	return value, true
}

func queryInt(c *gin.Context, key string, fallback int) (int, bool) {
	value := c.Query(key)
	if value == "" {
		return fallback, true
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		respondBadRequest(c, key, "The "+key+" query parameter must be a positive integer.")
		return 0, false
	}
	return parsed, true
}

func queryIntPtr(c *gin.Context, key string) (*int, bool) {
	value := c.Query(key)
	if value == "" {
		return nil, true
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		respondBadRequest(c, key, "The "+key+" query parameter must be a positive integer.")
		return nil, false
	}
	return &parsed, true
}

func respondBadRequest(c *gin.Context, field string, message string) {
	respondError(c, domain.NewFieldError(
		domain.ErrBadRequest,
		http.StatusText(http.StatusBadRequest),
		map[string][]string{
			field: {message},
		},
	))
}

func parseRole(value string) *domain.UserRole {
	if value == "" {
		return nil
	}
	role := domain.UserRole(value)
	if !role.IsValid() {
		return nil
	}
	return &role
}

func parseTaskStatus(value string) *domain.TaskStatus {
	if value == "" {
		return nil
	}
	status := domain.TaskStatus(value)
	if !status.IsValid() {
		return nil
	}
	return &status
}

func parseTaskPriority(value string) *domain.TaskPriority {
	if value == "" {
		return nil
	}
	priority := domain.TaskPriority(value)
	if !priority.IsValid() {
		return nil
	}
	return &priority
}
