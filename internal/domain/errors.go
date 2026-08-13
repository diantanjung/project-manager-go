package domain

import "errors"

var (
	ErrBadRequest        = errors.New("bad request")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrInvalidRefresh    = errors.New("invalid refresh token")
	ErrMalformedToken    = errors.New("malformed token")
	ErrNoToken           = errors.New("no token provided")
	ErrNotFound          = errors.New("not found")
	ErrRefreshNotOwned   = errors.New("refresh token not found or does not belong to user")
	ErrDuplicate         = errors.New("duplicate")
	ErrValidation        = errors.New("validation failed")
)

type AppError struct {
	Code        string
	Message     string
	FieldErrors map[string][]string
	Err         error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewError(err error, message string) error {
	return &AppError{
		Message: message,
		Err:     err,
	}
}

func NewValidationError(message string, fieldErrors map[string][]string) error {
	return NewFieldError(ErrValidation, message, fieldErrors)
}

func NewFieldError(err error, message string, fieldErrors map[string][]string) error {
	return &AppError{
		Message:     message,
		FieldErrors: fieldErrors,
		Err:         err,
	}
}
