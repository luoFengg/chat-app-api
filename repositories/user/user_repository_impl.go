package user

import (
	"chatapp-api/models/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// userRepositoryImpl is an implementation of UserRepository interface
type userRepositoryImpl struct {
	db *gorm.DB
}

// NewuSerRepository make a New UserRepository Instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create implements UserRepository.
func (userRepo *userRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	return userRepo.db.WithContext(ctx).Create(user).Error
}

// FindByID implements UserRepository
func (userRepo *userRepositoryImpl) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := userRepo.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail implements UserRepository
func (userRepo *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := userRepo.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update implements UserRepository
func (userRepo *userRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	return userRepo.db.WithContext(ctx).Save(user).Error
}

// Delete implements UserRepository
func (userRepo *userRepositoryImpl) Delete(ctx context.Context, id string) error {
	return userRepo.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error
}

// UpdateOnlineStatus implements UserRepository
func (userRepo *userRepositoryImpl) UpdateOnlineStatus(ctx context.Context, id string, isOnline bool) error {
	updates := map[string]interface{}{
		"is_online": isOnline,
	}

	if !isOnline {
		updates["last_seen"] = time.Now()
	}

	return userRepo.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(updates).Error
}