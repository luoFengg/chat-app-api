package message

import (
	"chatapp-api/middleware"
	"chatapp-api/models/web"
	messageService "chatapp-api/services/message"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// messageControllerImpl implements MessageController interface
type messageControllerImpl struct {
	messageService messageService.MessageService
}

// NewMessageController Create new instance of MessageController
func NewMessageController(messageService messageService.MessageService) MessageController {
	return &messageControllerImpl{
		messageService: messageService,
	}
}

// SendMessage handles POST /conversations/:conversationId/messages
func (controller *messageControllerImpl) SendMessage(ctx *gin.Context) {
	// 1. Get userID from context (from JWT middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}	

	// 2. Get conversation ID from URL parameter
	conversationID := ctx.Param("id")

	// 3. Bind request body
	var req web.SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error: err.Error(),
		})
		return
	}

	// 4. Set conversation ID from URL
	req.ConversationID = conversationID

	// 5. Call service to send message
	message, err := controller.messageService.SendMessage(ctx.Request.Context(), userID, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 6. Return response
	ctx.JSON(http.StatusCreated, web.ApiResponse{
		Success: true,
		Message: "Message sent successfully",
		Data:    message,
	})
}

// GetMessages handles GET /conversation/:conversationId/messages
func (controller *messageControllerImpl) GetMessages(ctx *gin.Context) {
	// 1. Get user ID from context
	userID, err := middleware.GetUserIDFromContext(ctx)

	// 2. Get conversation ID from URL parameter
	conversationID := ctx.Param("id")

	// 3. Parse query parameters (cursor & limit)
	var cursor *time.Time
	if cursorStr := ctx.Query("cursor"); cursorStr != "" {
		parsedCursor, err := time.Parse(time.RFC3339Nano, cursorStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.ErrorResponse{
				Success: false,
				Message: "Invalid cursor format",
				Error: err.Error(),
			})
			return
		}
		cursor = &parsedCursor
	}

	limit := 20 // default value
	if limitStr := ctx.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// 4. Call service
	messages, cursorMeta, err := controller.messageService.GetMessages(ctx.Request.Context(), userID, conversationID, cursor, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Return response
	ctx.JSON(http.StatusOK, web.CursorResponse{
		Success: true,
		Message: "Messages fetched successfully",
		Data: messages,
		Cursor: *cursorMeta,
	})
}

// GetMessageByID handles GET /messages/:messageId
func (controller *messageControllerImpl) GetMessageByID(ctx *gin.Context) {
	// 1. Get userID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 2. Get message ID from URL parameter
	messageID := ctx.Param("messageId")

	// 3. Call service
	message, err := controller.messageService.GetMessageByID(ctx.Request.Context(), userID, messageID)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Return response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "Message fetched successfully",
		Data: message,
	})
}

// UpdateMessage handles PUT /messages/:messageId
func (controller *messageControllerImpl) UpdateMessage(ctx *gin.Context) {
	// 1. Get userID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 2. Get message ID from URL parameter
	messageID := ctx.Param("messageId")

	// 3. Bind request body
	var req web.UpdateMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error: err.Error(),
		})
		return
	}

	// 4. Call service to update message
	message, err := controller.messageService.UpdateMessage(ctx.Request.Context(), userID, messageID, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Return response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "Message updated successfully",
		Data: message,
	})
}

// DeleteMessage handles DELETE /messages/:messageId
func (controller *messageControllerImpl) DeleteMessage(ctx *gin.Context) {
	// 1. Get userID from context
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	// 2. Get message ID from URL parameter
	messageID := ctx.Param("messageId")
	// 3. Call service (soft delete)
	err = controller.messageService.DeleteMessage(ctx.Request.Context(), userID, messageID)
	if err != nil {
		ctx.Error(err)
		return
	}
	// 4. Return response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "Message deleted successfully",
	})
}