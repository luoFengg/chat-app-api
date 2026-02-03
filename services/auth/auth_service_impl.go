package auth

import (
	"chatapp-api/config"
	"chatapp-api/exceptions"
	"chatapp-api/models/domain"
	"chatapp-api/models/web"
	"chatapp-api/utils"
	"context"
	"errors"

	userRepo "chatapp-api/repositories/user"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// authServiceImpl implements AuthService interface
type authServiceImpl struct {
	userRepo userRepo.UserRepository
	config  *config.Config
}

// NewAuthService Create new Instance of AuthService
func NewAuthService(userRepo userRepo.UserRepository, config *config.Config) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
		config:  config,
	}
}

// Register to register new user
func (authService *authServiceImpl) Register(ctx context.Context, req *web.RegisterRequest) (*web.AuthResponse, error) {
	// 1. Check if email already exists
	existingUser, err := authService.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingUser != nil {
		return nil, exceptions.NewConflictError("Email already registered")
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Create user
	user := &domain.User{
		Name: req.Name,
		Email: req.Email,
		Password: string(hashedPassword),
	}

	if err := authService.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 4. Generate tokens and return response
	return authService.generateAuthResponse(user)
}

// Login to Authenticate User
func (authService *authServiceImpl) Login(ctx context.Context, req *web.LoginRequest) (*web.AuthResponse, error) {
	// 1. Find user by email
	user, err := authService.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewUnauthorizedError("Invalid email or password")
		}
		return nil, err
	}

	// 2. Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, exceptions.NewUnauthorizedError("Invalid email or password")
	}

	// 3. Generate tokens and return response
	return authService.generateAuthResponse(user)
}

// RefreshToken to Create New Access Token from Refresh Token
func (authService *authServiceImpl) RefreshToken(ctx context.Context, req *web.RefreshTokenRequest) (*web.TokenResponse, error) {
	// 1. Validate refresh token
	claims, err := utils.ValidateToken(req.RefreshToken, authService.config.JWT.RefreshSecret)
	if err != nil {
		return nil, exceptions.NewUnauthorizedError("Invalid or expired refresh token")
	}	

	// 2. Generate new access token
	accessToken, expiresAt, err := utils.GenerateAccessToken(claims.UserID, authService.config.JWT.Secret)
	if err != nil {
		return nil, err
	}

	return &web.TokenResponse{
		AccessToken: accessToken,
		ExpiresAt: expiresAt,
	}, nil
}

// GetUserByID to Get User by ID
func (authService *authServiceImpl) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	// 1. Find user by ID
	user, err := authService.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError("User not found")
		}
		return nil, err
	}

	return user, nil
}

// generateAuthResponse is a helper function to create response with tokens
func (authService *authServiceImpl) generateAuthResponse(user *domain.User) (*web.AuthResponse, error) {
	// Generate access token
	accessToken, expiresAt, err := utils.GenerateAccessToken(user.ID, authService.config.JWT.Secret)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID, authService.config.JWT.RefreshSecret)
	if err != nil {
		return nil, err
	}

	// Return response
	return &web.AuthResponse{
		User: web.UserResponse{
			ID: user.ID,
			Name: user.Name,
			Email: user.Email,
			AvatarURL: user.AvatarURL,
			IsOnline: user.IsOnline,
			LastSeen: user.LastSeen,
		},
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		ExpiresAt: expiresAt,
	}, nil
}