package routes

import (
	"kubey/server/internal/config"
	"kubey/server/internal/handlers/clusters"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, cfg *config.ApiConfig) {
	// CORS configuration
	var corsConfig cors.Config

	if cfg.Environment == "production" {
		// Production: Only allow specific origins from env var
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "https://yourdomain.com"  // Change this default
		}
		origins := strings.Split(allowedOrigins, ",")
		corsConfig = cors.Config{
			AllowOrigins:     origins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}
	} else {
		// Development: Allow all localhost on any port
		corsConfig = cors.Config{
			AllowOriginFunc: func(origin string) bool {
				return strings.HasPrefix(origin, "http://localhost:") ||
					strings.HasPrefix(origin, "https://localhost:") ||
					strings.HasPrefix(origin, "http://127.0.0.1:") ||
					strings.HasPrefix(origin, "https://127.0.0.1:")
			},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}
	}

	router.Use(cors.New(corsConfig))

	// API routes
	api := router.Group("/api")
	{
		api.GET("/clusters", clusters.GetClusters)
		api.GET("/clusters/:id", clusters.GetCluster)
		api.GET("/clusters/:id/nodes", clusters.GetClusterNodes)
		api.GET("/clusters/:id/pods", clusters.GetClusterPods)
		api.GET("/clusters/:id/services", clusters.GetClusterServices)
		api.GET("/clusters/:id/deployments", clusters.GetClusterDeployments)
		api.GET("/clusters/:id/namespaces", clusters.GetClusterNamespaces)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})
}
