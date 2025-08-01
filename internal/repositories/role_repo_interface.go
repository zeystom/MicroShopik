package repositories

import "MicroShopik/internal/domain"

type RoleRepository interface {
	GetByName(name string) (*domain.Role, error)
	GetAll() ([]domain.Role, error)
	Create(role *domain.Role) error
	Update(role *domain.Role) error
	Delete(id int) error
}
