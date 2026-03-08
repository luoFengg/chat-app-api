package web

// RegisterRequest for New Register Body Request
type RegisterRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateProfileRequest for Updating Profile Body Request
type UpdateProfileRequest struct {
	Name      *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}


// LoginRequest for Login Body Request
type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest for Refresh Token Body Request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
