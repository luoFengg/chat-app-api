package message

import "github.com/gin-gonic/gin"

// MessageController interface for message HTTP handlers
type MessageController interface {
	// SendMessage POST /conversations/:conversationId/messages
	SendMessage(ctx *gin.Context)

	// GetMessages GET /conversations/:conversationId/messages
	GetMessages(ctx *gin.Context)

	// GetMessageByID GET /messages/:messageId
	GetMessageByID(ctx *gin.Context)

	// UpdateMessage PUT /messages/:messageId
	UpdateMessage(ctx *gin.Context)

	// DeleteMessage DELETE /messages/:messageId
	DeleteMessage(ctx *gin.Context)
}