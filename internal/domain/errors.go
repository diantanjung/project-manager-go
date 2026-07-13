package domain

import "errors"

var (
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
	Code    string
	Message string
	Err     error
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
