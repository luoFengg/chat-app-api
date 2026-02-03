package web

import "time"

// ConversationResponse for Detail Conversation
type ConversationResponse struct {
	ID string `json:"id"`
	Name *string `json:"name,omitempty"`
	Type string `json:"type"`
	DisplayName string `json:"display_name"` // Name of Group or Opposite User Name
	DisplayAvatar *string `json:"display_avatar,omitempty"` // Avatar of Group or Opposite User
	CreatedBy string `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Participants []ParticipantResponse `json:"participants,omitempty"`
	LastMessage *MessageBriefResponse `json:"last_message,omitempty"`
}

// ConversationListItem for List Conversation (more compact)
type ConversationListItem struct {
	ID string `json:"id"`
	Type string `json:"type"`
	DisplayName string `json:"display_name"`
	DisplayAvatar *string `json:"display_avatar,omitempty"`
	LastMessage *MessageBriefResponse `json:"last_message,omitempty"`
	UnreadCount int `json:"unread_count"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ParticipantResponse for Participant Data
type ParticipantResponse struct {
	User UserBriefResponse `json:"user"`
	Role string `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}