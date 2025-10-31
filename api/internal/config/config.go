package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type ApiConfig struct {
	Environment      string
	Host             string
	Port             string
	LogLevel         string
	KubeConfig       string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration
	AllowedOrigins   []string
	RequestIDHeader  string
}

func LoadApi() *ApiConfig {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &ApiConfig{
		Environment:      getEnv("ENVIRONMENT", "development"),
		Host:             getEnv("HOST", "localhost"),
		Port:             getEnv("PORT", "8080"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		KubeConfig:       getEnv("KUBECONFIG", ""), // Use default kubeconfig if not specified
		HTTPReadTimeout:  getDurationEnv("HTTP_READ_TIMEOUT", 10*time.Second),
		HTTPWriteTimeout: getDurationEnv("HTTP_WRITE_TIMEOUT", 10*time.Second),
		HTTPIdleTimeout:  getDurationEnv("HTTP_IDLE_TIMEOUT", 30*time.Second),
		AllowedOrigins:   getSliceEnv("ALLOWED_ORIGINS", []string{"http://localhost:5173", "http://127.0.0.1:5173"}),
		RequestIDHeader:  getEnv("REQUEST_ID_HEADER", "X-Request-ID"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim spaces
		items := strings.Split(value, ",")
		slice := make([]string, 0)
		for _, item := range items {
			trimmed := strings.TrimSpace(item)
			if trimmed != "" {
				slice = append(slice, trimmed)
			}
		}
		// If we got valid items, return them; otherwise return default
		if len(slice) > 0 {
			return slice
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		// Parse as seconds
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
		log.Printf("Invalid duration for %s: %s, using default", key, value)
	}
	return defaultValue
}
