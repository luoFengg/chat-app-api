package message_receipt

import (
	"chatapp-api/models/domain"
	"context"
)

// MessageReceiptRepository interface: defines the contract for message receipt database operations
type MessageReceiptRepository interface {
	// Create a new receipt (called when a message is sent)
	Create(ctx context.Context, receipt *domain.MessageReceipt) error

	// CreateBatch creates multiple receipts at once (1 message -> multiple recipients)
	CreateBatch(ctx context.Context, receipts []*domain.MessageReceipt) error

	// UpdateStatus updates the receipt status (sent -> delivered -> read)
	UpdateStatus(ctx context.Context, messsageID, userID, status string) error
	
	// FindByMessageID returns all receipts for a spesific message
	// Useful for: "siapa aja yang udah baca pesan ini?"
	FindByMessageID(ctx context.Context, messageID string) ([]domain.MessageReceipt, error)

	
	// FindByMessageAndUser returns a single receipt for a specific message + user combo
	// Useful for: "apakah User B sudah baca pesan msg_001?"
	FindByMessageAndUser(ctx context.Context, messageID string, userID string) (*domain.MessageReceipt, error)
}