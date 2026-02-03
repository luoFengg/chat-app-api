package domain

import (
	"chatapp-api/utils"
	"time"

	"gorm.io/gorm"
)

// Conversation model
type Conversation struct {
	ID        string         `gorm:"type:varchar(32);primaryKey" json:"id"`
	Name      *string        `gorm:"type:varchar(100)" json:"name,omitempty"`
	Type      string         `gorm:"type:varchar(20);not null;default:'direct'" json:"type"`
	CreatedBy string         `gorm:"type:varchar(32);not null" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relations
	Participants []Participant `gorm:"foreignKey:ConversationID" json:"participants,omitempty"`
	Messages     []Message     `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

// TableName to define table name
func (conversation *Conversation) TableName() string {
	return "conversations"
}

// BeforeCreate hook to generate ID
func (conversation *Conversation) BeforeCreate(tx *gorm.DB) error {
	if conversation.ID == "" {
		conversation.ID = utils.GenerateID("conv")
	}
	return nil
}
