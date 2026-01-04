package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truegul/api-server/internal/dto"
	apperrors "github.com/truegul/api-server/internal/errors"
	"github.com/truegul/api-server/internal/service"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCookie, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				ErrorCode: apperrors.CodeUnauthorized,
				Message:   "Authentication required",
			})
			return
		}

		claims, err := authService.ValidateToken(tokenCookie)
		if err != nil {
			if appErr, ok := apperrors.IsAppError(err); ok {
				c.AbortWithStatusJSON(appErr.HTTPStatus, dto.ErrorResponse{
					ErrorCode: appErr.Code,
					Message:   appErr.Message,
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				ErrorCode: apperrors.CodeUnauthorized,
				Message:   "Invalid or expired token",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
