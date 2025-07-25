package repositories

import "MicroShopik/internal/domain"

type ProductQueryParams struct {
	SellerId    *int
	CategoryID  *int
	MinPrice    *int
	MaxPrice    *int
	IsActive    *bool
	Disposable  *bool
	SearchQuery *string
	Limit       *int
	Offset      *int
}
type ProductUpdateData struct {
	Title       *string
	Description *string
	Price       *int64
	CategoryID  *int
	IsActive    *bool
	Disposable  *bool
	MaxSales    *int
}

type ProductCRUDRepository interface {
	Create(product *domain.Product) (int, error)
	GetById(id int) (*domain.Product, error)
	Update(id int, data ProductUpdateData) error
	Delete(id int) error

	IsAvailable(id int) (bool, error)
	IncrementSoldCount(id int, delta int) error

	Find(params ProductQueryParams) ([]*domain.Product, error)
	Count(params ProductQueryParams) (int, error)
}
