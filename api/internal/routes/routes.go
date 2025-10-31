package routes

import (
	"kubey/api/internal/config"
	"kubey/api/internal/handlers"
	"kubey/api/internal/handlers/clusters"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, cfg *config.ApiConfig) {
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
	router.GET("/health", handlers.Health)
}
