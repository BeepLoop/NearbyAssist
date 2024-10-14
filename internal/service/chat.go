package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/auth"
	ws "nearbyassist/internal/service/websocket"
	"nearbyassist/internal/store/chat"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ChatService struct {
	store     chat.ChatStore
	ws        *ws.Websocket
	jwt       auth.Authenticator
	encryptor auth.Encryption
}

func NewChatService(store chat.ChatStore, jwt auth.Authenticator, ws *ws.Websocket, encryptor auth.Encryption) *ChatService {
	return &ChatService{
		store:     store,
		jwt:       jwt,
		ws:        ws,
		encryptor: encryptor,
	}
}

func (s *ChatService) Websocket(c echo.Context) error {
	conn, err := s.ws.Upgrade(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error upgrading connection",
			Error:   err.Error(),
		})
	}
	defer conn.Close()

	token := c.QueryParam("token")
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	s.ws.RegisterClient(userId, conn)

	for {
		message := new(models.MessageModel)
		err := conn.ReadJSON(message)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.ws.UnregisterClient(userId)

				return nil
			}

			continue
		}

		s.ws.NewMessage(message)
	}
}

func (s *ChatService) GetMessages(c echo.Context) error {
	otherUser := c.Param("otherUserId")
	if otherUser == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Other user ID is required",
			Error:   "Other user ID is required",
		})
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	messages, err := s.store.GetMessages(userId, otherUser)
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

func (s *ChatService) GetConversations(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	conversations, err := s.store.GetConversations(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting conversations",
			Error:   err.Error(),
		})
	}

	for _, conversation := range conversations {
		if plain, err := s.encryptor.DecryptString(conversation.Name); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error decrypting name",
				Error:   auth.DECRYPTION_ERR,
			})
		} else {
			conversation.Name = plain
		}
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"conversations": conversations,
	})
}
