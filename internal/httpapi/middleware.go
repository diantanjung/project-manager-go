package httpapi

import (
	"strings"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/domain"
	"project-manager-go/internal/service"
)

func (h *Handler) authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			respondError(c, domain.NewError(domain.ErrNoToken, "No token provided"))
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			respondError(c, domain.NewError(domain.ErrMalformedToken, "Malformed token"))
			c.Abort()
			return
		}

		user, err := h.service.Authenticate(parts[1])
		if err != nil {
			respondError(c, domain.NewError(domain.ErrInvalidRefresh, "Invalid token"))
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func (h *Handler) requireRole(role domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := currentUser(c)
		if !service.HasRole(user.Role, role) {
			respondError(c, domain.NewError(domain.ErrForbidden, "Access denied. Insufficient permissions."))
			c.Abort()
			return
		}
		c.Next()
	}
}
