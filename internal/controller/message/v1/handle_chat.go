package message

import (
	"encoding/json"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

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

	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}

		message := new(types.Message)
		err = json.Unmarshal(bytes, message)
		if err != nil {
			c.Logger().Error(err)
		}

		err = conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			c.Logger().Error(err)
		}

	}
}
