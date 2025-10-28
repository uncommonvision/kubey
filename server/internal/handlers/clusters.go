package handlers

import (
	"net/http"

	"kubey/server/internal/services"

	"github.com/gin-gonic/gin"
)

// GetClusters returns all clusters
func GetClusters(c *gin.Context) {
	clusters, err := services.GetClusters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, clusters)
}

// GetCluster returns a specific cluster by ID
func GetCluster(c *gin.Context) {
	clusterID := c.Param("id")

	cluster, err := services.GetCluster(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if cluster == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cluster not found",
		})
		return
	}

	c.JSON(http.StatusOK, cluster)
}

// GetClusterNodes returns nodes for a specific cluster
func GetClusterNodes(c *gin.Context) {
	clusterID := c.Param("id")

	nodes, err := services.GetClusterNodes(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nodes)
}

// GetClusterPods returns pods for a specific cluster
func GetClusterPods(c *gin.Context) {
	clusterID := c.Param("id")

	pods, err := services.GetClusterPods(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, pods)
}

// GetClusterServices returns services for a specific cluster
func GetClusterServices(c *gin.Context) {
	clusterID := c.Param("id")

	services, err := services.GetClusterServices(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, services)
}

// GetClusterDeployments returns deployments for a specific cluster
func GetClusterDeployments(c *gin.Context) {
	clusterID := c.Param("id")

	deployments, err := services.GetClusterDeployments(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, deployments)
}

// GetClusterNamespaces returns namespaces for a specific cluster
func GetClusterNamespaces(c *gin.Context) {
	clusterID := c.Param("id")

	namespaces, err := services.GetClusterNamespaces(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, namespaces)
}
