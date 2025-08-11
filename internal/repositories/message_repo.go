package repositories

import (
	"MicroShopik/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *domain.Message) error {
	if err := r.db.Create(message).Error; err != nil {
		return err
	}

	// Reload with associations to return a fully populated entity
	return r.db.Preload("Sender").Preload("Order").First(message, message.ID).Error
}

func (r *messageRepository) GetByID(id int) (*domain.Message, error) {
	var message domain.Message
	err := r.db.Preload("Sender").Preload("Order").Where("id = ?", id).First(&message).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("message not found")
		}
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) GetByConversationID(conversationID int, limit, offset int) ([]*domain.Message, error) {
	var messages []*domain.Message
	query := r.db.Preload("Sender").Where("conversation_id = ?", conversationID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) GetByOrderID(orderID int) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.db.Preload("Sender").Where("order_id = ?", orderID).Order("created_at ASC").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) Update(message *domain.Message) error {
	return r.db.Save(message).Error
}

func (r *messageRepository) Delete(id int) error {
	return r.db.Delete(&domain.Message{}, id).Error
}

func (r *messageRepository) GetSystemMessages(conversationID int) ([]*domain.Message, error) {
	var messages []*domain.Message
	err := r.db.Where("conversation_id = ? AND is_system = ?", conversationID, true).
		Order("created_at ASC").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}
