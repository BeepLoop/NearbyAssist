package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type EventType string

const (
	EVT_MSSG  = "message"
	EVT_NOTIF = "notification"
	EVT_SYNC  = "sync"
)

type EventModel struct {
	ReceiverId string      `json:"-"`
	Type       EventType   `json:"type"`
	Payload    interface{} `json:"payload"`
}

type Socket interface {
	Upgrade(echo.Context) (*websocket.Conn, error)
	RegisterClient(string, *websocket.Conn)
	RemoveClient(string)
	Send(*EventModel)
	StartListening()
	Stop()
}
