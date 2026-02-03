package conversation

import (
	"chatapp-api/exceptions"
	"chatapp-api/models/domain"
	"chatapp-api/models/web"
	convRepo "chatapp-api/repositories/conversation"
	userRepo "chatapp-api/repositories/user"
	"context"
	"errors"

	"gorm.io/gorm"
)

// conversationServiceImpl implements ConversationService interface
type conversationServiceImpl struct {
	convRepo convRepo.ConversationRepository	
	userRepo userRepo.UserRepository
}

// NewConversationService makes a new ConversationService instance
func NewConversationService(convRepo convRepo.ConversationRepository, userRepo userRepo.UserRepository) ConversationService {
	return &conversationServiceImpl{convRepo: convRepo,
		userRepo: userRepo, 
	}
}

// CreateConversation implements ConversationService
func (service *conversationServiceImpl) CreateConversation(ctx context.Context, userID string, req *web.CreateConversationRequest) (*web.ConversationResponse, error) {
	// 1. Validation: For DM, only 1 participant allowed (interlocutor)
	if req.Type == "direct" {
		if len(req.ParticipantIDs) != 1 {
			return nil, exceptions.NewBadRequestError("Direct message requires exactly 1 participant")
		}

		// Check if DM already exists between these two users
		targetUserID := req.ParticipantIDs[0]
		existingConv, err := service.convRepo.FindDirectConversation(ctx, userID, targetUserID)
		if err == nil && existingConv != nil {
			// DM already exists, return it
			return service.buildConversationResponse(existingConv, userID), nil
		}
	}

	// 2. Validation: For Group, name is required
	if req.Type == "group" {
		if req.Name == nil || *req.Name == "" {
			return nil, exceptions.NewBadRequestError("Group name is required")
		}
	}

	// 3. Validation: All participants must exist in database
	for _, participantID := range req.ParticipantIDs {
		_, err := service.userRepo.FindByID(ctx, participantID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, exceptions.NewNotFoundError("User " + participantID + " not found")
			}
			return nil, err
		}
	}

	// 4. Build participants list (include creator as admin)
	participants := make([]domain.Participant, 0, len(req.ParticipantIDs)+1)

	// Add creator as admin
	participants = append(participants, domain.Participant{
		UserID: userID,
		Role: "admin",
	})

	// Add other participants as member
	for _, participantID := range req.ParticipantIDs {
		// Skip if same as creator (prevent duplicate)
		if participantID == userID {
			continue
		}
		participants = append(participants, domain.Participant{
			UserID: participantID,
			Role: "member",
		})
	}

	// 5. Create conversation with participants
	conversation := &domain.Conversation{
		Type: req.Type,
		Name: req.Name,
		CreatedBy: userID,
		Participants: participants,
	}

	if err := service.convRepo.Create(ctx, conversation); err != nil {
		return nil, err
	}

	// 6. Reload conversation with complete relation for response
	createdConv, err := service.convRepo.FindByID(ctx, conversation.ID)
	if err != nil {
		return nil, err
	}

	return service.buildConversationResponse(createdConv, userID), nil
}

// GetConversations implements ConversationService
func (service *conversationServiceImpl) GetConversations(ctx context.Context, userID string) ([]web.ConversationListItem, error) {
	// 1. Get all conversations for the user
	conversations, err := service.convRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Convert every single conversation to ConversationListItem
	result := make([]web.ConversationListItem, len(conversations))
	for idx, conv := range conversations {
		result[idx] = service.buildConversationListItem(&conv, userID)
	}

	return result, nil
}

// GetConversationByID implements ConversationService
func (service *conversationServiceImpl) GetConversationByID(ctx context.Context, userID, conversationID string) (*web.ConversationResponse, error) {
	// 1. Find conversation by ID
	conversation, err := service.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError("Conversation not found")
		}
		return nil, err
	}

	// 2. Validation: Check if user is a participant
	isParticipant := false
	for _, participant := range conversation.Participants {
		if participant.UserID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, exceptions.NewForbiddenError("You are not a participant of this conversation")
	}

	// 3. Return conversation response
	return service.buildConversationResponse(conversation, userID), nil
}

