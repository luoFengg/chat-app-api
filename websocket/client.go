package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Maximum time for waiting message from client
	pongWait = 60 * time.Second

	// Send ping interval to client (must < pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size (bytes)
	maxMessageSize = 512

	// Buffer size for send channel
	sendBufferSize = 256
)

// Client represents a WebSocket connection
type Client struct {
	hub *Hub
	conn *websocket.Conn
	userID string
	send chan []byte
}

// NewClient creates a new Client instance
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub: hub,
		conn: conn,
		userID: userID,
		send: make(chan []byte, sendBufferSize),
	}
}

// ReadPum reads messages from the WebSocket connection
func (client *Client) ReadPump() {
	defer func() {
		client.hub.unregister <- client
		client.conn.Close()
	}()

	// Set read limit and deadline
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Start reading messages
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, 
			websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket error for user %s:%v", client.userID, err)
			}
			break
		}

		log.Printf("Received message from user %s: %s", client.userID, string(message))
		// TODO: Handle incoming messages (typing indicator, read receipt, etc.)
		// Parse JSON message from client
		var wsMessage WSMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			log.Printf("Invalid message format from user %s: %v", client.userID, err)
			continue // Skip message with incorrect formatting, continue reading the next message
		}
		
		// Handle based on event type
		switch wsMessage.Event{
		case EventTypingStart, EventTypingStop:
			// Forward typing events to other participants
			client.hub.HandleTypingEvent(client.userID, wsMessage.ConversationID, wsMessage.Event)

		case EventMessageRead:
		// Client sends: {"event":"message_read", "conversation_id":"conv_xxx", "data":{"message_id":"msg_xxx"}}
		dataMap, ok := wsMessage.Data.(map[string]interface{})
		if !ok {
			log.Printf("Invalid message_read data from user %s", client.userID)
			break
		}
		messageID, ok := dataMap["message_id"].(string)
		if !ok {
			log.Printf("Missing message_id in message_read from user %s", client.userID)
			break
		}
		client.hub.HandleMessageReadEvent(client.userID, wsMessage.ConversationID, messageID)
		case EventMessageDelivered:
		// Client sends: {"event":"message_delivered", "conversation_id":"conv_xxx", "data":{"message_id":"msg_xxx"}}
		dataMap, ok := wsMessage.Data.(map[string]interface{})
		if !ok {
			log.Printf("Invalid message_delivered data from user %s", client.userID)
			break
		}
		messageID, ok := dataMap["message_id"].(string)
		if !ok {
			log.Printf("Missing message_id in message_delivered from user %s", client.userID)
			break
		}
		client.hub.HandleMessageDeliveredEvent(client.userID, wsMessage.ConversationID, messageID)

		default:
			log.Printf("Unknown event from user %s: %s", client.userID, wsMessage.Event)
		}
	}
}

// WritePump writes messages to the WebSocket connection
func (client *Client) WritePump() {
	// Create ticker for sending ping every 54 seconds
	ticker := time.NewTicker(pingPeriod)
	
	// Cleanup when function ends
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	// Main loop: listen for messages and ping ticker
	for {
		select {
			// Case 1: There is a message to send to client
		case message, ok := <- client.send:
			// If channel is closed by Hub (unregister), send close message and stop
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Write message to WebSocket connection
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}

			// Case 2: Ping ticker fired, send ping to check if client is alive
		case <- ticker.C:
			if err := client.conn.WriteMessage(websocket.PingMessage, nil);
			err != nil {
				return
			}
		}
	}
}