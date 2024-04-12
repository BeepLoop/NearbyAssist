package handlers

import (
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"

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

func (h *chatHandler) HandleGetMessages(c echo.Context) error {
	params := c.QueryString()

	values, err := models.MessageValueMapFactory(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	messages, err := h.server.DB.GetMessages(values["sender"], values["receiver"])
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, messages)
}

func (h *chatHandler) HandleWebsocket(c echo.Context) error {
	conn, err := h.server.Websocket.Upgrade(c)
	defer conn.Close()

	user := c.QueryParam("userId")
	userId, err := strconv.Atoi(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	h.server.Websocket.Clients[userId] = conn
	fmt.Printf("userId: %d connected\n", userId)

	for {
		message := models.NewMessageModel()
		err := conn.ReadJSON(message)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if _, ok := h.server.Websocket.Clients[userId]; ok {
					delete(h.server.Websocket.Clients, userId)
				}

				fmt.Printf("client: %d disconnected\n", userId)
				return nil
			}

			fmt.Printf("error reading message: %s\n", err.Error())
			continue
		}

		h.server.Websocket.MessageChan <- *message
	}
}

func (h *chatHandler) HandleGetConversations(c echo.Context) error {
	userId := c.QueryParam("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	conversations, err := h.server.DB.GetAllUserConversations(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, conversations)
}
