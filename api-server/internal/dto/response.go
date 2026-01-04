package dto

import "github.com/google/uuid"

type ErrorResponse struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

type AuthResponse struct {
	User      UserResponse `json:"user"`
	CSRFToken string       `json:"csrf_token"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// Error codes
const (
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeUserExists     = "USER_EXISTS"
	ErrCodeInternalServer = "INTERNAL_SERVER_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
)
