package repositories

import (
	"MicroShopik/internal/domain"
)

type MessageRepository interface {
	Create(message *domain.Message) error
	GetByID(id int) (*domain.Message, error)
	GetByConversationID(conversationID int, limit, offset int) ([]*domain.Message, error)
	GetByOrderID(orderID int) ([]*domain.Message, error)
	Update(message *domain.Message) error
	Delete(id int) error
	GetSystemMessages(conversationID int) ([]*domain.Message, error)
}
