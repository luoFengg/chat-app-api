package domain

import (
	"chatapp-api/utils"
	"time"

	"gorm.io/gorm"
)

// Message model
type Message struct {
	ID             string         `gorm:"type:varchar(32);primaryKey" json:"id"`
	ConversationID string         `gorm:"type:varchar(32);not null" json:"conversation_id"`
	SenderID       string         `gorm:"type:varchar(32);not null" json:"sender_id"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	Caption        *string        `gorm:"type:text" json:"caption,omitempty"`
	Type           string         `gorm:"type:varchar(20);default:'text'" json:"type"`
	IsEdited       bool           `gorm:"type:boolean;default:false" json:"is_edited"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relations
	Sender       User         `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Conversation Conversation `gorm:"foreignKey:ConversationID" json:"-"`
}

// TableName to define table name
func (message *Message) TableName() string {
	return "messages"
}

// BeforeCreate hook to generate ID
func (message *Message) BeforeCreate(tx *gorm.DB) error {
	if message.ID == "" {
		message.ID = utils.GenerateID("msg")
	}
	return nil
}

