package clusters

import (
	"net/http"

	"kubey/api/internal/services/kubernetes"

	"github.com/gin-gonic/gin"
)

// GetClusters returns all clusters
func GetClusters(c *gin.Context) {
	clusters, err := kubernetes.GetClusters()
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

	cluster, err := kubernetes.GetCluster(clusterID)
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

	nodes, err := kubernetes.GetClusterNodes(clusterID)
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

	pods, err := kubernetes.GetClusterPods(clusterID)
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

	services, err := kubernetes.GetClusterServices(clusterID)
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

	deployments, err := kubernetes.GetClusterDeployments(clusterID)
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

	namespaces, err := kubernetes.GetClusterNamespaces(clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, namespaces)
}
