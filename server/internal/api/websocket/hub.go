package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// Event represents a WebSocket event
type Event struct {
	Type      string      `json:"type"`
	ClusterID string      `json:"clusterId,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("WebSocket client connected. Total clients: %d", len(h.clients))

			// Send welcome message
			welcomeEvent := Event{
				Type:      "connection_status",
				Data:      map[string]string{"status": "connected"},
				Timestamp: time.Now(),
			}
			h.sendToClient(client, welcomeEvent)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			clientCount := len(h.clients)
			h.mu.RUnlock()

			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()

			if clientCount > 0 {
				log.Printf("Broadcast message to %d clients", clientCount)
			}
		}
	}
}

// Broadcast sends an event to all connected clients
func (h *Hub) Broadcast(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling WebSocket event: %v", err)
		return
	}

	select {
	case h.broadcast <- data:
	default:
		log.Printf("WebSocket broadcast channel is full, dropping message")
	}
}

// sendToClient sends an event to a specific client
func (h *Hub) sendToClient(client *Client, event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling WebSocket event: %v", err)
		return
	}

	select {
	case client.send <- data:
	default:
		log.Printf("Client send channel is full, dropping message")
	}
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// BroadcastClusterUpdate sends a cluster-specific update
func (h *Hub) BroadcastClusterUpdate(clusterID string, updateType string, data interface{}) {
	event := Event{
		Type:      updateType,
		ClusterID: clusterID,
		Data:      data,
		Timestamp: time.Now(),
	}
	h.Broadcast(event)
}
