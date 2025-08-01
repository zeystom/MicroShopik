package repositories

import (
	"MicroShopik/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByName(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetAll() ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

func (r *roleRepository) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) Update(role *domain.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id int) error {
	return r.db.Delete(&domain.Role{}, id).Error
}
