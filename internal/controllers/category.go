package controllers

import (
	"MicroShopik/internal/domain"
	domain2 "MicroShopik/internal/services/domain"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CategoryController struct {
	categoryService domain2.CategoryService
}

func NewCategoryController(s domain2.CategoryService) *CategoryController {
	return &CategoryController{categoryService: s}
}

// Create @Summary Create a new category
// @Description Create a new category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param category body domain.Category true "Category object"
// @Success 201 {object} domain.Category
// @Failure 400 {object} map[string]string
// @Security ApiKeyAuth
// @Router /categories [post]
func (cc *CategoryController) Create(c echo.Context) error {
	var category domain.Category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := cc.categoryService.Create(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, category)
}

// GetById @Summary Get a category by ID
// @Description Get category details by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} domain.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /categories/{id} [get]
func (cc *CategoryController) GetById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category id"})
	}

	category, err := cc.categoryService.GetCategoryById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, category)
}

// GetAll @Summary Get all categories
// @Description Get a list of all categories
// @Tags categories
// @Produce json
// @Success 200 {array} domain.Category
// @Failure 400 {object} map[string]string
// @Router /categories [get]
func (cc *CategoryController) GetAll(c echo.Context) error {
	categories, err := cc.categoryService.GetAllCategories()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, categories)
}

// Update @Summary Update a category
// @Description Update a category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body domain.Category true "Category object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /categories/{id} [put]
func (cc *CategoryController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category id"})
	}

	var category domain.Category
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	category.ID = id

	if err := cc.categoryService.Update(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "category updated successfully"})
}

// Delete @Summary Delete a category
// @Description Delete a category (admin only)
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /categories/{id} [delete]
func (cc *CategoryController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category id"})
	}

	category := domain.Category{ID: id}
	if err := cc.categoryService.Delete(&category); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "category deleted successfully"})
}
