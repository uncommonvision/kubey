package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleWebSocket handles WebSocket upgrade requests
func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection to WebSocket: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
			return
		}

		client := NewClient(hub, conn)
		hub.register <- client

		// Start the client's read and write pumps
		client.Start()
	}
}

// HandleWebSocketStatus returns WebSocket connection status
func HandleWebSocketStatus(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"connected_clients": hub.GetClientCount(),
			"status":            "active",
		})
	}
}
