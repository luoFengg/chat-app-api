package websocket

import (
	"log"
	"sync"
)

// Hub manages all active websocket connections
type Hub struct {
	// Registered clients, key = userID
	clients map[string]*Client

	// Channel to register new client
	register chan *Client

	// Channel to unregister client
	unregister chan *Client

	// Mutex to protect clients map
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (hub *Hub) Run() {
	for {
		select {
		case client := <- hub.register:
			hub.mu.Lock()
			hub.clients[client.userID] = client
			hub.mu.Unlock()
			log.Printf("User %s connected via WebSocket", client.userID)

		case client := <- hub.unregister:
			hub.mu.Lock()
			if _, ok := hub.clients[client.userID]; ok {
				delete(hub.clients, client.userID)
				close(client.send)
			}
			hub.mu.Unlock()
			log.Printf("User %s disconnected from WebSocket", client.userID)
		}
	}
}

// GetClient returns a client by userID
func (hub *Hub) GetClient(userID string) (*Client, bool) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	client, ok := hub.clients[userID]
	return client, ok
}

// SendToUser sends a message to a spesific user
func (hub *Hub) SendToUser(userID string, message []byte) {
	hub.mu.RLock()
	client, ok := hub.clients[userID]
	hub.mu.RUnlock()

	if ok {
		select {
		case client.send <- message:
		default:
			// Channel is full, may be slow/stuck
			log.Printf("Failed to send message to user %s: channel full", userID)
		}
	}
}

// SendToUsers sends a message to multiple users (group)
func (hub *Hub) SendToUsers(userIDs []string, message []byte) {
	for _, userID := range userIDs {
		hub.SendToUser(userID, message)
	}
}