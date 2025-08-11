package services

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
	"errors"
	"fmt"
)

type ConversationService interface {
	Create(conversation *domain.Conversation, participantIDs []int) error
	GetByID(id int) (*domain.Conversation, error)
	GetByUserID(userID int) ([]*domain.Conversation, error)
	Update(conversation *domain.Conversation) error
	Delete(id int) error
	AddParticipant(conversationID, userID int) error
	RemoveParticipant(conversationID, userID int) error
	IsParticipant(conversationID, userID int) (bool, error)
}

type conversationService struct {
	conversationRepo repositories.ConversationRepository
	participantRepo  repositories.ParticipantRepository
	userRepo         repositories.UserRepository
}

func NewConversationService(cRepo repositories.ConversationRepository, pRepo repositories.ParticipantRepository, uRepo repositories.UserRepository) ConversationService {
	return &conversationService{
		conversationRepo: cRepo,
		participantRepo:  pRepo,
		userRepo:         uRepo,
	}
}

func (s *conversationService) Create(conversation *domain.Conversation, participantIDs []int) error {
	for _, userID := range participantIDs {
		_, err := s.userRepo.GetByID(userID)
		if err != nil {
			return fmt.Errorf("invalid participant ID %d", userID)
		}
	}

	tx, err := s.conversationRepo.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.conversationRepo.CreateTx(tx, conversation); err != nil {
		return err
	}

	for _, userID := range participantIDs {
		if err := s.conversationRepo.AddParticipantTx(tx, conversation.ID, userID); err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (s *conversationService) GetByID(id int) (*domain.Conversation, error) {
	return s.conversationRepo.GetByID(id)
}

func (s *conversationService) GetByUserID(userID int) ([]*domain.Conversation, error) {
	return s.conversationRepo.GetByUserID(userID)
}

func (s *conversationService) Update(conversation *domain.Conversation) error {
	return s.conversationRepo.Update(conversation)
}

func (s *conversationService) Delete(id int) error {
	return s.conversationRepo.Delete(id)
}

func (s *conversationService) AddParticipant(conversationID, userID int) error {
	isParticipant, err := s.IsParticipant(conversationID, userID)
	if err != nil {
		return err
	}
	if isParticipant {
		return errors.New("user is already a participant")
	}

	return s.conversationRepo.AddParticipant(conversationID, userID)
}

func (s *conversationService) RemoveParticipant(conversationID, userID int) error {
	isParticipant, err := s.IsParticipant(conversationID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.New("user is not a participant")
	}

	return s.conversationRepo.RemoveParticipant(conversationID, userID)
}

func (s *conversationService) IsParticipant(conversationID, userID int) (bool, error) {
	return s.participantRepo.IsParticipant(conversationID, userID)
}
