package conversation

import (
	"chatapp-api/models/domain"
	"context"

	"gorm.io/gorm"
)

// conversationRepositoryImpl is an implementation of ConversationRepository interface
type conversationRepositoryImpl struct {
	db *gorm.DB
}

// NewConversationRepository makes a new ConversationRepository instance
func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepositoryImpl{db: db}
}

// Create implements ConversationRepository
func (repo *conversationRepositoryImpl) Create(ctx context.Context, conv *domain.Conversation) error {
	return repo.db.WithContext(ctx).Create(conv).Error
}

// FindByID implements ConversationRepository
func (repo *conversationRepositoryImpl) FindByID(ctx context.Context, id string) (*domain.Conversation, error) {
	var conv domain.Conversation
	err := repo.db.WithContext(ctx).
	Preload("Participants").
	Preload("Participants.User").
	Where("id = ?", id).
	First(&conv).Error

	if err != nil {
		return nil, err
	}

	return &conv, nil
}

// FindByUserID implements ConversationRepository
func (repo *conversationRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]domain.Conversation, error) {
	var conversations []domain.Conversation

	err := repo.db.WithContext(ctx).
	Joins("JOIN participants ON participants.conversation_id = conversations.id").
	Where("participants.user_id = ?", userID).
	Preload("Participants").
	Preload("Participants.User").
	Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(1)
	}).
	Order("updated_at DESC").
	Find(&conversations).Error

	if err != nil {
		return nil, err
	}

	return conversations, nil
}

// FindDirectConversation implements ConversationRepository
func (repo *conversationRepositoryImpl) FindDirectConversation(ctx context.Context, userID1, userID2 string) (*domain.Conversation, error) {
	var conv domain.Conversation

	err := repo.db.WithContext(ctx).
	Where("type = ?", "direct").
	Joins("JOIN participants p1 ON p1.conversation_id = conversations.id AND p1.user_id = ?", userID1).
	Joins("JOIN participants p2 ON p2.conversation_id = conversations.id AND p2.user_id = ?", userID2).
	First(&conv).Error

	if err != nil {
		return nil, err
	}

	return &conv, nil

}

// Update implements ConversationRepository
func (repo *conversationRepositoryImpl) Update(ctx context.Context, conv *domain.Conversation) error {
	return repo.db.WithContext(ctx).Save(conv).Error
}

// Delete implements ConversationRepository
func (repo *conversationRepositoryImpl) Delete(ctx context.Context, id string) error {
	return repo.db.WithContext(ctx).Delete(&domain.Conversation{}, "id = ?", id).Error
}

// AddParticipant implements ConversationRepository
func (repo *conversationRepositoryImpl) AddParticipant(ctx context.Context, participant *domain.Participant) error {
	return repo.db.WithContext(ctx).Create(participant).Error
}

// RemoveParticipant removes a participant from a conversation
func (repo *conversationRepositoryImpl) RemoveParticipant(ctx context.Context, conversationID, userID string) error {
    return repo.db.WithContext(ctx).
        Where("conversation_id = ? AND user_id = ?", conversationID, userID).
        Delete(&domain.Participant{}).Error
}