package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Socket interface {
	Upgrade(echo.Context) (*websocket.Conn, error)
	RegisterClient(string, *websocket.Conn)
	RemoveClient(string)
	Send(*EventModel)
	StartListening()
	Stop()
}
