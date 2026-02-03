package middleware

import (
	"chatapp-api/config"
	"chatapp-api/models/web"
	"chatapp-api/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Context keys to store user data
const (
	ContextKeyUserID = "userID"
)

// AuthMiddleware to validate JWT token
func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get Authorization hedar
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, 
			web.ApiResponse{
				Success: false,
				Message: "Authorization header is required",
			})
			return
		}

		// 2. Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
			web.ApiResponse{
				Success: false,
				Message: "Invalid Authorization header format. User: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]

		// 3. Validate token
		claims, err := utils.ValidateToken(tokenString, config.JWT.Secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, 
			web.ApiResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			return
		}

		// 4. Set user ID to context
		c.Set(ContextKeyUserID, claims.UserID)

		// 5. Continue to next handler
		c.Next()
	}
}

// GetUserIDFromContext to retrieve user ID from context
func GetUserIDFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return "", errors.New("user ID not found in context")
	}
	return userID.(string), nil
}