package services

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
	"errors"
)

type MessageService interface {
	Create(message *domain.Message) error
	GetByID(id int) (*domain.Message, error)
	GetByConversationID(conversationID int, limit, offset int) ([]*domain.Message, error)
	GetByOrderID(orderID int) ([]*domain.Message, error)
	Update(message *domain.Message) error
	Delete(id int) error
	GetSystemMessages(conversationID int) ([]*domain.Message, error)
	SendSystemMessage(conversationID int, text string, orderID *int) error
}

type messageService struct {
	messageRepo      repositories.MessageRepository
	conversationRepo repositories.ConversationRepository
	participantRepo  repositories.ParticipantRepository
	orderRepo        repositories.OrderRepository
}

func NewMessageService(mRepo repositories.MessageRepository, cRepo repositories.ConversationRepository, pRepo repositories.ParticipantRepository, oRepo repositories.OrderRepository) MessageService {
	return &messageService{
		messageRepo:      mRepo,
		conversationRepo: cRepo,
		participantRepo:  pRepo,
		orderRepo:        oRepo,
	}
}

func (s *messageService) Create(message *domain.Message) error {

	_, err := s.conversationRepo.GetByID(message.ConversationID)
	if err != nil {
		return errors.New("conversation not found")
	}

	if message.SenderID != nil {
		isParticipant, err := s.participantRepo.IsParticipant(message.ConversationID, *message.SenderID)
		if err != nil {
			return err
		}
		if !isParticipant {
			return errors.New("sender is not a participant in this conversation")
		}
	}

	if message.OrderID != nil {
		_, err := s.orderRepo.GetByID(*message.OrderID)
		if err != nil {
			return errors.New("order not found")
		}
	}

	return s.messageRepo.Create(message)
}

func (s *messageService) GetByID(id int) (*domain.Message, error) {
	return s.messageRepo.GetByID(id)
}

func (s *messageService) GetByConversationID(conversationID int, limit, offset int) ([]*domain.Message, error) {
	return s.messageRepo.GetByConversationID(conversationID, limit, offset)
}

func (s *messageService) GetByOrderID(orderID int) ([]*domain.Message, error) {
	return s.messageRepo.GetByOrderID(orderID)
}

func (s *messageService) Update(message *domain.Message) error {
	return s.messageRepo.Update(message)
}

func (s *messageService) Delete(id int) error {
	return s.messageRepo.Delete(id)
}

func (s *messageService) GetSystemMessages(conversationID int) ([]*domain.Message, error) {
	return s.messageRepo.GetSystemMessages(conversationID)
}

func (s *messageService) SendSystemMessage(conversationID int, text string, orderID *int) error {
	message := &domain.Message{
		ConversationID: conversationID,
		Text:           text,
		IsSystem:       true,
		OrderID:        orderID,
	}

	return s.messageRepo.Create(message)
}