// UpdateConversation implements ConversationService
func (service *conversationServiceImpl) UpdateConversation(ctx context.Context, userID, conversationID string, req *web.UpdateConversationRequest) (*web.ConversationResponse, error) {
	// 1. Find conversation
	conv, err := service.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError("Conversation not found")
		}
		return nil, err
	}	

	// 2. Validation: Only group can be updated (direct can't be updated)
	if conv.Type != "group" {
		return nil, exceptions.NewBadRequestError("Cannot update direct message conversation")
	}

	// 3. Validation: User must be an admin
	isAdmin := false
	for _, participant := range conv.Participants {
		if participant.UserID == userID && participant.Role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return nil, exceptions.NewForbiddenError("Only admin can update conversation")
	}

	// 4. Update conversation name
	conv.Name = &req.Name
	
	if err := service.convRepo.Update(ctx, conv); err != nil {
		return nil, err
	}

	// 5. Reload and return
	updatedConv, err := service.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	return service.buildConversationResponse(updatedConv, userID), nil
}

// AddParticipants implements ConversationService
func (service *conversationServiceImpl) AddParticipants(ctx context.Context, userID, conversationID string, req *web.AddParticipantRequest) error {
	// 1. Find conversation
	conv, err := service.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewNotFoundError("Conversation not found")
		}
		return err
	}

	// 2. Validation: Only group can add participants
	if conv.Type != "group" {
		return exceptions.NewBadRequestError("Cannot add participants to direct message")
	}

	// 3. Validation: User must be an admin
	isAdmin := false
	for _, participant := range conv.Participants {
		if participant.UserID == userID && participant.Role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return exceptions.NewForbiddenError("Only admin can add participants")
	}

	// 4. Create map existing participants for check duplication
	existingParticipants := make(map[string]bool)
	for _, participant := range conv.Participants {
		existingParticipants[participant.UserID] = true
	}

	// 5. Add new participants
	for _, newUserID :=  range req.UserIDs {
		// Skip if user is already a participant
		if existingParticipants[newUserID] {
			continue
		}

		// Validation: user must exist in database
		_, err := service.userRepo.FindByID(ctx, newUserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exceptions.NewNotFoundError("User " + newUserID + " not found")
			}
			return err
		}

		// Add as member
		newParticipant := &domain.Participant{
			UserID: newUserID,
			ConversationID: conversationID,
			Role: "member",
		}

		if err := service.convRepo.AddParticipant(ctx, newParticipant); err != nil {
			return err
		}
	}

	return nil
}

// LeaveConversation implements ConversationService
func (service *conversationServiceImpl) LeaveConversation(ctx context.Context, userID, conversationID string) error {
	// 1. Find conversation
	conv, err := service.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewNotFoundError("Conversation not found")
		}
		return err
	}

	// 2. Validaton: Cannot leave DM
	if conv.Type == "direct" {
		return exceptions.NewBadRequestError("Cannot leave direct message conversation")
	}

	// 3. Find participant user and check if user is a participant
	var userParticipant *domain.Participant
	for idx, participant := range conv.Participants {
		if participant.UserID == userID {
			userParticipant = &conv.Participants[idx]
			break
		}
	}

	if userParticipant == nil {
		return exceptions.NewForbiddenError("You are not a participant of this conversation")
	}

	// 4. If user is the only admin, promote another member become admin
	if userParticipant.Role == "admin" {
		// Count the number of admins
		adminCount := 0 
		for _, participant := range conv.Participants {
			if participant.Role == "admin" {
				adminCount++
			}
		}

		// If there is only one admin (this user), promote another member become admin
		if adminCount == 1 {
			for _, participant := range conv.Participants {
				if participant.UserID != userID && participant.Role == "member" {
					participant.Role = "admin"
					if err := service.convRepo.Update(ctx, conv); err != nil {
						return err
					}
					break
				}
			}
		}
	}

	// 5. Remove user from participant
	if err := service.convRepo.RemoveParticipant(ctx, conversationID, userID); err != nil {
		return err
	}

	return nil
}

// KickParticipant implements ConversationService
func (service *conversationServiceImpl) KickParticipant(ctx context.Context, adminUserID, conversationID, targetUserID string) error {
	// 1. Find conversation
	conv, err := service.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewNotFoundError("Conversation not found")
		}
		return err
	}

	// 2. Validation: Cannot kick at DM
	if conv.Type == "direct" {
		return exceptions.NewBadRequestError("Cannot kick participant from direct message")
	}

	// 3. Validation: User must be an admin
	isAdmin := false
	for _, participant := range conv.Participants {
		if participant.UserID == adminUserID && participant.Role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return exceptions.NewForbiddenError("Only admin can kick participants")
	}

	// 4. Validation: Cannot kick yourself (Use LeaveConversation instead)
	if targetUserID == adminUserID {
		return exceptions.NewBadRequestError("Cannot kick yourself, use leave instead")
	}

	// 5. Find the target and validate that he is a participant
	var targetParticipant *domain.Participant
	for idx, participant := range conv.Participants {
		if participant.UserID == targetUserID {
			targetParticipant = &conv.Participants[idx]
			break
		}
	}

	if targetParticipant == nil {
		return exceptions.NewNotFoundError("Target user is not a participant")
	}

	// 6. Validation: Cannot kick another admin
	if targetParticipant.Role == "admin" {
		return exceptions.NewConflictError("Cannot kick another admin")
	}

	// 7. Remove target from participant
	if err := service.convRepo.RemoveParticipant(ctx, conversationID, targetUserID); err != nil {
		return err
	}

	return nil
}


