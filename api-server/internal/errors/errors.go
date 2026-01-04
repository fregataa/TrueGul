package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code string, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

func Wrap(err error, code string, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

const (
	CodeValidation     = "VALIDATION_ERROR"
	CodeUnauthorized   = "UNAUTHORIZED"
	CodeForbidden      = "FORBIDDEN"
	CodeNotFound       = "NOT_FOUND"
	CodeConflict       = "CONFLICT"
	CodeInternalServer = "INTERNAL_SERVER_ERROR"
	CodeContentTooLong = "CONTENT_TOO_LONG"
)

func Validation(message string) *AppError {
	return New(CodeValidation, message, http.StatusBadRequest)
}

func Unauthorized(message string) *AppError {
	return New(CodeUnauthorized, message, http.StatusUnauthorized)
}

func Forbidden(message string) *AppError {
	return New(CodeForbidden, message, http.StatusForbidden)
}

func NotFound(message string) *AppError {
	return New(CodeNotFound, message, http.StatusNotFound)
}

func Conflict(message string) *AppError {
	return New(CodeConflict, message, http.StatusConflict)
}

func InternalServer(message string) *AppError {
	return New(CodeInternalServer, message, http.StatusInternalServerError)
}

func InternalServerWrap(err error, message string) *AppError {
	return Wrap(err, CodeInternalServer, message, http.StatusInternalServerError)
}

func ContentTooLong(message string) *AppError {
	return New(CodeContentTooLong, message, http.StatusBadRequest)
}

func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}
