package handlers

import (
	"fmt"
	"nearbyassist/internal/encryption"
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type chatHandler struct {
	server *server.Server
}

func NewChatHandler(server *server.Server) *chatHandler {
	return &chatHandler{
		server: server,
	}
}

func (h *chatHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Chat base route",
	})
}

func (h *chatHandler) HandleWebsocket(c echo.Context) error {
	conn, err := h.server.Websocket.Upgrade(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer conn.Close()

	token := c.QueryParam("token")
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID not found in JWT")
	}

	h.server.Websocket.Clients[id] = conn
	fmt.Printf("id: %s connected\n", id)

	for {
		message := models.NewMessageModel(h.server.IdGen, h.server.DB)
		err := conn.ReadJSON(message)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if _, ok := h.server.Websocket.Clients[id]; ok {
					delete(h.server.Websocket.Clients, id)
				}

				fmt.Printf("client: %s disconnected\n", id)
				return nil
			}

			fmt.Printf("error reading message: %s\n", err.Error())
			continue
		}

		h.server.Websocket.MessageChan <- *message
	}
}

func (h *chatHandler) HandleGetMessages(c echo.Context) error {
	otherUser := c.Param("otherUserId")
	if otherUser == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID must be a number")
	}

	token := c.Request().Header.Get("Authorization")
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID not found in JWT")
	}

	user := models.NewUserModelWithId(id, h.server.DB)
	messages, err := user.GetMessages(otherUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"messages": messages,
	})
}

func (h *chatHandler) HandleGetConversations(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID not found in JWT")
	}

	user := models.NewUserModelWithId(id, h.server.DB)
	conversations, err := user.GetConversations()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error retrieving conversations")
	}

	for _, conversation := range conversations {
		if _, err := conversation.DecryptName(h.server.Encrypt.DecryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
		}
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"conversations": conversations,
	})
}
