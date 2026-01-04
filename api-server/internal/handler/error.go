package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truegul/api-server/internal/dto"
	apperrors "github.com/truegul/api-server/internal/errors"
)

func handleError(c *gin.Context, err error) {
	if appErr, ok := apperrors.IsAppError(err); ok {
		c.JSON(appErr.HTTPStatus, dto.ErrorResponse{
			ErrorCode: appErr.Code,
			Message:   appErr.Message,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
		ErrorCode: apperrors.CodeInternalServer,
		Message:   "An unexpected error occurred",
	})
}

func handleValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, dto.ErrorResponse{
		ErrorCode: apperrors.CodeValidation,
		Message:   message,
	})
}
