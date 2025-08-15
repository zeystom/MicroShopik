package repositories

import (
	"MicroShopik/internal/domain"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	Create(conversation *domain.Conversation) error
	GetByID(id int) (*domain.Conversation, error)
	GetByUserID(userID int) ([]*domain.Conversation, error)
	GetByProductID(productID int) ([]*domain.Conversation, error) // Новый метод для поиска по товару
	Update(conversation *domain.Conversation) error
	Delete(id int) error
	AddParticipant(conversationID, userID int) error
	RemoveParticipant(conversationID, userID int) error

	BeginTx() (*gorm.DB, error)
	CreateTx(tx *gorm.DB, conversation *domain.Conversation) error
	AddParticipantTx(tx *gorm.DB, conversationID, userID int) error
	Commit(tx *gorm.DB) error
	Rollback(tx *gorm.DB) error
}
