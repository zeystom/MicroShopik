package domain

import (
	"MicroShopik/internal/domain"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserService interface {
	Register(user *domain.User) error
	Login(email, password string) (string, *domain.User, error)
	GetByID(id int) (*domain.User, error)
	GetAll() ([]domain.User, error)
	Create(user *domain.User) error
	Update(user *domain.User) error
	Delete(id int) error
	ValidateUserExists(userID int) error
}
type userService struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewUserService(userRepo domain.UserRepository, jwtSecret string) UserService {
	return &userService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *userService) GetByID(id int) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetAll() ([]domain.User, error) {
	return s.userRepo.GetAll()
}

func (s *userService) Create(user *domain.User) error {
	return s.userRepo.Create(user)
}

func (s *userService) Update(user *domain.User) error {
	if user.Password != "" && !strings.HasPrefix(user.Password, "$2") {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashed)
	}
	return s.userRepo.Update(user)
}

func (s *userService) Delete(id int) error {
	return s.userRepo.Delete(id)
}

func (s *userService) ValidateUserExists(userID int) error {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	return nil
}

func (s *userService) Register(user *domain.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	user.CreatedAt = time.Now()

	err = s.userRepo.Create(user)
	if err != nil {
		return err
	}

	err = s.userRepo.AssignRole(user.ID, "buyer")
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Login(email, password string) (string, *domain.User, error) {
	findUser, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(password)); err != nil {
		return "", nil, errors.New("incorrect credentials")
	}

	findUser, err = s.userRepo.GetByID(findUser.ID)
	if err != nil {
		return "", nil, errors.New("failed to get user with roles")
	}

	var roleNames []string
	for _, role := range findUser.Roles {
		roleNames = append(roleNames, role.Name)
	}

	claims := domain.JWTClaims{
		UserID: findUser.ID,
		Roles:  roleNames,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour))},
	}

	if err := s.userRepo.LastLoginUpdate(findUser.ID); err != nil {
		return "", nil, errors.New("failed to update last login time")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, findUser, nil
}
