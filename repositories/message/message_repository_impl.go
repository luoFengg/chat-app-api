package message

import (
	"chatapp-api/models/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// messageRepositoryImpl is an implementation of MessageRepository interface
type messageRepositoryImpl struct {
	db *gorm.DB
}

// NewMessageRepository makes a new MessageRepository instance
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepositoryImpl{db: db}
}

// Create implements MessageRepository
func (repo *messageRepositoryImpl) Create(ctx context.Context, message *domain.Message) error {
	return repo.db.WithContext(ctx).Create(message).Error
}

// FindByID implements MessageRepository
func (repo *messageRepositoryImpl) FindByID(ctx context.Context, id string) (*domain.Message, error) {
	var message domain.Message
	err := repo.db.WithContext(ctx).Preload("Sender").
	Where("id = ?", id).First(&message).Error

	if err != nil {
		return nil, err
	}

	return &message, nil
}

// FindByConversationID implements MessageRepository
func (repo *messageRepositoryImpl) FindByConversationID(ctx context.Context, conversationID string, limit, offset int) ([]domain.Message, error) {
	var messages []domain.Message

	err := repo.db.WithContext(ctx).
	Preload("Sender").
	Where("conversation_id = ?", conversationID).
	Order("created_at DESC").
	Limit(limit).
	Offset(offset).
	Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

// FindByConversationIDWithCursor implements MessageRepository
func (repo *messageRepositoryImpl) FindByConversationIDWithCursor(ctx context.Context, conversationID string, cursor *time.Time, limit int) ([]domain.Message, error) {
	var messages []domain.Message

	query := repo.db.WithContext(ctx).
	Preload("Sender").
	Where("conversation_id = ?", conversationID).
	Order("created_at DESC").
	Limit(limit)

	// if cursor not nil, filter messages older than cursor
	if cursor != nil {
		query = query.Where("created_at < ?", cursor)
	}

	err := query.Find(&messages).Error
	if err != nil {
		return nil, err
	}
	
	return messages, nil
}

// Update implements MessageRepository
func (repo *messageRepositoryImpl) Update(ctx context.Context, message *domain.Message) error {
	return repo.db.WithContext(ctx).Save(message).Error
}

// Delete implements MessageRepository
func (repo *messageRepositoryImpl) Delete(ctx context.Context, id string) error {
	return repo.db.WithContext(ctx).Delete(&domain.Message{}, "id = ?", id).Error
}

// CountByConversationID implements MessageRepository
func (repo *messageRepositoryImpl) CountByConversationID(ctx context.Context, conversationID string) (int64, error) {
	var count int64
	
	err := repo.db.WithContext(ctx).
	Model(&domain.Message{}).
	Where("conversation_id = ?", conversationID).
	Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

// FindLastByConversationID implements MessageRepository
func (repo *messageRepositoryImpl) FindLastByConversationID(ctx context.Context, conversationID string) (*domain.Message, error) {
	var message domain.Message

	err := repo.db.WithContext(ctx).
	Preload("Sender").
	Where("conversation_id = ?", conversationID).
	Order("created_at DESC").
	First(&message).Error

	if err != nil {
		return nil, err
	}

	return &message, nil
}