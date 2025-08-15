package repositories

import (
	"MicroShopik/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id int) (*domain.Order, error)
	GetByCustomerID(customerID int) ([]*domain.Order, error)
	GetBySellerID(sellerID int) ([]*domain.Order, error)
	GetByStatus(status string) ([]*domain.Order, error)
	Update(order *domain.Order) error
	Delete(id int) error
	UpdateStatus(id int, status string) error
	GetByProductID(productID int) ([]*domain.Order, error)
	GetAll() ([]*domain.Order, error)
	BeginTx() *gorm.DB
	UpdateStatusTx(tx *gorm.DB, id int, status string) error
}
