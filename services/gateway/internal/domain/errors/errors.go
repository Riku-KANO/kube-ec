package errors

import "errors"

var (
	// Domain errors
	ErrInvalidInput     = errors.New("invalid input")
	ErrUserNotFound     = errors.New("user not found")
	ErrNotFound         = errors.New("not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrEmailExists      = errors.New("email already exists")
	ErrInternalError    = errors.New("internal server error")
)
