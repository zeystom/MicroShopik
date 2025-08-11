package controllers

import (
	"MicroShopik/internal/domain"
	"MicroShopik/internal/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ParticipantController struct {
	participantService services.ParticipantService
}

func NewParticipantController(s services.ParticipantService) *ParticipantController {
	return &ParticipantController{participantService: s}
}

// Create @Summary Create a new participant
// @Description Add a user as participant to a conversation
// @Tags participants
// @Accept json
// @Produce json
// @Param participant body domain.Participant true "Participant object"
// @Success 201 {object} domain.Participant
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /participants [post]
func (pc *ParticipantController) Create(c echo.Context) error {
	var participant domain.Participant
	if err := c.Bind(&participant); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := pc.participantService.Create(&participant); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, participant)
}

// GetByConversationID @Summary Get participants by conversation ID
// @Description Get all participants for a specific conversation
// @Tags participants
// @Produce json
// @Param conversationID path int true "Conversation ID"
// @Success 200 {array} domain.Participant
// @Failure 400 {object} map[string]string
// @Router /conversations/{conversationID}/participants [get]
func (pc *ParticipantController) GetByConversationID(c echo.Context) error {
	conversationID, err := strconv.Atoi(c.Param("conversationID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
	}

	participants, err := pc.participantService.GetByConversationID(conversationID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, participants)
}

// GetByUserID @Summary Get conversations by user ID
// @Description Get all conversations where a user is a participant
// @Tags participants
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {array} domain.Participant
// @Failure 400 {object} map[string]string
// @Router /users/{userID}/participations [get]
func (pc *ParticipantController) GetByUserID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	participants, err := pc.participantService.GetByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, participants)
}

// Delete @Summary Remove participant from conversation
// @Description Remove a user as participant from a conversation
// @Tags participants
// @Produce json
// @Param conversationID path int true "Conversation ID"
// @Param userID path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /conversations/{conversationID}/participants/{userID} [delete]
func (pc *ParticipantController) Delete(c echo.Context) error {
	conversationID, err := strconv.Atoi(c.Param("conversationID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	if err := pc.participantService.Delete(conversationID, userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "participant removed successfully"})
}

// IsParticipant @Summary Check if user is participant
// @Description Check if a user is a participant in a conversation
// @Tags participants
// @Produce json
// @Param conversationID path int true "Conversation ID"
// @Param userID path int true "User ID"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Router /conversations/{conversationID}/participants/{userID}/check [get]
func (pc *ParticipantController) IsParticipant(c echo.Context) error {
	conversationID, err := strconv.Atoi(c.Param("conversationID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid conversation id"})
	}

	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	isParticipant, err := pc.participantService.IsParticipant(conversationID, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]bool{"is_participant": isParticipant})
}
