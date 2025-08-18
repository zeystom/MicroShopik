package repositories

import (
	"MicroShopik/internal/domain"

	"gorm.io/gorm"
)

type participantRepository struct {
	db *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) domain.ParticipantRepository {
	return &participantRepository{db: db}
}

func (r *participantRepository) Create(participant *domain.Participant) error {
	return r.db.Create(participant).Error
}

func (r *participantRepository) GetByConversationID(conversationID int) ([]*domain.Participant, error) {
	var participants []*domain.Participant
	err := r.db.Preload("User").Where("conversation_id = ?", conversationID).Find(&participants).Error
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *participantRepository) GetByUserID(userID int) ([]*domain.Participant, error) {
	var participants []*domain.Participant
	err := r.db.Preload("Conversation").Where("user_id = ?", userID).Find(&participants).Error
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *participantRepository) Delete(conversationID, userID int) error {
	return r.db.Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Delete(&domain.Participant{}).Error
}

func (r *participantRepository) IsParticipant(conversationID, userID int) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Participant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
