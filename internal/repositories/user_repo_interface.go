package repositories

import (
	"MicroShopik/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	LastLoginUpdate(userID int) error
	GetRoles(userID int) ([]string, error)
	AssignRole(userID int, roleName string) error
	GetByEmail(email string) (*domain.User, error)
	GetByID(userID int) (*domain.User, error)
}
