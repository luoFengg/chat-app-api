package websocket

// Event types for WebSocket communication
const (
	// Server to client events
	EventNewMessage = "new_message"
	EventUserOnline = "user_online"
	EventUserOffline = "user_offline"

	// Client to server events (and forwarded to other clients)
	EventTypingStart = "typing_start"
	EventTypingStop = "typing_stop"
	EventMessageRead = "message_read"
	EventMessageDelivered = "message_delivered"
)

// WSMessage is the standard WebSocket message format
// Every message sent/received through WebSocket follows this format
type WSMessage struct {
	Event string `json:"event"`
	ConversationID string `json:"conversation_id,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// TypingData is the payload for typing events
type TypingData struct {
	UserID string `json:"user_id"`
	DisplayName string `json:"display_name"`
}

// OnlineStatusData is the payload for online/offline events
type OnlineStatusData struct {
	UserID string `json:"user_id"`
}

// MessageReadData is the payload for message read events
type MessageReadData struct {
	UserID string `json:"user_id"`
	MessageID string `json:"message_id"`
}
