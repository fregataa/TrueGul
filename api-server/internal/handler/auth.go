package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/truegul/api-server/internal/dto"
	"github.com/truegul/api-server/internal/service"
)

type AuthHandler struct {
	authService  *service.AuthService
	isProduction bool
}

func NewAuthHandler(authService *service.AuthService, environment string) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		isProduction: environment == "production",
	}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			ErrorCode: dto.ErrCodeValidation,
			Message:   err.Error(),
		})
		return
	}

	user, err := h.authService.Signup(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				ErrorCode: dto.ErrCodeUserExists,
				Message:   "User with this email already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			ErrorCode: dto.ErrCodeInternalServer,
			Message:   "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			ErrorCode: dto.ErrCodeValidation,
			Message:   err.Error(),
		})
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				ErrorCode: dto.ErrCodeUnauthorized,
				Message:   "Invalid email or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			ErrorCode: dto.ErrCodeInternalServer,
			Message:   "Failed to login",
		})
		return
	}

	csrfToken, err := service.GenerateCSRFToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			ErrorCode: dto.ErrCodeInternalServer,
			Message:   "Failed to generate CSRF token",
		})
		return
	}

	c.SetCookie(
		"token",
		token,
		3600,
		"/",
		"",
		h.isProduction,
		true,
	)

	c.SetCookie(
		"csrf_token",
		csrfToken,
		3600,
		"/",
		"",
		h.isProduction,
		false,
	)

	c.JSON(http.StatusOK, dto.AuthResponse{
		User: dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
		CSRFToken: csrfToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", h.isProduction, true)
	c.SetCookie("csrf_token", "", -1, "/", "", h.isProduction, false)

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Logged out successfully",
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			ErrorCode: dto.ErrCodeUnauthorized,
			Message:   "Not authenticated",
		})
		return
	}

	email, _ := c.Get("email")

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:    userID.(uuid.UUID),
		Email: email.(string),
	})
}
