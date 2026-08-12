package httpapi

import (
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

func pageFilter(c *gin.Context) service.PageFilter {
	return service.NormalizePage(service.PageFilter{
		Page:  queryInt(c, "page", 1),
		Limit: queryInt(c, "limit", 10),
	})
}

func pathInt(c *gin.Context, key string) int {
	value, _ := strconv.Atoi(c.Param(key))
	return value
}

func queryInt(c *gin.Context, key string, fallback int) int {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func queryIntPtr(c *gin.Context, key string) *int {
	value := c.Query(key)
	if value == "" {
		return nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &parsed
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
