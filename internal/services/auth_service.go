package services

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/repositories"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *domain.User) error
	Login(email, password string) (string, *domain.User, error)
}

type authService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(r repositories.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: r, jwtSecret: jwtSecret}
}

func (a *authService) Register(user *domain.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	user.CreatedAt = time.Now()

	err = a.userRepo.Create(user)
	if err != nil {
		return err
	}

	err = a.userRepo.AssignRole(user.ID, "buyer")
	if err != nil {
		return err
	}

	return nil
}

func (a *authService) Login(email, password string) (string, *domain.User, error) {
	findUser, err := a.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(password)); err != nil {
		return "", nil, errors.New("incorrect credentials")
	}

	findUser, err = a.userRepo.GetByID(findUser.ID)
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

	if err := a.userRepo.LastLoginUpdate(findUser.ID); err != nil {
		return "", nil, errors.New("failed to update last login time")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, findUser, nil
}
