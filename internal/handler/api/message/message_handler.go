package message

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	message_service "nearbyassist/internal/service/message"
	"nearbyassist/internal/utils"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
)

type messageHandler struct {
	messageService *message_service.Service
}

func NewHandler(messageService *message_service.Service) *messageHandler {
	return &messageHandler{messageService: messageService}
}

func (h *messageHandler) GetMessages(c echo.Context) error {
	otherUserId := c.Param("otherUserId")
	if otherUserId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Other user ID is required",
			Error:   "Other user ID is required",
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)
	messages, err := h.messageService.GetMessages(bearerToken, otherUserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error getting messages",
			Error:   err.Error(),
		})
	}

	resp := slices.AppendSeq(
		make([]response.Message, 0),
		utils.Map(messages, func(message *models.MessageModel) response.Message {
			return response.Message{
				Id:        message.Id,
				Sender:    message.Sender,
				Receiver:  message.Receiver,
				Content:   message.Content,
				Seen:      message.Seen,
				SeenAt:    message.SeenAt.String,
				CreatedAt: message.CreatedAt,
			}
		}),
	)

	return c.JSON(http.StatusOK, utils.Mapper{
		"messages": resp,
	})
}

func (h *messageHandler) MarkSeen(c echo.Context) error {
	req := new(request.MarkMessageSeenPayload)
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

	bearerToken := utils.BearerTokenFromHeader(c)
	if err := h.messageService.MarkSeen(bearerToken, req.MessageIDs); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error marking messages as seen",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *messageHandler) GetConversationList(c echo.Context) error {
	bearerToken := utils.BearerTokenFromHeader(c)
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
		Model: models.Model{
			Id:        req.Id,
			CreatedAt: req.CreatedAt,
		},
		Sender:   req.Sender,
		Receiver: req.Receiver,
		Content:  req.Content,
		Seen:     false,
	}
	if err := h.messageService.SendMessage(message); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error sending message",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, http.StatusNoContent)
}
