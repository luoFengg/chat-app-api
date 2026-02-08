package web

import "time"

// MessageResponse for Full Detail Message
type MessageResponse struct {
	ID             string            `json:"id"`
	ConversationID string            `json:"conversation_id"`
	Sender         UserBriefResponse `json:"sender"`
	Content        string            `json:"content"`
	Caption        *string           `json:"caption,omitempty"`
	Type           string            `json:"type"`
	IsEdited       bool              `json:"is_edited"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

// MessageBriefResponse for Quick View Message (e.g in list conversation)
type MessageBriefResponse struct {
	ID        string            `json:"id"`
	Sender    UserBriefResponse `json:"sender"`
	Content   string            `json:"content"`
	Caption   *string           `json:"caption,omitempty"`
	Type      string            `json:"type"`
	CreatedAt time.Time         `json:"created_at"`
}