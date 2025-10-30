package recovery

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Recover returns a Gin middleware that recovers from panics and logs them
func Recover() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic with request ID if available
		requestID := c.GetString("RequestID")
		if requestID == "" {
			requestID = "no-request-id"
		}

		// Log the panic details
		log.Printf("[%s] PANIC: %v", requestID, recovered)

		// Abort the request with a 500 status
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		// Continue to let other middleware complete
		c.Next()
	})
}
