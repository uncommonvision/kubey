package logging

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a Gin middleware for logging HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request ID if available
		requestID := c.GetString("RequestID")
		if requestID == "" {
			requestID = "no-request-id"
		}

		// Log request details
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		log.Printf("[%s] %s %s %d %v - %s",
			requestID,
			method,
			path,
			status,
			latency,
			clientIP,
		)
	}
}
