package domain

import (
	"chatapp-api/utils"
	"time"

	"gorm.io/gorm"
)

// Device Model
type Device struct {
	ID         string    `gorm:"type:varchar(32);primaryKey" json:"id"`
	UserID     string    `gorm:"type:varchar(32);not null" json:"user_id"`
	FCMToken   string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"fcm_token"`
	DeviceType string    `gorm:"type:varchar(20);default:'android'" json:"device_type"`
	DeviceName *string   `gorm:"type:varchar(100)" json:"device_name,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	// Relations
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName to define table name
func (device *Device) TableName() string {
	return "devices"
}

// BeforeCreate hook to generate ID
func (device *Device) BeforeCreate(tx *gorm.DB) error {
	if device.ID == "" {
		device.ID = utils.GenerateID("dev")
	}
	return nil
}

