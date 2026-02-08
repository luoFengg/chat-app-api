package message

import (
	"chatapp-api/models/domain"
	"chatapp-api/models/web"
	"context"
	"time"
)

// MessageService interface for message business logic
type MessageService interface {
	// SendMessage sends a new message to a conversation
	SendMessage(ctx context.Context, senderID string, req *web.SendMessageRequest) (*domain.Message, error)

	// GetMessages gets messages with cursor-based pagination
	GetMessages(ctx context.Context, userID, conversationID string, cursor *time.Time, limit int) ([]domain.Message, *web.CursorMeta, error)

	// GetMessageByID gets a single message by ID
	GetMessageByID(ctx context.Context, userID, messageID string) (*domain.Message, error)

	// UpdateMessage updates a message
	UpdateMessage(ctx context.Context, userID, messageID string, req *web.UpdateMessageRequest) (*domain.Message, error)

	// DeleteMessage deletes a message
	DeleteMessage(ctx context.Context, userID, messageID string) error
}