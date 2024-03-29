package message

import (
	"fmt"
	"nearbyassist/internal/types"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var clients = make(map[int]*websocket.Conn)
var messageChan = make(chan types.Message)
var broadcastChan = make(chan types.Message)

func HandleChat(c echo.Context) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	defer conn.Close()

	user := c.QueryParam("userId")
	if user == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing data",
		})
	}

	userId, err := strconv.Atoi(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
	}

	clients[userId] = conn
	fmt.Printf("userId: %d connected\n", userId)

	for {
		message := new(types.Message)
		err := conn.ReadJSON(message)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("client: %d disconnected\n", userId)
				delete(clients, userId)
				return nil
			}
			fmt.Printf("error reading message: %s\n", err.Error())
			continue
		}

		messageChan <- *message
	}
}
