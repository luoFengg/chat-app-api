package message

import (
	"chatapp-api/exceptions"
	"chatapp-api/models/domain"
	"chatapp-api/models/web"
	conversationRepo "chatapp-api/repositories/conversation"
	messageRepo "chatapp-api/repositories/message"
	receiptRepo "chatapp-api/repositories/message_receipt"
	"chatapp-api/websocket"
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// messageServiceImpl implements MessageService interface
type messageServiceImpl struct {
	messageRepo        messageRepo.MessageRepository
	conversationRepo   conversationRepo.ConversationRepository
	receiptRepo        receiptRepo.MessageReceiptRepository
	hub                *websocket.Hub
}

// NewMessageService creates a new MessageService instance
func NewMessageService(
	messageRepo messageRepo.MessageRepository, 
	conversationRepo conversationRepo.ConversationRepository, 
	receiptRepo receiptRepo.MessageReceiptRepository, 
	hub *websocket.Hub) MessageService {
	return &messageServiceImpl{
		messageRepo: messageRepo, 
		conversationRepo: conversationRepo,
		receiptRepo: receiptRepo,
		hub: hub,
	}
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

	// 7. Create receipts for all recipients (everyone except sender)
	// Each recipient gets a receipt with initial status "sent"
	var receipts []*domain.MessageReceipt
	for _, participant := range conv.Participants {
		// Sender doesn't need a rececipt for their own message
		if participant.UserID != senderID {
			receipts = append(receipts, &domain.MessageReceipt{
				MessageID: message.ID,
				UserID: participant.UserID,
				Status: "sent",
			})
		}
	} 

	// Save all receipts in one batch query (for efficiency)
	if len(receipts) > 0 {
		if err := service.receiptRepo.CreateBatch(ctx, receipts); err != nil {
			// Log error but don't fail the whole operation
			// Receipt creation failure shouldn't block message delivery
			log.Printf("Failed to create message receipts %s: %v", message.ID, err)
		}
	}

	// 8. Reload message with sender info
	// Retrieved the saved message complete with sender data
	savedMessage, err := service.messageRepo.FindByID(ctx, message.ID)
	if err != nil {
		return nil, err
	}

	// 8. Broadcast new message to all participants via WebSocket
	// Send realtime notification to all online participants
	if service.hub != nil {
		// 8a. Wrap the message in standard WSMessage format
		wsMessage := websocket.WSMessage{
			Event: websocket.EventNewMessage, // Event type: "new_message"
			ConversationID: savedMessage.ConversationID, // Destionation conversation ID
			Data: savedMessage, // Complete message data
		}

		// 8b. Convert struct to JSON bytes to send via WebSocket
		jsonData, err := json.Marshal(wsMessage)
		if err == nil {
			// 8c. Collcet all participants for this conversation
			var participantIDs []string
			for _, participant := range conv.Participants {
				participantIDs = append(participantIDs, participant.UserID)
			}

			// 8d. Send to all online participants
			service.hub.SendToUsers(participantIDs, jsonData)
			log.Printf("Broadcasted new message to %d participants", len(participantIDs))
		}
	}
	
	return savedMessage, nil
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

// GetMessageReceipts implements MessageService
// Returns all receipts for a message (only accessible by conversation participants)
func (service *messageServiceImpl) GetMessageReceipts(ctx context.Context, userID, messageID string) ([]domain.MessageReceipt, error) {
	// 1. Find the message
	message, err := service.messageRepo.FindByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFoundError("Message not found")
		}
		return nil, err
	}

	// 2. Validate: only the sender can view receipts
	if message.SenderID != userID {
		return nil, exceptions.NewForbiddenError("Only the message sender can view receipts")
	}

	// 3. Get all receipts for this message
	receipts, err := service.receiptRepo.FindByMessageID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	return receipts, nil
}

