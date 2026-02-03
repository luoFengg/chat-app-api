package auth

import (
	"chatapp-api/middleware"
	"chatapp-api/models/web"
	authService "chatapp-api/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// authControllerImpl implements AuthController interface
type authControllerImpl struct {
	authService authService.AuthService
}

// NewAuthController Create new instance of AuthController
func NewAuthController(authService authService.AuthService) AuthController {
	return &authControllerImpl{
		authService: authService,
	}
}

// Register Handles POST /auth/register
func (controller *authControllerImpl) Register(ctx *gin.Context) {
	var req web.RegisterRequest
	
	// 1. Bind JSON body to struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, web.ApiResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// 2. Call service to register user
	result, err := controller.authService.Register(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err) // Middleware handle it
		return	
	}

	// 3. Return success response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "User registered successfully",
		Data: result,
	})
}

// Login Handles POST /auth/login
func (controller *authControllerImpl) Login(ctx *gin.Context) {
	var req web.LoginRequest

	// 1. Bind JSON body to struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error: err.Error(),
		})
		return
	}

	// 2. Call service to login user
	result, err := controller.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err) // Middleware handle it
		return	
	}

	// 3. Return success response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "User logged in successfully",
		Data: result,
	})
}

// RefreshToken Handles POST /auth/refresh-token
func (controller *authControllerImpl) RefreshToken(ctx *gin.Context) {
	var req web.RefreshTokenRequest

	// 1. Bind JSON body to struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, web.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
			Error: err.Error(),
		})
		return
	}

	// 2. Call service to refresh token
	result, err := controller.authService.RefreshToken(ctx.Request.Context(), &req)
	if err != nil {
		ctx.Error(err) // Middleware handle it
		return	
	}

	// 3. Return success response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: result,
	})
}

// GetMe Handles GET /auth/me
func (controller *authControllerImpl) GetMe(ctx *gin.Context) {
	// 1. Get user ID from context (set by middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 2. Get user from service
	user, err := controller.authService.GetUserByID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 3. Return user response
	ctx.JSON(http.StatusOK, web.ApiResponse{
		Success: true,
		Message: "User profile retrieved",
		Data: web.UserResponse{
			ID: user.ID,
			Name: user.Name,
			Email: user.Email,
			AvatarURL: user.AvatarURL,
			IsOnline: user.IsOnline,
			LastSeen: user.LastSeen,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

