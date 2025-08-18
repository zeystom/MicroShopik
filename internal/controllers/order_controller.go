package controllers

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/services/application"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OrderController struct {
	orderAppService *application.OrderApplicationService
}

func NewOrderController(orderAppService *application.OrderApplicationService) *OrderController {
	return &OrderController{orderAppService: orderAppService}
}

func (oc *OrderController) GetMyOrders(c echo.Context) error {
	userID := c.Get("user_id").(int)

	orders, err := oc.orderAppService.GetMyOrders(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orders)
}

func (oc *OrderController) GetMyOrdersAsSeller(c echo.Context) error {
	userID := c.Get("user_id").(int)

	orders, err := oc.orderAppService.GetMyOrdersAsSeller(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orders)
}

func (oc *OrderController) Create(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var order domain.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	order.CustomerID = &userID

	if err := oc.orderAppService.CreateOrder(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, order)
}

func (oc *OrderController) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	order, err := oc.orderAppService.GetOrderByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, order)
}

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

	if err := oc.orderAppService.UpdateOrderStatus(id, status); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order status updated successfully"})
}

func (oc *OrderController) ProcessOrder(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	if err := oc.orderAppService.ProcessOrder(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order processed successfully"})
}

func (oc *OrderController) CancelOrder(c echo.Context) error {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	customerID, err := strconv.Atoi(c.Param("customerID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid customer id"})
	}

	if err := oc.orderAppService.CancelOrder(orderID, customerID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order cancelled successfully"})
}

func (oc *OrderController) ConfirmOrder(c echo.Context) error {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	userID := c.Get("user_id").(int)

	if err := oc.orderAppService.ConfirmOrder(orderID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "order payment confirmed successfully"})
}
