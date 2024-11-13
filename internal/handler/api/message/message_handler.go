package message

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	message_service "nearbyassist/internal/service/message"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type messageHandler struct {
	messageService *message_service.Service
}

func NewHandler(messageService *message_service.Service) *messageHandler {
	return &messageHandler{messageService: messageService}
}

func (h *messageHandler) ConnectWebsocket(c echo.Context) error {
	return h.messageService.ConnectWebsocket(c)
}

func (h *messageHandler) GetMessages(c echo.Context) error {
	otherUserId := c.Param("otherUserId")
	if otherUserId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Other user ID is required",
			Error:   "Other user ID is required",
		})
	}

	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	messages, err := h.messageService.GetMessages(bearerToken, otherUserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error getting messages",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"messages": messages,
	})
}

func (h *messageHandler) GetConversationList(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	conversations, err := h.messageService.GetConversationList(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving conversations",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"conversations": conversations,
	})
}

func (h *messageHandler) SendMessage(c echo.Context) error {
	req := new(request.NewMessagePayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	message := &models.MessageModel{
		Sender:   req.Sender,
		Receiver: req.Receiver,
		Content:  req.Content,
	}
	if err := h.messageService.SendMessage(message); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error sending message",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, http.StatusNoContent)
}
