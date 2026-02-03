package web

import "time"

// AuthResponse for successful Login/Register Response
type AuthResponse struct {
	User UserResponse `json:"user"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// TokenResponse for Refresh Token Response
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt time.Time `json:"expires_at"`
}