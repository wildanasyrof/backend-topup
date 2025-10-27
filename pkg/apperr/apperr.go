// pkg/apperror/apperror.go
package apperror

import (
	"errors"
	"fmt"
)

type Code string

const (
	// 4xx
	CodeBadRequest    Code = "BAD_REQUEST"
	CodeUnauthorized  Code = "UNAUTHORIZED"
	CodeForbidden     Code = "FORBIDDEN"
	CodeNotFound      Code = "RESOURCE_NOT_FOUND"
	CodeConflict      Code = "RESOURCE_CONFLICT"
	CodeRateLimited   Code = "RATE_LIMITED"
	CodeUnprocessable Code = "UNPROCESSABLE_ENTITY" // semantic/validation errors

	// 5xx
	CodeInternal    Code = "INTERNAL_ERROR"
	CodeUnavailable Code = "SERVICE_UNAVAILABLE"
	CodeTimeout     Code = "REQUEST_TIME_OUT"
)

type AppError struct {
	Code   Code
	Msg    string
	Fields map[string]string // for validation; nil otherwise
	Cause  error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Msg, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func (e *AppError) Unwrap() error { return e.Cause }

func New(code Code, msg string, cause error) *AppError {
	return &AppError{Code: code, Msg: msg, Cause: cause}
}

func Validation(fields map[string]string) *AppError {
	return &AppError{Code: CodeUnprocessable, Msg: "Validation failed", Fields: fields}
}

var (
	ErrNotFound     = &AppError{Code: CodeNotFound, Msg: "Resource not found"}
	ErrUnauthorized = &AppError{Code: CodeUnauthorized, Msg: "Unauthorized"}
	ErrForbidden    = &AppError{Code: CodeForbidden, Msg: "Forbidden"}
	ErrConflict     = &AppError{Code: CodeConflict, Msg: "Conflict"}
)

// helpers
func Is(err error, code Code) bool {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae.Code == code
	}
	return false
}
