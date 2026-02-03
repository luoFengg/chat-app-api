package web

// CreateConversationRequest for Creating New Conversation
type CreateConversationRequest struct {
	Type string `json:"type" binding:"required,oneof=direct group"`
	Name *string `json:"name"` // Mandatory for Group
	ParticipantIDs []string `json:"participant_ids" binding:"required,min=1"`
}

// UpdateConversationRequest for Updating Conversation (rename group)
type UpdateConversationRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

// AddParticipantRequest for Adding Participant to Conversation
type AddParticipantRequest struct {
	UserIDs []string `json:"user_ids" binding:"required,min=1"`
}
