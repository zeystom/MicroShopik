package controllers

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/dto"
	"MicroShopik/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(s services.AuthService) *AuthController {

	return &AuthController{authService: s}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User object"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (a *AuthController) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := a.authService.Register(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	resp := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return c.JSON(http.StatusCreated, resp)

}

// Login godoc
// @Summary Login a user
// @Description Login with email and password to receive a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body domain.AuthRequest true "Login credentials"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/login [post]
func (a *AuthController) Login(c echo.Context) error {
	var auth domain.AuthRequest
	if err := c.Bind(&auth); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	token, err := a.authService.Login(auth.Email, auth.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})

	}
	return c.JSON(http.StatusCreated, map[string]string{"token": token})
}
