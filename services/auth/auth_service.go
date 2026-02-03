package auth

import (
	"chatapp-api/models/domain"
	"chatapp-api/models/web"
	"context"
)

// AuthService interface for Authentication Service
type AuthService interface {
	Register(ctx context.Context, req *web.RegisterRequest) (*web.AuthResponse, error)
	Login(ctx context.Context, req *web.LoginRequest) (*web.AuthResponse, error)
	RefreshToken(ctx context.Context, req *web.RefreshTokenRequest) (*web.TokenResponse, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}