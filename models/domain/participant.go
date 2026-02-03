package domain

import (
	"chatapp-api/utils"
	"time"

	"gorm.io/gorm"
)

// Participant model
type Participant struct {
	ID             string    `gorm:"type:varchar(32);primaryKey" json:"id"`
	UserID         string    `gorm:"type:varchar(32);not null" json:"user_id"`
	ConversationID string    `gorm:"type:varchar(32);not null" json:"conversation_id"`
	Role           string    `gorm:"type:varchar(20);default:'member'" json:"role"`
	JoinedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"joined_at"`
	
	// Relations
	User         User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Conversation Conversation `gorm:"foreignKey:ConversationID" json:"-"`
}

// TableName to define table name
func (participant *Participant) TableName() string {
	return "participants"
}

// BeforeCreate hook to generate ID
func (participant *Participant) BeforeCreate(tx *gorm.DB) error {
	if participant.ID == "" {
		participant.ID = utils.GenerateID("part")
	}
	return nil
}
