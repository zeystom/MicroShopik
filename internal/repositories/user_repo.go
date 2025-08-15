package repositories

import (
	"MicroShopik/internal/domain"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) AssignRole(userID int, roleName string) error {
	var role domain.Role
	if err := r.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	return r.db.Model(&domain.User{ID: userID}).Association("Roles").Append(&role)
}

func (r *userRepository) RemoveRole(userID int, roleName string) error {
	var role domain.Role
	if err := r.db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	return r.db.Model(&domain.User{ID: userID}).Association("Roles").Delete(&role)
}

func (r *userRepository) GetRoles(userID int) ([]string, error) {
	var user domain.User
	err := r.db.Preload("Roles").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}
	return roleNames, nil
}

func (r *userRepository) LastLoginUpdate(userID int) error {
	return r.db.Model(&domain.User{}).
		Where("id = ?", userID).
		UpdateColumn("last_login", gorm.Expr("NOW()")).Error
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) GetByID(userID int) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Roles").Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", err)
		}

		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Delete(userID int) error {
	var user domain.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	return r.db.Delete(&user).Error
}
func (r *userRepository) GetAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Preload("Roles").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}
