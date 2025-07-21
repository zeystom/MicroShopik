package controllers

import (
	"MicroShopik/internal/domain"
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
// @Param user body domain.User true "User object"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]string
// @Router /register [post]
func (a *AuthController) Register(c echo.Context) error {
	var user domain.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := a.authService.Register(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
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
// @Router /login [post]
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
