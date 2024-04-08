package message

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func websocketUpgrader(c echo.Context) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
