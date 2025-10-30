package request

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID returns a Gin middleware that generates and tracks request IDs
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID header is already present
		requestID := c.GetHeader("X-Request-ID")

		// If not present, generate a new UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store request ID in Gin context
		c.Set("RequestID", requestID)

		// Add request ID to response header
		c.Header("X-Request-ID", requestID)

		// Continue processing
		c.Next()
	}
}
