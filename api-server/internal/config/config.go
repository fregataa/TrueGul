package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   time.Duration
	MLServerURL string
	Environment string
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

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		JWTExpiry:   time.Duration(jwtExpiryHours) * time.Hour,
		MLServerURL: getEnv("ML_SERVER_URL", "http://localhost:8000"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
