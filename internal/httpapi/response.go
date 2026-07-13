package httpapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/domain"
)

func bindJSON(c *gin.Context, dest any) bool {
	if err := c.ShouldBindJSON(dest); err != nil {
		respondError(c, domain.NewError(domain.ErrValidation, err.Error()))
		return false
	}
	return true
}

func respond(c *gin.Context, status int, payload any, err error) {
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(status, payload)
}

func respondError(c *gin.Context, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		switch {
		case errors.Is(appErr.Err, domain.ErrInvalidCredential),
			errors.Is(appErr.Err, domain.ErrInvalidRefresh),
			errors.Is(appErr.Err, domain.ErrNoToken),
			errors.Is(appErr.Err, domain.ErrMalformedToken):
			c.JSON(http.StatusUnauthorized, gin.H{"message": appErr.Message})
		case errors.Is(appErr.Err, domain.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"message": appErr.Message})
		case errors.Is(appErr.Err, domain.ErrDuplicate):
			c.JSON(http.StatusConflict, gin.H{"message": appErr.Message})
		case errors.Is(appErr.Err, domain.ErrValidation):
			c.JSON(http.StatusBadRequest, gin.H{"message": appErr.Message})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": appErr.Message})
		}
		return
	}
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Resource not found"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
}
