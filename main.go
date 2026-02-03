package main

import (
	"chatapp-api/apps/database"
	"chatapp-api/config"
	"chatapp-api/routes"
	"log"

	authController "chatapp-api/controllers/auth"
	conversationController "chatapp-api/controllers/conversation"
	conversationRepo "chatapp-api/repositories/conversation"
	userRepo "chatapp-api/repositories/user"
	authService "chatapp-api/services/auth"
	conversationService "chatapp-api/services/conversation"
)

func main() {
	// 1. Load Configuration
	config := config.LoadConfig()
	log.Printf("Starting %s on port %s...", config.App.Name, config.App.Port)
	
	// 2. Connect to Database
	db := database.ConnectDatabase(config)

	// 3. Initialize repositories
	userRepository := userRepo.NewUserRepository(db)
	conversationRepository := conversationRepo.NewConversationRepository(db)

	// 4. Initialize services
	authService := authService.NewAuthService(userRepository, config)
	conversationService := conversationService.NewConversationService(conversationRepository, userRepository)


	// 5. Initialize controllers
	authController := authController.NewAuthController(authService)
	conversationController := conversationController.NewConversationController(conversationService)

	// 6. Setup router
	router := routes.SetupRouter(config, authController, conversationController)

	// 7. Start server
	log.Printf("⏳ Attempting to start server on port %s...", config.App.Port)
	if err := router.Run(":" + config.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}