package message_receipt

import (
	"chatapp-api/models/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// messageReceiptRepositoryImpl implements MessageReceiptRepository
type messageReceiptRepositoryImpl struct {
	db *gorm.DB
}

// NewMessageReceiptRepository creates a new message receipt repository
func NewMessageReceiptRepository(db *gorm.DB) MessageReceiptRepository {
	return &messageReceiptRepositoryImpl{db: db}
}

// Create saves a single receipt to database
func (repo *messageReceiptRepositoryImpl) Create(ctx context.Context, receipt *domain.MessageReceipt) error {
	return repo.db.WithContext(ctx).Create(receipt).Error
}

// CreateBatch saves multiple receipts to database
func (repo *messageReceiptRepositoryImpl) CreateBatch(ctx context.Context, receipts []*domain.MessageReceipt) error {
	return repo.db.WithContext(ctx).Create(&receipts).Error
}

// UpdateStatus updates receipt status and sets the corresponding timestamp
func (repo *messageReceiptRepositoryImpl) UpdateStatus(ctx context.Context, messageID, userID, status string) error {
	// Build update data based on status
	updates := map[string]interface{}{
		"status": status,
	}

	// Set timestamp based on status
	now := time.Now()
	switch status {
	case "delivered":
		// Message arriced at recipient's device
		updates["delivered_at"] = now
	case "read":
		// Recipient opened and read the message
		updates["delivered_at"] = now // Also mark as delivered if not already
		updates["read_at"] = now
	}

	return repo.db.WithContext(ctx).
	Model(&domain.MessageReceipt{}).
	Where("message_id = ? AND user_id = ?", messageID, userID).
	Updates(updates).Error
}

// FindByMessageID returns all receipts for a message (with user info preloaded)
func (repo *messageReceiptRepositoryImpl) FindByMessageID(ctx context.Context, messageID string) ([]domain.MessageReceipt, error) {
	var receipts []domain.MessageReceipt

	err := repo.db.WithContext(ctx).
	Preload("User").
	Where("message_id = ?", messageID).
	Find(&receipts).Error

	if err != nil {
		return nil, err
	}

	return receipts, nil
}

// FindByMessageAndUser returns a single receipt for a specific message + user
func (repo *messageReceiptRepositoryImpl) FindByMessageAndUser(ctx context.Context, messageID string, userID string) (*domain.MessageReceipt, error) {
	var receipt domain.MessageReceipt
	err := repo.db.WithContext(ctx).
		Where("message_id = ? AND user_id = ?", messageID, userID).
		First(&receipt).Error
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}