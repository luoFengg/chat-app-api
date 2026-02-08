package routes

import (
	"chatapp-api/config"
	"chatapp-api/controllers/auth"
	"chatapp-api/controllers/conversation"
	"chatapp-api/controllers/message"
	"chatapp-api/exceptions"
	"chatapp-api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter
func SetupRouter(
	config *config.Config, 
	authController auth.AuthController, 
	convController conversation.ConversationController,
	messageController message.MessageController,) *gin.Engine {
	// Create router
	router := gin.Default()

	// Global middleware
	router.Use(exceptions.ErrorHandler())

	// API v1
	    v1 := router.Group("/api/v1")
    {
        // Auth routes
        authRoutes := v1.Group("/auth")
        {
            // Public
            authRoutes.POST("/register", authController.Register)
            authRoutes.POST("/login", authController.Login)
            authRoutes.POST("/refresh", authController.RefreshToken)
            // Protected (butuh token)
            authRoutes.Use(middleware.AuthMiddleware(config))
            authRoutes.GET("/me", authController.GetMe)
        }

        // Conversation routes
        conversationRoutes := v1.Group("/conversations")
        conversationRoutes.Use(middleware.AuthMiddleware(config))
        {
            conversationRoutes.POST("", convController.CreateConversation)
			conversationRoutes.GET("", convController.GetConversations)
			conversationRoutes.GET("/:id", convController.GetConversationByID)
			conversationRoutes.PUT("/:id", convController.UpdateConversation)
			conversationRoutes.POST("/:id/participants", convController.AddParticipants)
			conversationRoutes.DELETE("/:id/leave", convController.LeaveConversation)
			conversationRoutes.DELETE("/:id/participants/:userId", convController.KickParticipant)
        
			// Message routes (nested under conversations)
			//POST & GET messages wihtin a conversation
			conversationRoutes.POST("/:id/messages", messageController.SendMessage)
			conversationRoutes.GET("/:id/messages", messageController.GetMessages)
		}

		// Message routes (direct - for single message opeations)
		messageRoutes := v1.Group("/messages")
		messageRoutes.Use(middleware.AuthMiddleware(config))
		{
			messageRoutes.GET("/:messageId", messageController.GetMessageByID)
			messageRoutes.PUT("/:messageId", messageController.UpdateMessage)
			messageRoutes.DELETE("/:messageId", messageController.DeleteMessage)
		}
    }
	return router
}