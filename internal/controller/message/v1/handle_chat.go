package message

import (
	"fmt"
	"nearbyassist/internal/db/models"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var clients = make(map[int]*websocket.Conn)
var messageChan = make(chan models.MessageModel)
var broadcastChan = make(chan models.MessageModel)

func HandleChat(c echo.Context) error {
	conn, err := websocketUpgrader(c)
	defer conn.Close()

	user := c.QueryParam("userId")
	userId, err := strconv.Atoi(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	clients[userId] = conn
	fmt.Printf("userId: %d connected\n", userId)

	for {
		message := models.NewMessageModel()
		err := conn.ReadJSON(message)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				if _, ok := clients[userId]; ok {
					delete(clients, userId)
				}

				fmt.Printf("client: %d disconnected\n", userId)
				return nil
			}

			fmt.Printf("error reading message: %s\n", err.Error())
			continue
		}

		messageChan <- *message
	}
}
