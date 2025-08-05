package controllers

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RoleController struct {
	roleService services.RoleService
}

func NewRoleController(s services.RoleService) *RoleController {
	return &RoleController{roleService: s}
}

// GetAll @Summary Get all roles
// @Description Get a list of all available roles
// @Tags roles
// @Produce json
// @Success 200 {array} domain.Role
// @Failure 400 {object} map[string]string
// @Router /roles [get]
func (rc *RoleController) GetAll(c echo.Context) error {
	roles, err := rc.roleService.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, roles)
}

// GetByName @Summary Get role by name
// @Description Get a specific role by name
// @Tags roles
// @Produce json
// @Param name path string true "Role name"
// @Success 200 {object} domain.Role
// @Failure 404 {object} map[string]string
// @Router /roles/{name} [get]
func (rc *RoleController) GetByName(c echo.Context) error {
	name := c.Param("name")
	role, err := rc.roleService.GetByName(name)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, role)
}

// Create @Summary Create a new role
// @Description Create a new role (admin only)
// @Tags roles
// @Accept json
// @Produce json
// @Param role body domain.Role true "Role object"
// @Success 201 {object} domain.Role
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /roles [post]
func (rc *RoleController) Create(c echo.Context) error {
	var role domain.Role
	if err := c.Bind(&role); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := rc.roleService.Create(&role)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, role)
}

// AssignRoleToUser @Summary Assign role to user
// @Description Assign a role to a specific user (admin only)
// @Tags roles
// @Produce json
// @Param user_id path int true "User ID"
// @Param role_name path string true "Role name"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /users/{user_id}/roles/{role_name} [post]
func (rc *RoleController) AssignRoleToUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	roleName := c.Param("role_name")
	err = rc.roleService.AssignRoleToUser(userID, roleName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "role assigned successfully"})
}

// GetUserRoles @Summary Get user roles
// @Description Get all roles assigned to a specific user
// @Tags roles
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} string
// @Failure 400 {object} map[string]string
// @Router /users/{user_id}/roles [get]
func (rc *RoleController) GetUserRoles(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	roles, err := rc.roleService.GetUserRoles(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, roles)
}
