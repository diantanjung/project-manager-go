package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

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
		respondError(c, bindingError(dest, err))
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
	fieldErrors := map[string][]string(nil)

	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		switch {
		case errors.Is(appErr.Err, domain.ErrBadRequest):
			status = http.StatusBadRequest
		case errors.Is(appErr.Err, domain.ErrInvalidCredential),
			errors.Is(appErr.Err, domain.ErrInvalidRefresh),
			errors.Is(appErr.Err, domain.ErrNoToken),
			errors.Is(appErr.Err, domain.ErrMalformedToken):
			status = http.StatusUnauthorized
		case errors.Is(appErr.Err, domain.ErrForbidden):
			status = http.StatusForbidden
		case errors.Is(appErr.Err, domain.ErrDuplicate):
			status = http.StatusConflict
		case errors.Is(appErr.Err, domain.ErrNotFound):
			status = http.StatusNotFound
		case errors.Is(appErr.Err, domain.ErrValidation):
			status = http.StatusUnprocessableEntity
		default:
			status = http.StatusInternalServerError
		}
		message = appErr.Message
		fieldErrors = appErr.FieldErrors
	} else if errors.Is(err, domain.ErrNotFound) {
		status = http.StatusNotFound
		message = "Resource not found"
	}

	if message == "" {
		message = http.StatusText(status)
	}
	c.JSON(status, errorEnvelope{
		Message: message,
		Errors:  fieldErrors,
	})
}

func bindingError(dest any, err error) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return domain.NewValidationError(
			"The given data was invalid.",
			formatValidationErrors(dest, validationErrors),
		)
	}

	var syntaxError *json.SyntaxError
	var typeError *json.UnmarshalTypeError
	if errors.As(err, &syntaxError) || errors.As(err, &typeError) {
		return domain.NewValidationError(
			"The given data was invalid.",
			map[string][]string{
				"body": {"The request body must be valid JSON."},
			},
		)
	}

	return domain.NewValidationError(
		"The given data was invalid.",
		map[string][]string{
			"body": {"The request body is invalid."},
		},
	)
}

func formatValidationErrors(dest any, validationErrors validator.ValidationErrors) map[string][]string {
	fieldErrors := make(map[string][]string, len(validationErrors))
	for _, validationError := range validationErrors {
		field := jsonFieldName(dest, validationError.StructField())
		fieldErrors[field] = append(fieldErrors[field], validationMessage(field, validationError))
	}
	return fieldErrors
}

func jsonFieldName(dest any, structField string) string {
	destType := reflect.TypeOf(dest)
	for destType.Kind() == reflect.Pointer {
		destType = destType.Elem()
	}
	if destType.Kind() == reflect.Struct {
		if field, ok := destType.FieldByName(structField); ok {
			jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
			if jsonName != "" && jsonName != "-" {
				return jsonName
			}
		}
	}
	return lowerFirst(structField)
}

func validationMessage(field string, validationError validator.FieldError) string {
	switch validationError.Tag() {
	case "required":
		return "The " + field + " field is required."
	case "email":
		return "The " + field + " field must be a valid email address."
	case "min":
		return "The " + field + " field must be at least " + validationError.Param() + " characters."
	case "oneof":
		return "The " + field + " field is invalid."
	case "gt":
		return "The " + field + " field must be greater than " + validationError.Param() + "."
	case "gte":
		return "The " + field + " field must be greater than or equal to " + validationError.Param() + "."
	case "datetime":
		return "The " + field + " field must be a valid date."
	default:
		return "The " + field + " field is invalid."
	}
}

func lowerFirst(value string) string {
	if value == "" {
		return value
	}
	return strings.ToLower(value[:1]) + value[1:]
}
