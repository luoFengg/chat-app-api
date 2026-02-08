package message

import (
	"chatapp-api/models/domain"
	"context"
	"time"
)

// MessageRepository inferface for message operations
type MessageRepository interface {
	// Create creates a new message
	Create(ctx context.Context, message *domain.Message) error

	// FindByID finds a message by ID
	FindByID(ctx context.Context, id string) (*domain.Message, error)

	// FindByConversationID finds all messages in a conversation
	FindByConversationID(ctx context.Context, conversationID string, limit, offset int) ([]domain.Message, error)

	// FindByConversationIDWithCursor finds messages using cursor-based pagination
	FindByConversationIDWithCursor(ctx context.Context, conversationID string, cursor *time.Time, limit int) ([]domain.Message, error)
	// Update updates a message
	Update(ctx context.Context, message *domain.Message) error

	// Delete deletes a message
	Delete(ctx context.Context, id string) error

	// CountByConversationID counts messages in a conversation
	CountByConversationID(ctx context.Context, conversationID string) (int64, error)

	// FindLastByConversationID finds the last message in a conversation
	FindLastByConversationID(ctx context.Context, conversationID string) (*domain.Message, error)
}