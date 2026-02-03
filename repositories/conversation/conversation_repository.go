package conversation

import (
	"chatapp-api/models/domain"
	"context"
)

// ConversationRepository interface for conversation operations
type ConversationRepository interface {
	// Create creates a new conversation
	Create(ctx context.Context, conv *domain.Conversation) error

	// FingByID finds a conversation by ID
	FindByID(ctx context.Context, id string) (*domain.Conversation, error)

	// FindByUserID finds all conversations for a user
	FindByUserID(ctx context.Context, userID string) ([]domain.Conversation, error)

	// FindDirectConversation finds existing DM between two users
	FindDirectConversation(ctx context.Context, userID1, userID2 string) (*domain.Conversation, error)

	// Update updates a conversation
	Update(ctx context.Context, conv *domain.Conversation) error

	// Delete deletes a conversation
	Delete(ctx context.Context, id string) error

	// Participant operations
	// AddParticipant adds new participant to conversation 
	AddParticipant(ctx context.Context, participant *domain.Participant) error

	// RemoveParticipant removes participant from conversation
	RemoveParticipant(ctx context.Context, conversationID, userID string) error
}