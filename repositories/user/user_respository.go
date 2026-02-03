package user

import (
	"chatapp-api/models/domain"
	"context"
)

// UserRepository interface for User Repository
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error
	
	// FindByID finds a user by ID
	FindByID(ctx context.Context, id string) (*domain.User, error)

	// FindByEmail finds a user by email (for duplicate check)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update updates a user's data
	Update(ctx context.Context, user *domain.User) error

	// Delete deletes a user
	Delete(ctx context.Context, id string) error

	// UpdateOnlineStatus updates a user's online status
	UpdateOnlineStatus(ctx context.Context, id string, isOnline bool) error
}