// buildConversationResponse converts domain.Conversation to web.ConversationResponse
func (s *conversationServiceImpl) buildConversationResponse(conv *domain.Conversation, currentUserID string) *web.ConversationResponse {
	// Initialize
	participants := make([]web.ParticipantResponse, len(conv.Participants))
	var displayName string
	var displayAvatar *string

	// LOOP: Process each participant
	for idx, participant := range conv.Participants {
		// Convert to DTO
		participants[idx] = web.ParticipantResponse{
			User: web.UserBriefResponse{
				ID:        participant.User.ID,
				Name:      participant.User.Name,
				AvatarURL: participant.User.AvatarURL,
				IsOnline:  participant.User.IsOnline,
			},
			Role:     participant.Role,
			JoinedAt: participant.JoinedAt,
		}

		// For DM: find the OTHER user (not yourself)
		if conv.Type == "direct" && participant.UserID != currentUserID {
			displayName = participant.User.Name
			displayAvatar = participant.User.AvatarURL
		}
	}
	// END OF LOOP

	// For Group: use conversation name 
	if conv.Type == "group" && conv.Name != nil {
		displayName = *conv.Name
	}

	// Build last message if any 
	var lastMessage *web.MessageBriefResponse
	if len(conv.Messages) > 0 {
		msg := conv.Messages[0]
		lastMessage = &web.MessageBriefResponse{
			ID:      msg.ID,
			Content: msg.Content,
			Type:    msg.Type,
			Sender: web.UserBriefResponse{
				ID:        msg.Sender.ID,
				Name:      msg.Sender.Name,
				AvatarURL: msg.Sender.AvatarURL,
				IsOnline:  msg.Sender.IsOnline,
			},
			CreatedAt: msg.CreatedAt,
		}
	}

	// Build final response
	return &web.ConversationResponse{
		ID:            conv.ID,
		Name:          conv.Name,
		Type:          conv.Type,
		DisplayName:   displayName,
		DisplayAvatar: displayAvatar,
		CreatedBy:     conv.CreatedBy,
		CreatedAt:     conv.CreatedAt,
		UpdatedAt:     conv.UpdatedAt,
		Participants:  participants,
		LastMessage:   lastMessage,
	}
}

// buildConversationListItem converts domain.Conversation to web.ConversationListItem
func (service *conversationServiceImpl) buildConversationListItem(conv *domain.Conversation, currentUserID string) web.ConversationListItem {
	var displayName string
	var displayAvatar *string

	// Determine DisplayName and DisplayAvatar based on conversation type
	if conv.Type == "direct" {
		// For DM: find the OTHER user (not yourself)
		for _, participant := range conv.Participants {
			if participant.UserID != currentUserID {
				displayName = participant.User.Name
				displayAvatar = participant.User.AvatarURL
				break
			}
		}
	} else if conv.Type == "group" && conv.Name != nil {
		// For Group: user conversation or group name
		displayName = *conv.Name
	}

	// Build last message (if any)
	var lastMessage *web.MessageBriefResponse
	if len(conv.Messages) > 0 {
		msg := conv.Messages[0]
		lastMessage = &web.MessageBriefResponse{
			ID: msg.ID,
			Content: msg.Content,
			Type: msg.Type,
			Sender: web.UserBriefResponse{
				ID: msg.Sender.ID,
				Name: msg.Sender.Name,
				AvatarURL: msg.Sender.AvatarURL,
				IsOnline: msg.Sender.IsOnline,
			},
			CreatedAt: msg.CreatedAt,
		}
	}

		return web.ConversationListItem{
		ID:            conv.ID,
		Type:          conv.Type,
		DisplayName:   displayName,
		DisplayAvatar: displayAvatar,
		LastMessage:   lastMessage,
		UnreadCount:   0, // TODO: implement unread count logic
		UpdatedAt:     conv.UpdatedAt,
	}
}