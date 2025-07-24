package repositories

import "MicroShopik/internal/domain"

type ProductCRUDRepository interface {
	CreateProduct(product *domain.Product) (int, error)
	GetProductById(productId int) (*domain.Product, error)
	UpdateProduct(product *domain.Product) error
	DeleteProduct(productId int) error
}
