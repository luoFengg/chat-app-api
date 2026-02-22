package main

import (
	"chatapp-api/apps/database"
	"chatapp-api/apps/redis"
	"chatapp-api/config"
	"chatapp-api/routes"
	"chatapp-api/websocket"
	"log"

	authController "chatapp-api/controllers/auth"
	conversationController "chatapp-api/controllers/conversation"
	messageController "chatapp-api/controllers/message"
	conversationRepo "chatapp-api/repositories/conversation"
	messageRepo "chatapp-api/repositories/message"
	userRepo "chatapp-api/repositories/user"
	authService "chatapp-api/services/auth"
	conversationService "chatapp-api/services/conversation"
	messageService "chatapp-api/services/message"
)

func main() {
	// 1. Load Configuration
	config := config.LoadConfig()
	log.Printf("Starting %s on port %s...", config.App.Name, config.App.Port)
	
	// 2. Connect to Database
	db := database.ConnectDatabase(config)

	// 3. Connect to Redis
	redisClient := redis.ConnectRedis(config)
	_ = redisClient

	// 4. Initialize WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// 5. Initialize repositories
	userRepository := userRepo.NewUserRepository(db)
	conversationRepository := conversationRepo.NewConversationRepository(db)
	messageRepository := messageRepo.NewMessageRepository(db)

	// 6. Initialize services
	authService := authService.NewAuthService(userRepository, config)
	conversationService := conversationService.NewConversationService(conversationRepository, userRepository)
	messageService := messageService.NewMessageService(messageRepository, conversationRepository)

	// 7. Initialize controllers
	authController := authController.NewAuthController(authService)
	conversationController := conversationController.NewConversationController(conversationService)
	messageController := messageController.NewMessageController(messageService)

	// 8. Setup router
	router := routes.SetupRouter(config, authController, conversationController, messageController, hub)

	// 9. Start server
	log.Printf("⏳ Attempting to start server on port %s...", config.App.Port)
	if err := router.Run(":" + config.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}