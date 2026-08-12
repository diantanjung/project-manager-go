package httpapi

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"project-manager-go/internal/domain"
)

type responseEnvelope struct {
	Data       any                `json:"data"`
	Pagination *domain.Pagination `json:"pagination,omitempty"`
	Meta       any                `json:"meta,omitempty"`
}

type errorEnvelope struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

type paginatedPayload interface {
	Items() any
	PageInfo() domain.Pagination
}

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
	respondData(c, status, payload, nil)
}

func respondWithMeta(c *gin.Context, status int, payload any, meta any, err error) {
	if err != nil {
		respondError(c, err)
		return
	}
	respondData(c, status, payload, meta)
}

func respondData(c *gin.Context, status int, payload any, meta any) {
	envelope := responseEnvelope{Data: payload, Meta: meta}
	if paginated, ok := payload.(paginatedPayload); ok {
		pageInfo := paginated.PageInfo()
		envelope.Data = paginated.Items()
		envelope.Pagination = &pageInfo
	}
	c.JSON(status, envelope)
}

func respondError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "Internal server error"

	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		switch {
		case errors.Is(appErr.Err, domain.ErrInvalidCredential),
			errors.Is(appErr.Err, domain.ErrInvalidRefresh),
			errors.Is(appErr.Err, domain.ErrNoToken),
			errors.Is(appErr.Err, domain.ErrMalformedToken):
			status = http.StatusUnauthorized
		case errors.Is(appErr.Err, domain.ErrForbidden):
			status = http.StatusForbidden
		case errors.Is(appErr.Err, domain.ErrDuplicate):
			status = http.StatusConflict
		case errors.Is(appErr.Err, domain.ErrValidation):
			status = http.StatusUnprocessableEntity
		default:
			status = http.StatusInternalServerError
		}
		message = appErr.Message
	} else if errors.Is(err, domain.ErrNotFound) {
		status = http.StatusNotFound
		message = "Resource not found"
	}

	if message == "" {
		message = http.StatusText(status)
	}
	c.JSON(status, errorEnvelope{Message: message})
}
