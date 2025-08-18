package repositories

import (
	"MicroShopik/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type conversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) domain.ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) Create(conversation *domain.Conversation) error {
	return r.db.Create(conversation).Error
}

func (r *conversationRepository) GetByID(id int) (*domain.Conversation, error) {
	var conversation domain.Conversation
	err := r.db.Preload("Product").Preload("Participants.User").Preload("Messages.Sender").Where("id = ?", id).First(&conversation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("conversation not found")
		}
		return nil, err
	}
	return &conversation, nil
}

func (r *conversationRepository) GetByUserID(userID int) ([]*domain.Conversation, error) {
	var conversations []*domain.Conversation
	err := r.db.Joins("JOIN participants ON conversations.id = participants.conversation_id").
		Where("participants.user_id = ?", userID).
		Preload("Product").
		Preload("Participants.User").
		Preload("Messages.Sender").
		Find(&conversations).Error
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *conversationRepository) GetByProductID(productID int) ([]*domain.Conversation, error) {
	var conversations []*domain.Conversation
	err := r.db.Where("product_id = ?", productID).
		Preload("Product").
		Preload("Participants.User").
		Preload("Messages.Sender").
		Find(&conversations).Error
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *conversationRepository) Update(conversation *domain.Conversation) error {
	return r.db.Save(conversation).Error
}

func (r *conversationRepository) Delete(id int) error {
	return r.db.Delete(&domain.Conversation{}, id).Error
}

func (r *conversationRepository) AddParticipant(conversationID, userID int) error {
	participant := &domain.Participant{
		ConversationID: conversationID,
		UserID:         userID,
	}
	return r.db.Create(participant).Error
}

func (r *conversationRepository) RemoveParticipant(conversationID, userID int) error {
	return r.db.Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Delete(&domain.Participant{}).Error
}

func (r *conversationRepository) BeginTx() (*gorm.DB, error) {
	tx := r.db.Begin()
	return tx, tx.Error
}

func (r *conversationRepository) CreateTx(tx *gorm.DB, conversation *domain.Conversation) error {
	return tx.Create(conversation).Error
}

func (r *conversationRepository) AddParticipantTx(tx *gorm.DB, conversationID, userID int) error {
	participant := &domain.Participant{
		ConversationID: conversationID,
		UserID:         userID,
	}
	return tx.Create(participant).Error
}

func (r *conversationRepository) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *conversationRepository) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}
