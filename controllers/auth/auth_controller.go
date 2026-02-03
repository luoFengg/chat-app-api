package auth

import "github.com/gin-gonic/gin"

// AuthController interface for Authentication Controller
type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	GetMe(ctx *gin.Context)
}