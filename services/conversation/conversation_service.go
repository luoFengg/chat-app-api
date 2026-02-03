package conversation

import (
	"chatapp-api/models/web"
	"context"
)

// ConversationService interface for conversation operations
type ConversationService interface {
	// CreateDirectConversation creates or retrieve existing DM between two users
	CreateConversation(ctx context.Context, userID string, req *web.CreateConversationRequest) (*web.ConversationResponse, error)

	// GetConversations retrieves all conversations for a user (list view)
	GetConversations(ctx context.Context, userID string) ([]web.ConversationListItem, error)

	// GetConversationByID retrieves a single conversation detail
	GetConversationByID(ctx context.Context, userID, conversationID string) (*web.ConversationResponse, error)

	// UpdateConversation updates conversation name (group only)
	UpdateConversation(ctx context.Context, userID, conversationID string, req *web.UpdateConversationRequest) (*web.ConversationResponse, error)

	// AddParticipants adds new participants to a group conversation
	AddParticipants(ctx context.Context, userID, conversationID string, req *web.AddParticipantRequest) error

	// LeaveConversation leaves a conversation
	LeaveConversation(ctx context.Context, userID, conversationID string) error

	// KickParticipant removes a participant from a conversation (amdin only)
	KickParticipant(ctx context.Context, adminUserID, conversationID, targetUserID string) error
}