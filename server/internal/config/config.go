package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Host        string
	Port        string
	LogLevel    string
	KubeConfig  string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Host:        getEnv("HOST", "localhost"),
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		KubeConfig:  getEnv("KUBECONFIG", ""), // Use default kubeconfig if not specified
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
