package repositories

import (
	"MicroShopik/internal/domain"
	"gorm.io/gorm"
)

type ProductItemRepository interface {
	Create(productItem *domain.ProductItem) error
	GetByID(id int) (*domain.ProductItem, error)
	GetByProductID(productID int) ([]*domain.ProductItem, error)
	GetAvailableByProductID(productID int) ([]*domain.ProductItem, error)
	Update(productItem *domain.ProductItem) error
	Delete(id int) error
	MarkAsUsed(id int) error
	MarkAsUsedTx(tx *gorm.DB, id int) error
}
