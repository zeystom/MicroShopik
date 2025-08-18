package application

import (
	"MicroShopik/internal/domain"
	domain2 "MicroShopik/internal/services/domain"
	"errors"
)

type UserApplicationService struct {
	userService domain2.UserService
	roleService domain2.RoleService
}

func NewUserApplicationService(
	userService domain2.UserService,
	roleService domain2.RoleService,
) *UserApplicationService {
	return &UserApplicationService{
		userService: userService,
		roleService: roleService,
	}
}

func (s *UserApplicationService) RegisterUser(user *domain.User) error {
	if err := s.userService.Register(user); err != nil {
		return err
	}
	return nil
}

func (s *UserApplicationService) LoginUser(email, password string) (string, *domain.User, error) {
	token, user, err := s.userService.Login(email, password)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *UserApplicationService) CreateUserWithRoles(user *domain.User, roles []string) error {
	if err := s.userService.Create(user); err != nil {
		return err
	}

	for _, roleName := range roles {
		if err := s.roleService.AssignRoleToUser(user.ID, roleName); err != nil {
			return errors.New("failed to assign role: " + roleName)
		}
	}
	return nil
}

func (s *UserApplicationService) DeleteUserAccount(userID int) error {
	_, err := s.userService.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := s.userService.Delete(userID); err != nil {
		return err
	}

	return nil
}

func (s *UserApplicationService) UpdateUserProfile(user *domain.User) error {
	_, err := s.userService.GetByID(user.ID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := s.userService.Update(user); err != nil {
		return err
	}

	return nil
}

func (s *UserApplicationService) ValidateUserAccess(userID int) error {
	if err := s.userService.ValidateUserExists(userID); err != nil {
		return errors.New("user not found or access denied")
	}
	return nil
}

func (s *UserApplicationService) GetUserByID(userID int) (*domain.User, error) {
	return s.userService.GetByID(userID)
}
