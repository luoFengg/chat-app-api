package conversation

import "github.com/gin-gonic/gin"

// ConversationController interface for conversation HTTP handlers
type ConversationController interface {
	// CreateConversation POST /conversation
	CreateConversation(ctx *gin.Context)

	// GetConversations GET /conversation
	GetConversations(ctx *gin.Context)

	// GetConversationByID GET /conversation/:id
	GetConversationByID(ctx *gin.Context)

	// UpdateConversation PUT /conversation/:id
	UpdateConversation(ctx *gin.Context)

	// AddParticipants POST /conversation/:id/participants
	AddParticipants(ctx *gin.Context)
	
	// LeaveConversation handles DELETE /conversations/:id/leave
	LeaveConversation(ctx *gin.Context)
	
	// KickParticipant handles DELETE /conversations/:id/participants/:userId
	KickParticipant(ctx *gin.Context)
}