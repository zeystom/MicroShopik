package repositories

import (
	"MicroShopik/internal/domain"
)

type ParticipantRepository interface {
	Create(participant *domain.Participant) error
	GetByConversationID(conversationID int) ([]*domain.Participant, error)
	GetByUserID(userID int) ([]*domain.Participant, error)
	Delete(conversationID, userID int) error
	IsParticipant(conversationID, userID int) (bool, error)
}
