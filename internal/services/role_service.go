package services

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
	"errors"
)

type RoleService interface {
	GetByName(name string) (*domain.Role, error)
	GetAll() ([]domain.Role, error)
	Create(role *domain.Role) error
	Update(role *domain.Role) error
	Delete(id int) error
	AssignRoleToUser(userID int, roleName string) error
	RemoveRoleFromUser(userID int, roleName string) error
	GetUserRoles(userID int) ([]string, error)
}

type roleService struct {
	roleRepo repositories.RoleRepository
	userRepo repositories.UserRepository
}

func NewRoleService(roleRepo repositories.RoleRepository, userRepo repositories.UserRepository) RoleService {
	return &roleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

func (s *roleService) GetByName(name string) (*domain.Role, error) {
	return s.roleRepo.GetByName(name)
}

func (s *roleService) GetAll() ([]domain.Role, error) {
	return s.roleRepo.GetAll()
}

func (s *roleService) Create(role *domain.Role) error {
	if role.Name == "" {
		return errors.New("role name is required")
	}
	return s.roleRepo.Create(role)
}

func (s *roleService) Update(role *domain.Role) error {
	if role.ID == 0 {
		return errors.New("role ID is required")
	}
	return s.roleRepo.Update(role)
}

func (s *roleService) Delete(id int) error {
	return s.roleRepo.Delete(id)
}

func (s *roleService) AssignRoleToUser(userID int, roleName string) error {
	return s.userRepo.AssignRole(userID, roleName)
}

func (s *roleService) RemoveRoleFromUser(userID int, roleName string) error {
	return errors.New("remove role functionality not implemented yet")
}

func (s *roleService) GetUserRoles(userID int) ([]string, error) {
	return s.userRepo.GetRoles(userID)
}
