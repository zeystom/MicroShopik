package controllers

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OrderController struct {
	orderService services.OrderService
}

func NewOrderController(s services.OrderService) *OrderController {
	return &OrderController{orderService: s}
}

// GetMyOrders @Summary Get orders for current user
// @Description Get all orders for the authenticated user
// @Tags orders
// @Produce json
// @Success 200 {array} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /orders [get]
func (oc *OrderController) GetMyOrders(c echo.Context) error {
	userID := c.Get("user_id").(int)

	orders, err := oc.orderService.GetByCustomerID(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orders)
}

// GetMyOrdersAsSeller @Summary Get orders for seller's products
// @Description Get all orders for products belonging to the authenticated seller
// @Tags orders
// @Produce json
// @Success 200 {array} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /orders/seller [get]
func (oc *OrderController) GetMyOrdersAsSeller(c echo.Context) error {
	userID := c.Get("user_id").(int)

	orders, err := oc.orderService.GetBySellerID(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orders)
}

// Create @Summary Create a new order
// @Description Create a new order for a customer
// @Tags orders
// @Accept json
// @Produce json
// @Param order body domain.Order true "Order object"
// @Success 201 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /orders [post]
func (oc *OrderController) Create(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var order domain.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Set the customer ID from JWT token
	order.CustomerID = &userID

	if err := oc.orderService.Create(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, order)
}

// GetByID @Summary Get an order by ID
// @Description Get order details by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func (oc *OrderController) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	order, err := oc.orderService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, order)
}

// GetByCustomerID @Summary Get orders by customer ID
// @Description Get all orders for a specific customer
// @Tags orders
// @Produce json
// @Param customerID path int true "Customer ID"
// @Success 200 {array} domain.Order
// @Failure 400 {object} map[string]string
// @Router /customers/{customerID}/orders [get]
func (oc *OrderController) GetByCustomerID(c echo.Context) error {
	customerID, err := strconv.Atoi(c.Param("customerID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid customer id"})
	}

	orders, err := oc.orderService.GetByCustomerID(customerID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orders)
}

// GetByStatus @Summary Get orders by status
// @Description Get all orders with a specific status
// @Tags orders
// @Produce json
// @Param status path string true "Order status"
// @Success 200 {array} domain.Order
// @Failure 400 {object} map[string]string
// @Router /orders/status/{status} [get]
func (oc *OrderController) GetByStatus(c echo.Context) error {
	status := c.Param("status")
	if status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "status is required"})
	}

	orders, err := oc.orderService.GetByStatus(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orders)
}

// Update @Summary Update an order
// @Description Update an order (admin/seller only)
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param order body domain.Order true "Order object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /orders/{id} [put]
func (oc *OrderController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	var order domain.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	order.ID = id
	if err := oc.orderService.Update(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order updated successfully"})
}

// Delete @Summary Delete an order
// @Description Delete an order (admin only)
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /orders/{id} [delete]
func (oc *OrderController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	if err := oc.orderService.Delete(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order deleted successfully"})
}

// UpdateStatus @Summary Update order status
// @Description Update the status of an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body map[string]string true "Status object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /orders/{id}/status [put]
func (oc *OrderController) UpdateStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	var request map[string]string
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	status, exists := request["status"]
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "status is required"})
	}

	if err := oc.orderService.UpdateStatus(id, status); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order status updated successfully"})
}

// GetByProductID @Summary Get orders by product ID
// @Description Get all orders for a specific product
// @Tags orders
// @Produce json
// @Param productID path int true "Product ID"
// @Success 200 {array} domain.Order
// @Failure 400 {object} map[string]string
// @Router /products/{productID}/orders [get]
func (oc *OrderController) GetByProductID(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("productID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	orders, err := oc.orderService.GetByProductID(productID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orders)
}

// ProcessOrder @Summary Process an order
// @Description Process a pending order (admin/seller only)
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /orders/{id}/process [post]
func (oc *OrderController) ProcessOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	if err := oc.orderService.ProcessOrder(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order processed successfully"})
}

// CancelOrder @Summary Cancel an order
// @Description Cancel a pending order (customer only)
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Param customerID path int true "Customer ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /orders/{id}/cancel/{customerID} [post]
func (oc *OrderController) CancelOrder(c echo.Context) error {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	customerID, err := strconv.Atoi(c.Param("customerID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid customer id"})
	}

	if err := oc.orderService.CancelOrder(orderID, customerID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order cancelled successfully"})
}

// ConfirmOrder @Summary Confirm an order payment
// @Description Confirm payment for a pending order (customer only)
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /orders/{id}/confirm [post]
func (oc *OrderController) ConfirmOrder(c echo.Context) error {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	userID := c.Get("user_id").(int)
	if err := oc.orderService.ConfirmOrder(orderID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order payment confirmed successfully"})
}
