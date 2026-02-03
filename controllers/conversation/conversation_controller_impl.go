package conversation

import (
	"chatapp-api/middleware"
	"chatapp-api/models/web"
	convService "chatapp-api/services/conversation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type conversationControllerImpl struct {
	convService convService.ConversationService
}

func NewConversationController(convService convService.ConversationService) ConversationController {
	return &conversationControllerImpl{
		convService: convService,
	}
}

// CreateConversation handles POST /conversations
func (controller *conversationControllerImpl) CreateConversation(ctx *gin.Context) {
	// 1. Get userID from context (from JWT middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}	

	// 2. Bind request body
	var req web.CreateConversationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error: err.Error(),
		})
	}

	// 3. Call service to create conversation
	result, err := controller.convService.CreateConversation(ctx.Request.Context(), userID, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Reutrn response 
	ctx.JSON(http.StatusCreated, web.ApiResponse{
		Success: true,
		Message: "Conversation created successfully",
		Data: result,
	})
}

// GetConversations handles GET /conversations
func (controller *conversationControllerImpl) GetConversations(ctx *gin.Context) {
	// 1. Get User ID from context (from JWT middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 2. Call service to get conversations
	result, err := controller.convService.GetConversations(ctx.Request.Context(), userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 3. Return response 
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "Conversations fetched successfully",
		Data: result,
	})
}

// GetConversationByID handles GET /conversations/:id
func (controller *conversationControllerImpl) GetConversationByID(ctx *gin.Context) {
	// 1. Get User ID from context (from JWT middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	} 

	// 2. Get conversation ID from param
	conversationID := ctx.Param("id")

	result, err := controller.convService.GetConversationByID(ctx.Request.Context(), userID, conversationID)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 3. Return response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "Conversation fetched successfully",
		Data: result,
	})
}

// UpdateConversation handles PUT /conversations/:id
func (controller *conversationControllerImpl) UpdateConversation(ctx *gin.Context) {
	// 1. Get User ID from context (from JWT middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 2. Get conversation ID from param
	conversationID := ctx.Param("id")

	// 3. Bind request body
	var req web.UpdateConversationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrorResponse{
			Success: false,
			Message: "Invalid body request",
			Error: err.Error(),
		})
	}

	// 4. Call service to update conversation
	result, err := controller.convService.UpdateConversation(ctx.Request.Context(), userID, conversationID, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Return response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "Conversation updated successfully",
		Data: result,
	})
}

// AddParticipants handles POST /conversations/:id/participants
func (controller *conversationControllerImpl) AddParticipants(ctx *gin.Context) {
	// 1. Get User ID from context (from JWT middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 2. Get conversation ID from param
	conversationID := ctx.Param("id")

	// 3. Bind request body
	var req web.AddParticipantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrorResponse{
			Success: false,
			Message: "Invalid body request",
			Error: err.Error(),
		})
	}

	// 4. Call service to add participants
	err = controller.convService.AddParticipants(ctx.Request.Context(), userID, conversationID, &req)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 5. Return response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "Participants added succesfully",
	})
}

// LeaveConversation handles DELETE /conversations/:id/leave
func (controller *conversationControllerImpl) LeaveConversation(ctx *gin.Context) {
	// 1. Get User ID from context (from JWT middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 2. Get Conversation ID from param
	conversationID := ctx.Param("id")

	// 3. Call service to leave conversation
	err = controller.convService.LeaveConversation(ctx.Request.Context(), userID, conversationID)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Return response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "Left conversation successfully",
	})
}

// KickParticipant handles DELETE /conversations/:id/participants/:userId
func (controller *conversationControllerImpl) KickParticipant(ctx *gin.Context) {
	// 1. Get admin user ID from context
	adminUserID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 2. Get conversation ID and target user ID from param
	conversationID := ctx.Param("id")
	targetUserID := ctx.Param("userId")

	// 3. Call service to kick participant
	err = controller.convService.KickParticipant(ctx.Request.Context(), adminUserID, conversationID, targetUserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 4. Return response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "Participant kicked successfully",
	})
}