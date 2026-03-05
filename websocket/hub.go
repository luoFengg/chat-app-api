package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	conversationRepo "chatapp-api/repositories/conversation"
	receiptRepo "chatapp-api/repositories/message_receipt"
	userRepo "chatapp-api/repositories/user"
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

	// Repository to get conversation participants (for typing indicator)
	conversationRepo conversationRepo.ConversationRepository

	// Repository to update online status in DB
	userRepo userRepo.UserRepository

	// Repository to update message receipts (delivered/read)
	receiptRepo receiptRepo.MessageReceiptRepository
}

// NewHub creates a new Hub instance
func NewHub(conversationRepo conversationRepo.ConversationRepository, 
	userRepo userRepo.UserRepository,
	receiptRepo receiptRepo.MessageReceiptRepository,
) *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		register: make(chan *Client),
		unregister: make(chan *Client),
		conversationRepo: conversationRepo,
		userRepo: userRepo,
		receiptRepo: receiptRepo,
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

			// Update online status in DB
			if err := hub.userRepo.UpdateOnlineStatus(context.Background(), client.userID, true); err != nil {
				log.Printf("Failed to update online status for user %s: %v", client.userID, err)
			}

			// Broadcast to all online user that this user is now online
			hub.broadcastOnlineStatus(client.userID, EventUserOnline)

		case client := <- hub.unregister:
			hub.mu.Lock()
			if _, ok := hub.clients[client.userID]; ok {
				delete(hub.clients, client.userID)
				close(client.send)
			}
			hub.mu.Unlock()
			log.Printf("User %s disconnected from WebSocket", client.userID)
		
			// Update database: user offline + save last_seen (automatic in repo)
			if err := hub.userRepo.UpdateOnlineStatus(context.Background(), client.userID, false); err != nil {
				log.Printf("Failed to update offline status for user %s: %v", client.userID, err)
			}

			// Broadcast to all online user that this user is now offline or disconnected
			hub.broadcastOnlineStatus(client.userID, EventUserOffline)
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

// HandleTypingEvent forwards typing events to other participants in the conversation
func (hub *Hub) HandleTypingEvent(senderID, conversationID, event string) {
	// 1. Search conversation to get participants
	conv, err := hub.conversationRepo.FindByID(context.Background(), conversationID)
	if err != nil {
		log.Printf("Failed to find conversation %s: %v", conversationID, err)
		return
	}

	// 2. Create a typing message in WSMessage format
	wsMessage := WSMessage{
		Event: event, // "typing_start" or "typing_stop"
		ConversationID: conversationID,
		Data: TypingData{
			UserID: senderID,
		},
	}

	// 3. Convert to JSON bytes
	jsonData, err := json.Marshal(wsMessage)
	if err != nil {
		log.Printf("Failed to marshal typing event: %v", err)
		return
	}

	// 4. Send to all participants except the sender
	// User A doesn't need to see their own typing indicator
	for _, participant := range conv.Participants {
		if participant.UserID != senderID {
			hub.SendToUser(participant.UserID, jsonData)
		}
	}
}

// broadcastOnlineStatus sends online/offline status to all connected users
func (hub *Hub) broadcastOnlineStatus(userID, event string) {
	// 1. Create message in WSMessage format
	wsMessage := WSMessage{
		Event: event, // "user_online" or "user_offline"
		Data: OnlineStatusData{
			UserID: userID, // Which user changed status
		},
	}

	// 2. Convert to JSON bytes
	jsonData, err := json.Marshal(wsMessage)
	if err != nil {
		log.Printf("Failed to marshal online status: %v", err)
		return
	}

	// 3. Send to all online user (except the user who changed status)
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	for id, client := range hub.clients {
		if id != userID {
			select {
			case client.send <- jsonData:
			default:
				log.Printf("Failed to send online status to user %s: channel full", id)
			}
		}
	}
}

// HandleMessageReadEvent processes a "message_read" event from a client
// Updates the receipt status in DB and notifies other participants
func (hub *Hub) HandleMessageReadEvent(readerID, conversationID, messageID string) {
	// 1. Update receipt status in database: "sent"/"delivered" → "read"
	if err := hub.receiptRepo.UpdateStatus(context.Background(), messageID, readerID, "read"); err != nil {
		log.Printf("Failed to update read status for message %s by user %s: %v", messageID, readerID, err)
		return
	}
	
	// 2. Build read notification in WSMessage format
	wsMsg := WSMessage{
		Event:          EventMessageRead,
		ConversationID: conversationID,
		Data: MessageReadData{
			UserID:    readerID,
			MessageID: messageID,
		},
	}
	
	// 3. Convert to JSON
	jsonData, err := json.Marshal(wsMsg)
	if err != nil {
		log.Printf("Failed to marshal message read event: %v", err)
		return
	}
	
	// 4. Find conversation to get participants
	conv, err := hub.conversationRepo.FindByID(context.Background(), conversationID)
	if err != nil {
		log.Printf("Failed to find conversation %s: %v", conversationID, err)
		return
	}

	// 5. Notify all participants except the reader
	for _, p := range conv.Participants {
		if p.UserID != readerID {
			hub.SendToUser(p.UserID, jsonData)
		}
	}
}

// HandleMessageDeliveredEvent processes a "message_delivered" event from a client
// Updates the receipt status in DB and notifies other participants
func (hub *Hub) HandleMessageDeliveredEvent(recipientID, conversationID, messageID string) {
	// 1. Update receipt status in database: "sent" → "delivered"
	if err := hub.receiptRepo.UpdateStatus(context.Background(), messageID, recipientID, "delivered"); err != nil {
		log.Printf("Failed to update delivery status for message %s by user %s: %v", messageID, recipientID, err)
		return
	}

	// 2. Build WebSocket message to notify other participants
	wsMsg := WSMessage{
		Event:          EventMessageDelivered,
		ConversationID: conversationID,
		Data: MessageReadData{
			UserID:    recipientID,
			MessageID: messageID,
		},
	}

	// 3. Convert to JSON
	jsonData, err := json.Marshal(wsMsg)
	if err != nil {
		log.Printf("Failed to marshal message delivered event: %v", err)
		return
	}

	// 4. Find conversation to get participants
	conv, err := hub.conversationRepo.FindByID(context.Background(), conversationID)
	if err != nil {
		log.Printf("Failed to find conversation %s: %v", conversationID, err)
		return
	}

	// 5. Notify all participants except the recipient
	for _, p := range conv.Participants {
		if p.UserID != recipientID {
			hub.SendToUser(p.UserID, jsonData)
		}
	}
}