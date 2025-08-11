package services

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
	"errors"
)

type ParticipantService interface {
	Create(participant *domain.Participant) error
	GetByConversationID(conversationID int) ([]*domain.Participant, error)
	GetByUserID(userID int) ([]*domain.Participant, error)
	Delete(conversationID, userID int) error
	IsParticipant(conversationID, userID int) (bool, error)
	GetConversationParticipants(conversationID int) ([]*domain.Participant, error)
}

type participantService struct {
	participantRepo  repositories.ParticipantRepository
	conversationRepo repositories.ConversationRepository
	userRepo         repositories.UserRepository
}

func NewParticipantService(pRepo repositories.ParticipantRepository, cRepo repositories.ConversationRepository, uRepo repositories.UserRepository) ParticipantService {
	return &participantService{
		participantRepo:  pRepo,
		conversationRepo: cRepo,
		userRepo:         uRepo,
	}
}

func (s *participantService) Create(participant *domain.Participant) error {
	// Validate conversation exists
	_, err := s.conversationRepo.GetByID(participant.ConversationID)
	if err != nil {
		return errors.New("conversation not found")
	}

	// Check if user is already a participant
	isParticipant, err := s.IsParticipant(participant.ConversationID, participant.UserID)
	if err != nil {
		return err
	}
	if isParticipant {
		return errors.New("user is already a participant")
	}

	return s.participantRepo.Create(participant)
}

func (s *participantService) GetByConversationID(conversationID int) ([]*domain.Participant, error) {
	return s.participantRepo.GetByConversationID(conversationID)
}

func (s *participantService) GetByUserID(userID int) ([]*domain.Participant, error) {
	return s.participantRepo.GetByUserID(userID)
}

func (s *participantService) Delete(conversationID, userID int) error {
	// Check if user is a participant
	isParticipant, err := s.IsParticipant(conversationID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return errors.New("user is not a participant")
	}

	return s.participantRepo.Delete(conversationID, userID)
}

func (s *participantService) IsParticipant(conversationID, userID int) (bool, error) {
	return s.participantRepo.IsParticipant(conversationID, userID)
}

func (s *participantService) GetConversationParticipants(conversationID int) ([]*domain.Participant, error) {
	return s.participantRepo.GetByConversationID(conversationID)
}
