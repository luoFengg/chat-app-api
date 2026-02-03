package web

// SendMessageRequest for Sending Message
type SendMessageRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	Content string `json:"content" binding:"required,min=1"`
	Type           string `json:"type" binding:"omitempty,oneof=text image video audio file location"`
}

// UpdateMessageRequest for Updating Message (Optional)
type UpdateMessageRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}
