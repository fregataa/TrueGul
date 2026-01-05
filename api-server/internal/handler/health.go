package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db     *gorm.DB
	redis  *redis.Client
}

type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

func NewHealthHandler(db *gorm.DB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{
		db:     db,
		redis:  redis,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	services := make(map[string]string)
	overallStatus := "healthy"

	// Check database
	sqlDB, err := h.db.DB()
	if err != nil {
		services["database"] = "unhealthy"
		overallStatus = "unhealthy"
	} else if err := sqlDB.PingContext(ctx); err != nil {
		services["database"] = "unhealthy"
		overallStatus = "unhealthy"
	} else {
		services["database"] = "healthy"
	}

	// Check Redis
	if h.redis != nil {
		if err := h.redis.Ping(ctx).Err(); err != nil {
			services["redis"] = "unhealthy"
			overallStatus = "unhealthy"
		} else {
			services["redis"] = "healthy"
		}
	} else {
		services["redis"] = "not configured"
	}

	status := http.StatusOK
	if overallStatus != "healthy" {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, HealthResponse{
		Status:   overallStatus,
		Services: services,
	})
}

func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

func (h *HealthHandler) Readiness(c *gin.Context) {
	h.Check(c)
}
