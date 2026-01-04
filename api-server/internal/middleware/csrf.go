package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truegul/api-server/internal/dto"
)

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		csrfCookie, err := c.Cookie("csrf_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				ErrorCode: dto.ErrCodeUnauthorized,
				Message:   "CSRF token missing",
			})
			return
		}

		csrfHeader := c.GetHeader("X-CSRF-Token")
		if csrfHeader == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				ErrorCode: dto.ErrCodeUnauthorized,
				Message:   "CSRF token header missing",
			})
			return
		}

		if csrfCookie != csrfHeader {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				ErrorCode: dto.ErrCodeUnauthorized,
				Message:   "CSRF token mismatch",
			})
			return
		}

		c.Next()
	}
}
