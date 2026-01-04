package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/truegul/api-server/internal/config"
	"github.com/truegul/api-server/internal/database"
	"github.com/truegul/api-server/internal/handler"
	"github.com/truegul/api-server/internal/middleware"
	"github.com/truegul/api-server/internal/mq"
	"github.com/truegul/api-server/internal/repository"
	"github.com/truegul/api-server/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	publisher, err := mq.NewRedisPublisher(cfg.RedisURL, cfg.StreamName)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer publisher.Close()

	userRepo := repository.NewUserRepository(db)
	writingRepo := repository.NewWritingRepository(db)
	analysisRepo := repository.NewAnalysisRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	writingService := service.NewWritingService(writingRepo)
	analysisService := service.NewAnalysisService(analysisRepo, writingRepo, userRepo, publisher, cfg)

	authHandler := handler.NewAuthHandler(authService, cfg.Environment)
	writingHandler := handler.NewWritingHandler(writingService)
	analysisHandler := handler.NewAnalysisHandler(analysisService, cfg)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
		}

		internal := v1.Group("/internal")
		{
			internal.POST("/callback", analysisHandler.Callback)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		protected.Use(middleware.CSRFMiddleware())
		{
			protected.GET("/auth/me", authHandler.Me)

			writings := protected.Group("/writings")
			{
				writings.POST("", writingHandler.Create)
				writings.GET("", writingHandler.List)
				writings.GET("/:id", writingHandler.GetByID)
				writings.PUT("/:id", writingHandler.Update)
				writings.DELETE("/:id", writingHandler.Delete)
				writings.POST("/:id/submit", analysisHandler.Submit)
				writings.GET("/:id/analysis", analysisHandler.GetAnalysis)
			}
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
