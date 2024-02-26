package message

import (
	"encoding/json"
	"fmt"
	query "nearbyassist/internal/db/query/message"
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

	client := c.Request().RemoteAddr
	fmt.Printf("%s connected\n", client)

	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("client disconnected")
				return nil
			}
			fmt.Println("error reading message")
		}

		message := new(types.Message)
		err = json.Unmarshal(bytes, message)
		if err != nil {
			fmt.Println("error unmarshalling message")
		}

		err = query.NewMessage(*message)
		if err != nil {
			fmt.Println("error saving message")
		}

		err = conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			fmt.Println("error sending message")
		}

	}
}
