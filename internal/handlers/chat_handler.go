package handlers

import (
	"fmt"
	"nearbyassist/internal/db/models"
	"nearbyassist/internal/server"
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

func (h *chatHandler) HandleGetMessages(c echo.Context) error {
	params := c.QueryString()

	model, err := models.MessageModelFactory(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	messages, err := model.GetMessages()
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
		message := models.NewMessageModel(h.server.DB)
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
