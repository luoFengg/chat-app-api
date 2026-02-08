package message

import (
	"chatapp-api/exceptions"
	"chatapp-api/models/domain"
	"chatapp-api/models/web"
	conversationRepo "chatapp-api/repositories/conversation"
	messageRepo "chatapp-api/repositories/message"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

// messageServiceImpl implements MessageService interface
type messageServiceImpl struct {
	messageRepo        messageRepo.MessageRepository
	conversationRepo   conversationRepo.ConversationRepository
}

// NewMessageService creates a new MessageService instance
func NewMessageService(messageRepo messageRepo.MessageRepository, conversationRepo conversationRepo.ConversationRepository) MessageService {
	return &messageServiceImpl{
		messageRepo: messageRepo, 
		conversationRepo: conversationRepo}
}

// SendMessage implements MessageService interface
func (service *messageServiceImpl) SendMessage(ctx context.Context, senderID string, req *web.SendMessageRequest) (*domain.Message, error) {
	// 1. Validate: Check if conversation exists
	conv, err := service.conversationRepo.FindByID(ctx, req.ConversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError("conversation not found")
		}
		return nil, err
	}

	// 2. Validate: Check if sender is a participant
	isParticipant := false
	for _, participant := range conv.Participants {
		if participant.UserID == senderID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, exceptions.NewForbiddenError("You are not a participant in this conversation")
	}

	// 3. Validate message type
	validTypes := map[string]bool{
		"text": true,
		"image": true,
		"file": true,
	}

	messageType := req.Type
	if messageType == "" {
		messageType = "text"
	}

	if !validTypes[messageType] {
		return nil, exceptions.NewBadRequestError("Invalid message type. Allowed: text, image, file")
	}

	// 4. Validate content based on type
	if messageType == "text" && req.Content == "" {
		return nil, exceptions.NewBadRequestError("Message content is required for text type")
	}
	
	// 5. Create message
	message := &domain.Message{
		ConversationID: req.ConversationID,
		SenderID: senderID,
		Content: req.Content,
		Caption: req.Caption,
		Type: messageType,
	}

	// 6. Save to database
	if err := service.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	// 7. Reload message with sender info
	return service.messageRepo.FindByID(ctx, message.ID)
}

// GetMessages implements MessageService
func (service *messageServiceImpl) GetMessages(ctx context.Context, userID, conversationID string, cursor *time.Time, limit int) ([]domain.Message, *web.CursorMeta, error) {
	// 1. Validate: Check if conversation exists
	conv, err := service.conversationRepo.FindByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil , exceptions.NewNotFoundError("Conversation not found")
		}
		return nil, nil, err
	}

	// 2. Validate: Check if user is a participant
	isParticipant := false
	for _, participant := range conv.Participants {
		if participant.UserID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, nil, exceptions.NewForbiddenError("You are not a participant in this conversation")
	}

	// 3. Set default limit
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	// 4. Fetch messages (limit + 1 to check if has more)
	messages, err := service.messageRepo.FindByConversationIDWithCursor(ctx, conversationID, cursor, limit+1)
	if err != nil {
		return nil, nil, err
	}

	// 5. Build cursor metadata
	hasMore := len(messages) > limit
	var nextCursor *string

	if hasMore {
		// Remove the extra message
		messages = messages[:limit]
		lastCreatedAt := messages[len(messages)-1].CreatedAt.Format(time.RFC3339Nano)
		nextCursor = &lastCreatedAt
	}

	cursorMeta := &web.CursorMeta{
		HasMore: hasMore,
		NextCursor: nextCursor,
	}

	return messages, cursorMeta, nil
}

// GetMessageByID implements MessageService
func (service *messageServiceImpl) GetMessageByID(ctx context.Context, userID, messageID string) (*domain.Message, error) {
	// 1. Find the message
	message, err := service.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError("Message not found")
		}
		return nil, err
	}

	// 2. Validate: Check if user is a participant in this conversation
	conv, err := service.conversationRepo.FindByID(ctx, message.ConversationID)
	if err != nil {
		return nil, err
	}

	isParticipant := false
	for _, participant := range conv.Participants {
		if participant.UserID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, exceptions.NewForbiddenError("You are not a participant in this conversation")
	}

	// 3. Return the message
	return message, nil
}

// UpdateMessage implements MessageService
func (service *messageServiceImpl) UpdateMessage(ctx context.Context, userID, messageID string, req *web.UpdateMessageRequest) (*domain.Message, error) {
	// 1. Validate: At least one field must be provided
	if req.Content == nil && req.Caption == nil {
		return nil, exceptions.NewBadRequestError("At least one field (content or caption) must be provided")
	}

	// 2. Find the message
	message, err := service.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError("Message not found")
		}
		return nil, err
	}

	// 3. Validate: Only sender can edit their own message
	if message.SenderID != userID {
		return nil, exceptions.NewForbiddenError("You can only edit your own message")
	}

	// 4. Update based on message type
	if message.Type == "text" {
		// Text messages: can update content
		if req.Content != nil {
			message.Content = *req.Content
			message.IsEdited = true
		}
	} else {
		// Image/File messages: can only update caption, NOT content
		if req.Content != nil {
			return nil, exceptions.NewBadRequestError("Content can only be updated for text messages")
		}
		if req.Caption != nil {
			message.Caption = req.Caption
			message.IsEdited = true
		}
	}

	// 5. Save to database
	if err := service.messageRepo.Update(ctx, message); err != nil {
		return nil, err
	}

	// 6. Reload to get fresh data
	return service.messageRepo.FindByID(ctx, message.ID)
}

// DeleteMessage implements MessageService
func (service *messageServiceImpl) DeleteMessage(ctx context.Context, userID, messageID string) error {
	// 1. Find the message
	message, err := service.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exceptions.NewNotFoundError("Message not found")
		}
		return err
	}

	// 2. Validate: Only sender can delete their own message
	if message.SenderID != userID {
		return exceptions.NewForbiddenError("You can only delete your own message")
	}

	// 3. Delete the message
	return service.messageRepo.Delete(ctx, message.ID)
}