package domain

import (
	"chatapp-api/utils"
	"time"

	"gorm.io/gorm"
)

// MessageReceipt tracks the delivery and read status of a message per recipient
// Each message has one receipt per recipient (sender doesn't get a receipt)
type MessageReceipt struct {
	// Unique ID for this receipt (rcpt_xxx)
	ID string `gorm:"type:varchar(32);primaryKey" json:"id"`

	// Which message this receipt belongs to (FK to messages)
	MessageID string `gorm:"type:varchar(32);not null" json:"message_id"`

	// Which user received this message (FK to users)
	UserID string `gorm:"type:varchar(32);not null" json:"user_id"`

	// Current status: "sent", "delivered", or "read"
	Status string `gorm:"type:varchar(20);not null;default:'sent'" json:"status"`

	// When the message was delivered to the recipient's device (null = not yet)
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`

	// When the recipient read the message (null = not yet)
	ReadAt *time.Time `json:"read_at,omitempty"`

	// Relations
	Message Message `gorm:"foreignKey:MessageID" json:"-"`
	User    User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName defines the table name in database
func (receipt *MessageReceipt) TableName() string {
	return "message_receipts"
}

// BeforeCreate hook to auto-generate ID with "rcpt_" prefix
func (receipt *MessageReceipt) BeforeCreate(tx *gorm.DB) error {
	if receipt.ID == "" {
		receipt.ID = utils.GenerateID("rcpt")
	}
	return nil
}