package web

// SendMessageRequest for Sending Message
type SendMessageRequest struct {
	ConversationID string  `json:"-"` // Set from URL param, not from body
	Content        string  `json:"content"`
	Caption        *string `json:"caption,omitempty"`
	Type           string  `json:"type" binding:"omitempty,oneof=text image file"`
}

// UpdateMessageRequest for Updating Message (Optional)
type UpdateMessageRequest struct {
	Content *string `json:"content,omitempty" binding:"omitempty,min=1"`
	Caption *string `json:"caption,omitempty"`
}
