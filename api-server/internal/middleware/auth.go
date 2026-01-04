package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truegul/api-server/internal/dto"
	"github.com/truegul/api-server/internal/service"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCookie, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				ErrorCode: dto.ErrCodeUnauthorized,
				Message:   "Authentication required",
			})
			return
		}

		claims, err := authService.ValidateToken(tokenCookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				ErrorCode: dto.ErrCodeUnauthorized,
				Message:   "Invalid or expired token",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
