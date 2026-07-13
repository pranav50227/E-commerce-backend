package errors

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrBadRequest    = errors.New("bad request parameters")
	ErrInternalError = errors.New("internal server error")
)
