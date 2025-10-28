package routes

import (
	"kubey/server/internal/api/websocket"
	"kubey/server/internal/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, hub *websocket.Hub) {
	// CORS configuration for React frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// WebSocket routes
	router.GET("/ws", websocket.HandleWebSocket(hub))
	router.GET("/ws/status", websocket.HandleWebSocketStatus(hub))

	// API routes
	api := router.Group("/api")
	{
		api.GET("/clusters", handlers.GetClusters)
		api.GET("/clusters/:id", handlers.GetCluster)
		api.GET("/clusters/:id/nodes", handlers.GetClusterNodes)
		api.GET("/clusters/:id/pods", handlers.GetClusterPods)
		api.GET("/clusters/:id/services", handlers.GetClusterServices)
		api.GET("/clusters/:id/deployments", handlers.GetClusterDeployments)
		api.GET("/clusters/:id/namespaces", handlers.GetClusterNamespaces)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})
}
