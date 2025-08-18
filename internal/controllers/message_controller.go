package controllers

import (
	"MicroShopik/internal/domain"
	domain2 "MicroShopik/internal/services/domain"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MessageController struct {
	messageService domain2.MessageService
}

func NewMessageController(s domain2.MessageService) *MessageController {
	return &MessageController{messageService: s}
}

// Create @Summary Create a new message
// @Description Send a message in a conversation
// @Tags messages
// @Accept json
// @Produce json
// @Param message body domain.Message true "Message object"
// @Success 201 {object} domain.Message
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /messages [post]
func (mc *MessageController) Create(c echo.Context) error {
	var message domain.Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := mc.messageService.Create(&message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, message)
}

// GetByID @Summary Get a message by ID
// @Description Get message details by ID
// @Tags messages
// @Produce json
// @Param id path int true "Message ID"
// @Success 200 {object} domain.Message
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /messages/{id} [get]
func (mc *MessageController) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid message id"})
	}

	message, err := mc.messageService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, message)
}

// GetByConversationID @Summary Get messages by conversation ID
// @Description Get messages for a specific conversation with pagination
// @Tags messages
// @Produce json
// @Param conversationID path int true "Conversation ID"
// @Param limit query int false "Limit number of messages (default: 50)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Success 200 {array} domain.Message
// @Failure 400 {object} map[string]string
// @Router /conversations/{conversationID}/messages [get]
func (mc *MessageController) GetByConversationID(c echo.Context) error {
	conversationID, err := strconv.Atoi(c.Param("conversationID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 50
	}

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	messages, err := mc.messageService.GetByConversationID(conversationID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, messages)
}

// GetByOrderID @Summary Get messages by order ID
// @Description Get all messages related to a specific order
// @Tags messages
// @Produce json
// @Param orderID path int true "Order ID"
// @Success 200 {array} domain.Message
// @Failure 400 {object} map[string]string
// @Router /orders/{orderID}/messages [get]
func (mc *MessageController) GetByOrderID(c echo.Context) error {
	orderID, err := strconv.Atoi(c.Param("orderID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	messages, err := mc.messageService.GetByOrderID(orderID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, messages)
}

// Update @Summary Update a message
// @Description Update a message (only by the sender)
// @Tags messages
// @Accept json
// @Produce json
// @Param id path int true "Message ID"
// @Param message body domain.Message true "Message object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /messages/{id} [put]
func (mc *MessageController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid message id"})
	}

	var message domain.Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	message.ID = id
	if err := mc.messageService.Update(&message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "message updated successfully"})
}

// Delete @Summary Delete a message
// @Description Delete a message (only by the sender)
// @Tags messages
// @Produce json
// @Param id path int true "Message ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /messages/{id} [delete]
func (mc *MessageController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid message id"})
	}

	if err := mc.messageService.Delete(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "message deleted successfully"})
}

// GetSystemMessages @Summary Get system messages
// @Description Get all system messages for a conversation
// @Tags messages
// @Produce json
// @Param conversationID path int true "Conversation ID"
// @Success 200 {array} domain.Message
// @Failure 400 {object} map[string]string
// @Router /conversations/{conversationID}/messages/system [get]
func (mc *MessageController) GetSystemMessages(c echo.Context) error {
	conversationID, err := strconv.Atoi(c.Param("conversationID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
	}

	messages, err := mc.messageService.GetSystemMessages(conversationID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, messages)
}

// SendSystemMessage @Summary Send system message
// @Description Send a system message to a conversation
// @Tags messages
// @Accept json
// @Produce json
// @Param conversationID path int true "Conversation ID"
// @Param request body map[string]interface{} true "Message text and optional order ID"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /conversations/{conversationID}/messages/system [post]
func (mc *MessageController) SendSystemMessage(c echo.Context) error {
	conversationID, err := strconv.Atoi(c.Param("conversationID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
	}

	var request map[string]interface{}
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	text, ok := request["text"].(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "text is required"})
	}

	var orderID *int
	if orderIDVal, exists := request["order_id"]; exists && orderIDVal != nil {
		if orderIDInt, ok := orderIDVal.(float64); ok {
			orderIDInt := int(orderIDInt)
			orderID = &orderIDInt
		}
	}

	if err := mc.messageService.SendSystemMessage(conversationID, text, orderID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "system message sent successfully"})
}
