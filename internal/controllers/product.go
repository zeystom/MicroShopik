package controllers

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/services/application"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductController struct {
	productAppService *application.ProductApplicationService
}

func NewProductController(s *application.ProductApplicationService) *ProductController {
	return &ProductController{productAppService: s}
}

// @Summary Create a new product
// @Description Create a new product for the authenticated user
// @Tags products
// @Accept json
// @Produce json
// @Param product body *domain.Product true "Product object"
// @Success 201 {object} *domain.Product
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /products [post]

func (pc *ProductController) Create(c echo.Context) error {
	var product domain.Product

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	userID := c.Get("user_id").(int)

	if err := pc.productAppService.CreateProductWithValidation(&product, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, product)
}

// GetById GetProduct godoc
// @Summary Get a product by ID
// @Description Get product details by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (pc *ProductController) GetById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	product, err := pc.productAppService.GetProductWithFullDetails(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

// Update UpdateProduct godoc
// @Summary Update a product
// @Description Update a product (only by the seller)
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body domain.Product true "Product object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /products/{id} [put]

func (pc *ProductController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	var product domain.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID := c.Get("user_id").(int)
	product.ID = id
	err = pc.productAppService.UpdateProductWithInventoryCheck(&product, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "product updated successfully"})
}

// Delete DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product (only by the seller)
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /products/{id} [delete]
func (pc *ProductController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}
	err = pc.productAppService.DeactivateProductWithOrderCheck(id, c.Get("user_id").(int))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "product deleted successfully"})
}

// IsAvailable checks if a product is available
// @Summary Check product availability
// @Description Returns whether the product with the given ID is available for purchase
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{} "Returns availability status"
// @Failure 400 {object} map[string]string "Invalid product ID or internal error"
// @Router /products/{id}/available [get]

func (pc *ProductController) IsAvailable(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	err = pc.productAppService.ValidateProductForPurchase(id, 0)
	available := err == nil
	if !available {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"available":  false,
			"product_id": id,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"available":  available,
		"product_id": id,
	})
}

// Find GetProducts godoc
// @Summary Get products with filters
// @Description Get a list of products with optional filtering and pagination
// @Tags products
// @Produce json
// @Param seller_id query int false "Filter by seller ID"
// @Param category_id query int false "Filter by category ID"
// @Param min_price query int false "Minimum price filter"
// @Param max_price query int false "Maximum price filter"
// @Param is_active query bool false "Filter by active status"
// @Param disposable query bool false "Filter by disposable status"
// @Param search query string false "Search in title and description"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Success 200 {object} []domain.Product
// @Failure 400 {object} map[string]string
// @Router /products [get]
func (pc *ProductController) Find(c echo.Context) error {
	params := domain.ProductQueryParams{}

	if sellerID := c.QueryParam("seller_id"); sellerID != "" {
		if id, err := strconv.Atoi(sellerID); err == nil {
			params.SellerId = &id
		}
	}

	if categoryID := c.QueryParam("category_id"); categoryID != "" {
		if id, err := strconv.Atoi(categoryID); err == nil {
			params.CategoryID = &id
		}
	}

	if minPrice := c.QueryParam("min_price"); minPrice != "" {
		if price, err := strconv.Atoi(minPrice); err == nil {
			params.MinPrice = &price
		}
	}

	if maxPrice := c.QueryParam("max_price"); maxPrice != "" {
		if price, err := strconv.Atoi(maxPrice); err == nil {
			params.MaxPrice = &price
		}
	}

	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			params.IsActive = &active
		}
	}

	if disposable := c.QueryParam("disposable"); disposable != "" {
		if disp, err := strconv.ParseBool(disposable); err == nil {
			params.Disposable = &disp
		}
	}

	if search := c.QueryParam("search"); search != "" {
		params.SearchQuery = &search
	}

	if limit := c.QueryParam("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			params.Limit = &l
		}
	} else {
		defaultLimit := 20
		params.Limit = &defaultLimit
	}

	if offset := c.QueryParam("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			params.Offset = &o
		}
	}

	products, err := pc.productAppService.FindProducts(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, products)
}

// Count CountProducts godoc
// @Summary Count products with filters
// @Description Get the count of products matching the filters
// @Tags products
// @Produce json
// @Param seller_id query int false "Filter by seller ID"
// @Param category_id query int false "Filter by category ID"
// @Param min_price query int false "Minimum price filter"
// @Param max_price query int false "Maximum price filter"
// @Param is_active query bool false "Filter by active status"
// @Param disposable query bool false "Filter by disposable status"
// @Param search query string false "Search in title and description"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Router /products/count [get]
func (pc *ProductController) Count(c echo.Context) error {
	params := domain.ProductQueryParams{}

	if sellerID := c.QueryParam("seller_id"); sellerID != "" {
		if id, err := strconv.Atoi(sellerID); err == nil {
			params.SellerId = &id
		}
	}

	if categoryID := c.QueryParam("category_id"); categoryID != "" {
		if id, err := strconv.Atoi(categoryID); err == nil {
			params.CategoryID = &id
		}
	}

	if minPrice := c.QueryParam("min_price"); minPrice != "" {
		if price, err := strconv.Atoi(minPrice); err == nil {
			params.MinPrice = &price
		}
	}

	if maxPrice := c.QueryParam("max_price"); maxPrice != "" {
		if price, err := strconv.Atoi(maxPrice); err == nil {
			params.MaxPrice = &price
		}
	}

	if isActive := c.QueryParam("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			params.IsActive = &active
		}
	}

	if disposable := c.QueryParam("disposable"); disposable != "" {
		if disp, err := strconv.ParseBool(disposable); err == nil {
			params.Disposable = &disp
		}
	}

	if search := c.QueryParam("search"); search != "" {
		params.SearchQuery = &search
	}

	count, err := pc.productAppService.CountProducts(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]int{"count": count})
}
