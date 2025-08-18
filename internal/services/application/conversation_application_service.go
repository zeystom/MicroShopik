package application

import (
	"MicroShopik/internal/domain"
	domain2 "MicroShopik/internal/services/domain"
	"errors"
)

type ConversationApplicationService struct {
	conversationService domain2.ConversationService
	messageService      domain2.MessageService
	userService         domain2.UserService
	productService      domain2.ProductService
	participantService  domain2.ParticipantService
}

func NewConversationApplicationService(
	conversationService domain2.ConversationService,
	messageService domain2.MessageService,
	userService domain2.UserService,
	productService domain2.ProductService,
	participantService domain2.ParticipantService,
) *ConversationApplicationService {
	return &ConversationApplicationService{
		conversationService: conversationService,
		messageService:      messageService,
		userService:         userService,
		productService:      productService,
		participantService:  participantService,
	}
}

func (s *ConversationApplicationService) CreateConversationWithParticipants(conversation *domain.Conversation, participantIDs []int, productID *int) error {
	if len(participantIDs) < 2 {
		return errors.New("conversation must have at least 2 participants")
	}

	for _, userID := range participantIDs {
		if err := s.userService.ValidateUserExists(userID); err != nil {
			return errors.New("invalid participant ID")
		}
	}

	if productID != nil {
		if err := s.productService.ValidateProductExists(*productID); err != nil {
			return errors.New("invalid product ID")
		}
		conversation.ProductID = productID
	}

	if err := s.conversationService.Create(conversation, participantIDs); err != nil {
		return err
	}

	return nil
}

func (s *ConversationApplicationService) SendMessageWithValidation(conversationID int, senderID int, content string) error {
	_, err := s.conversationService.GetByID(conversationID)
	if err != nil {
		return errors.New("conversation not found")
	}

	participants, err := s.participantService.GetByConversationID(conversationID)
	if err != nil {
		return errors.New("failed to get conversation participants")
	}

	isParticipant := false
	for _, participant := range participants {
		if participant.UserID == senderID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return errors.New("user is not a participant in this conversation")
	}

	message := &domain.Message{
		ConversationID: conversationID,
		SenderID:       &senderID,
		Content:        content,
	}

	return s.messageService.Create(message)
}

func (s *ConversationApplicationService) GetConversationWithMessages(conversationID, userID int) (*domain.Conversation, []*domain.Message, error) {
	conversation, err := s.conversationService.GetByID(conversationID)
	if err != nil {
		return nil, nil, errors.New("conversation not found")
	}

	participants, err := s.participantService.GetByConversationID(conversationID)
	if err != nil {
		return nil, nil, errors.New("failed to get conversation participants")
	}

	isParticipant := false
	for _, participant := range participants {
		if participant.UserID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, nil, errors.New("user is not a participant in this conversation")
	}

	messages, err := s.messageService.GetByConversationID(conversationID)
	if err != nil {
		return nil, nil, err
	}

	return conversation, messages, nil
}

func (s *ConversationApplicationService) GetUserConversationsWithLastMessage(userID int) ([]*domain.Conversation, error) {
	if err := s.userService.ValidateUserExists(userID); err != nil {
		return nil, errors.New("user not found")
	}

	conversations, err := s.conversationService.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, conversation := range conversations {
		messages, err := s.messageService.GetByConversationID(conversation.ID)
		if err != nil {
			continue
		}

		if len(messages) > 0 {
			lastMessage := messages[len(messages)-1]
			conversation.LastMessage = lastMessage.Content
			conversation.LastMessageTime = &lastMessage.CreatedAt
		}
	}

	return conversations, nil
}

func (s *ConversationApplicationService) AddParticipantToConversation(conversationID, userID int) error {
	_, err := s.conversationService.GetByID(conversationID)
	if err != nil {
		return errors.New("conversation not found")
	}

	if err := s.userService.ValidateUserExists(userID); err != nil {
		return errors.New("user not found")
	}

	existingParticipants, err := s.participantService.GetByConversationID(conversationID)
	if err != nil {
		return err
	}

	for _, participant := range existingParticipants {
		if participant.UserID == userID {
			return errors.New("user is already a participant")
		}
	}

	participant := &domain.Participant{
		ConversationID: conversationID,
		UserID:         userID,
	}

	if err := s.participantService.Create(participant); err != nil {
		return err
	}

	systemMessage := &domain.Message{
		ConversationID: conversationID,
		SenderID:       nil,
		Content:        "User joined the conversation",
		IsSystem:       true,
	}

	return s.messageService.Create(systemMessage)
}

func (s *ConversationApplicationService) GetConversationsByProductID(productID int) ([]*domain.Conversation, error) {
	if err := s.productService.ValidateProductExists(productID); err != nil {
		return nil, errors.New("product not found")
	}

	return s.conversationService.GetByProductID(productID)
}

func (s *ConversationApplicationService) UpdateConversation(conversation *domain.Conversation) error {
	_, err := s.conversationService.GetByID(conversation.ID)
	if err != nil {
		return errors.New("conversation not found")
	}

	return s.conversationService.Update(conversation)
}

func (s *ConversationApplicationService) DeleteConversation(conversationID int) error {
	_, err := s.conversationService.GetByID(conversationID)
	if err != nil {
		return errors.New("conversation not found")
	}

	return s.conversationService.Delete(conversationID)
}

func (s *ConversationApplicationService) RemoveParticipantFromConversation(conversationID, userID int) error {
	_, err := s.conversationService.GetByID(conversationID)
	if err != nil {
		return errors.New("conversation not found")
	}

	participants, err := s.participantService.GetByConversationID(conversationID)
	if err != nil {
		return err
	}

	if len(participants) <= 2 {
		return errors.New("cannot remove participant from conversation with only 2 participants")
	}

	if err := s.participantService.DeleteByConversationAndUser(conversationID, userID); err != nil {
		return err
	}

	systemMessage := &domain.Message{
		ConversationID: conversationID,
		SenderID:       nil,
		Content:        "User left the conversation",
		IsSystem:       true,
	}

	return s.messageService.Create(systemMessage)
}
