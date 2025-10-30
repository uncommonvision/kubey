package security

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"kubey/server/internal/config"
)

// CORS returns a Gin middleware for handling Cross-Origin Resource Sharing (CORS)
func CORS(cfg *config.ApiConfig) gin.HandlerFunc {
	// Create CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Return the CORS middleware
	return cors.New(corsConfig)
}
