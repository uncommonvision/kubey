package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Health is a simple liveness-probe endpoint.
// It returns a JSON payload with a status flag and the current Unix timestamp.
func Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}
