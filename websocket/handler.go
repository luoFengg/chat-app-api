package websocket

import (
	"chatapp-api/config"
	"chatapp-api/models/web"
	"chatapp-api/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connection to WebSocket
var upgrader = gorilla.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	// CheckOrigin allows connections from any origin (for development)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleWebSocket handles WebSocket connetction requests
func HandleWebSocket(hub *Hub, config *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Get token from query parameter
		tokenString := ctx.Query("token")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, web.ApiResponse{
				Success: false,
				Message: "Token is required. Use ?token=<jwt_token>",
			})
			return
		}

		// 2. Validate JWT token
		claims, err := utils.ValidateToken(tokenString, config.JWT.Secret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, web.ApiResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			return
		}

		// 3. Upgrade HTTP connection to WebSocket
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Printf("Failed to upgrade connection for user %s: %v", claims.UserID, err)
			return
		}

		// 4. Create new client and register to hub
		client := NewClient(hub, conn, claims.UserID)
		hub.register <- client

		// 5. Start client goroutines
		go client.WritePump()
		go client.ReadPump()

		log.Printf("WebSocket connection established for user %s", claims.UserID)
	}
}