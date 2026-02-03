package domain

import (
	"chatapp-api/utils"
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	ID        string         `gorm:"type:varchar(32);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	AvatarURL *string        `gorm:"type:varchar(255)" json:"avatar_url,omitempty"`
	IsOnline  bool           `gorm:"default:false" json:"is_online"`
	LastSeen  *time.Time     `json:"last_seen,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relations
	Devices      []Device      `gorm:"foreignKey:UserID" json:"devices,omitempty"`
	Participants []Participant `gorm:"foreignKey:UserID" json:"-"`
}

// TableName to define table name
func (user *User) TableName() string {
	return "users"
}

// BeforeCreate hook to generate ID
func (user *User) BeforeCreate(tx *gorm.DB) error {
	if user.ID == "" {
		user.ID = utils.GenerateID("user")
	}
	return nil
}