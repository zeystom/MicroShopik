package repositories

import "MicroShopik/internal/domain"

type CategoryRepository interface {
	Create(category *domain.Category) error
	GetCategoryById(id int) (*domain.Category, error)
	GetAllCategories() (*[]domain.Category, error)
	Update(category *domain.Category) (*domain.Category, error)
	Delete(category *domain.Category) error
}
