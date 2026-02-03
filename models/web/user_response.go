package web

import "time"

// UserResponse for User Data Sent to the Client (with no password)
type UserResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	IsOnline bool `json:"is_online"`
	LastSeen *time.Time `json:"last_seen,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserBriefResponse for the quick view (e.g in list participant, sender, etc)
type UserBriefResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	IsOnline bool `json:"is_online"`
}