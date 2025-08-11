package controllers

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductItemController struct {
	productItemService services.ProductItemService
}

func NewProductItemController(s services.ProductItemService) *ProductItemController {
	return &ProductItemController{productItemService: s}
}

// Create @Summary Create a new product item
// @Description Create a new product item for the authenticated user
// @Tags product-items
// @Accept json
// @Produce json
// @Param productItem body domain.ProductItem true "ProductItem object"
// @Success 201 {object} domain.ProductItem
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /product-items [post]
func (pic *ProductItemController) Create(c echo.Context) error {
	var productItem domain.ProductItem
	if err := c.Bind(&productItem); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := pic.productItemService.Create(&productItem); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, productItem)
}

// CreateBulk @Summary Create multiple product items
// @Description Create multiple product items for a product
// @Tags product-items
// @Accept json
// @Produce json
// @Param productID path int true "Product ID"
// @Param dataItems body []string true "Array of data strings"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /products/{productID}/items/bulk [post]
func (pic *ProductItemController) CreateBulk(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("productID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	var dataItems []string
	if err := c.Bind(&dataItems); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := pic.productItemService.CreateBulk(productID, dataItems); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "product items created successfully"})
}

// GetByID @Summary Get a product item by ID
// @Description Get product item details by ID
// @Tags product-items
// @Produce json
// @Param id path int true "ProductItem ID"
// @Success 200 {object} domain.ProductItem
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /product-items/{id} [get]
func (pic *ProductItemController) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product item id"})
	}

	productItem, err := pic.productItemService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, productItem)
}

// GetByProductID @Summary Get product items by product ID
// @Description Get all product items for a specific product
// @Tags product-items
// @Produce json
// @Param productID path int true "Product ID"
// @Success 200 {array} domain.ProductItem
// @Failure 400 {object} map[string]string
// @Router /products/{productID}/items [get]
func (pic *ProductItemController) GetByProductID(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("productID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	productItems, err := pic.productItemService.GetByProductID(productID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, productItems)
}

// GetAvailableByProductID @Summary Get available product items by product ID
// @Description Get all available (unused) product items for a specific product
// @Tags product-items
// @Produce json
// @Param productID path int true "Product ID"
// @Success 200 {array} domain.ProductItem
// @Failure 400 {object} map[string]string
// @Router /products/{productID}/items/available [get]
func (pic *ProductItemController) GetAvailableByProductID(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("productID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	productItems, err := pic.productItemService.GetAvailableByProductID(productID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, productItems)
}

// Update @Summary Update a product item
// @Description Update a product item (only by the seller)
// @Tags product-items
// @Accept json
// @Produce json
// @Param id path int true "ProductItem ID"
// @Param productItem body domain.ProductItem true "ProductItem object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /product-items/{id} [put]
func (pic *ProductItemController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product item id"})
	}

	var productItem domain.ProductItem
	if err := c.Bind(&productItem); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	productItem.ID = id
	if err := pic.productItemService.Update(&productItem); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "product item updated successfully"})
}

// Delete @Summary Delete a product item
// @Description Delete a product item (only by the seller)
// @Tags product-items
// @Produce json
// @Param id path int true "ProductItem ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /product-items/{id} [delete]
func (pic *ProductItemController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product item id"})
	}

	if err := pic.productItemService.Delete(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "product item deleted successfully"})
}

// MarkAsUsed @Summary Mark product item as used
// @Description Mark a product item as used (admin/seller only)
// @Tags product-items
// @Produce json
// @Param id path int true "ProductItem ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /product-items/{id}/mark-used [post]
func (pic *ProductItemController) MarkAsUsed(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product item id"})
	}

	if err := pic.productItemService.MarkAsUsed(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "product item marked as used"})
}
