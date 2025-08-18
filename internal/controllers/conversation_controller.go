package controllers

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/services/application"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ConversationController struct {
	conversationAppService *application.ConversationApplicationService
}

func NewConversationController(s *application.ConversationApplicationService) *ConversationController {
	return &ConversationController{conversationAppService: s}
}

// Create @Summary Create a new conversation
// @Description Create a new conversation with participants and optional product
// @Tags conversations
// @Accept json
// @Produce json
// @Param conversation body domain.Conversation true "Conversation object"
// @Param participantIDs body []int true "Array of participant user IDs"
// @Param productID body int false "Product ID to link with conversation"
// @Success 201 {object} domain.Conversation
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /conversations [post]
func (cc *ConversationController) Create(c echo.Context) error {
	var request struct {
		Conversation   domain.Conversation `json:"conversation"`
		ParticipantIDs []int               `json:"participant_ids"`
		ProductID      *int                `json:"product_id,omitempty"` // ID товара для привязки
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Устанавливаем ProductID в conversation
	if request.ProductID != nil {
		request.Conversation.ProductID = request.ProductID
	}

	if err := cc.conversationAppService.CreateConversationWithParticipants(&request.Conversation, request.ParticipantIDs, request.ProductID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, request.Conversation)
}

// GetByID @Summary Get a conversation by ID
// @Description Get conversation details by ID
// @Tags conversations
// @Produce json
// @Param id path int true "Conversation ID"
// @Success 200 {object} domain.Conversation
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /conversations/{id} [get]
func (cc *ConversationController) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
	}

	userID := c.Get("user_id").(int)
	conversation, messages, err := cc.conversationAppService.GetConversationWithMessages(id, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	response := map[string]interface{}{
		"conversation": conversation,
		"messages":     messages,
	}

	return c.JSON(http.StatusOK, response)
}

// GetByUserID @Summary Get conversations by user ID
// @Description Get all conversations for a specific user
// @Tags conversations
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {array} domain.Conversation
// @Failure 400 {object} map[string]string
// @Router /users/{userID}/conversations [get]
func (cc *ConversationController) GetByUserID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	conversations, err := cc.conversationAppService.GetUserConversationsWithLastMessage(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, conversations)
}

// GetByProductID @Summary Get conversations by product ID
// @Description Get all conversations related to a specific product
// @Tags conversations
// @Produce json
// @Param productID path int true "Product ID"
// @Success 200 {array} domain.Conversation
// @Failure 400 {object} map[string]string
// @Router /products/{productID}/conversations [get]
func (cc *ConversationController) GetByProductID(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("productID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid product id"})
	}

	conversations, err := cc.conversationAppService.GetConversationsByProductID(productID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, conversations)
}

// Update @Summary Update a conversation
// @Description Update a conversation
// @Tags conversations
// @Accept json
// @Produce json
// @Param id path int true "Conversation ID"
// @Param conversation body domain.Conversation true "Conversation object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /conversations/{id} [put]
func (cc *ConversationController) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
	}

	var conversation domain.Conversation
	if err := c.Bind(&conversation); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	conversation.ID = id
	if err := cc.conversationAppService.UpdateConversation(&conversation); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "conversation updated successfully"})
}

// Delete @Summary Delete a conversation
// @Description Delete a conversation
// @Tags conversations
// @Produce json
// @Param id path int true "Conversation ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /conversations/{id} [delete]
func (cc *ConversationController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
	}

	if err := cc.conversationAppService.DeleteConversation(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "conversation deleted successfully"})
}

// AddParticipant @Summary Add participant to conversation
// @Description Add a user as participant to a conversation
// @Tags conversations
// @Produce json
// @Param id path int true "Conversation ID"
// @Param userID path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /conversations/{id}/participants/{userID} [post]
func (cc *ConversationController) AddParticipant(c echo.Context) error {
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	if err := cc.conversationAppService.AddParticipantToConversation(conversationID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "participant added successfully"})
}

// RemoveParticipant @Summary Remove participant from conversation
// @Description Remove a user as participant from a conversation
// @Tags conversations
// @Produce json
// @Param id path int true "Conversation ID"
// @Param userID path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /conversations/{id}/participants/{userID} [delete]
func (cc *ConversationController) RemoveParticipant(c echo.Context) error {
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	if err := cc.conversationAppService.RemoveParticipantFromConversation(conversationID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "participant removed successfully"})
}
