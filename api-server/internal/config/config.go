package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port             string
	DatabaseURL      string
	RedisURL         string
	JWTSecret        string
	JWTExpiry        time.Duration
	MLServerURL      string
	MLCallbackSecret string
	Environment      string
	StreamName       string
	CallbackBaseURL  string
	CallbackPath     string
	CORSOrigins      []string
}

func Load() *Config {
	jwtExpiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "1"))

	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	port := getEnv("PORT", "8080")
	callbackBaseURL := getEnv("CALLBACK_BASE_URL", "http://localhost:"+port)

	corsOriginsStr := getEnv("CORS_ORIGINS", "http://localhost:3000")
	corsOrigins := strings.Split(corsOriginsStr, ",")
	for i := range corsOrigins {
		corsOrigins[i] = strings.TrimSpace(corsOrigins[i])
	}

	return &Config{
		Port:             port,
		DatabaseURL:      databaseURL,
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:        jwtSecret,
		JWTExpiry:        time.Duration(jwtExpiryHours) * time.Hour,
		MLServerURL:      getEnv("ML_SERVER_URL", "http://localhost:8000"),
		MLCallbackSecret: getEnv("ML_CALLBACK_SECRET", ""),
		Environment:      getEnv("ENVIRONMENT", "development"),
		StreamName:       getEnv("STREAM_NAME", "analysis_tasks"),
		CallbackBaseURL:  callbackBaseURL,
		CallbackPath:     "/api/v1/internal/callback",
		CORSOrigins:      corsOrigins,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
