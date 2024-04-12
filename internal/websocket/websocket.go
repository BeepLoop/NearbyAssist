package websocket

import (
	"fmt"
	"nearbyassist/internal/db"
	"nearbyassist/internal/models"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Websocket struct {
	Clients       map[int]*websocket.Conn
	MessageChan   chan models.MessageModel
	BroadcastChan chan models.MessageModel
	DB            db.Database
}

func NewWebsocket(db db.Database) *Websocket {
	return &Websocket{
		Clients:       make(map[int]*websocket.Conn),
		MessageChan:   make(chan models.MessageModel),
		BroadcastChan: make(chan models.MessageModel),
		DB:            db,
	}
}

func (w *Websocket) SaveMessages() {
	for {
		message := <-w.MessageChan

		id, err := w.DB.NewMessage(message)
		if err != nil {
			fmt.Printf("error saving message: %s\n", err.Error())
			continue
		}

		message.Id = id
		w.BroadcastChan <- message
	}
}

func (w *Websocket) ForwardMessages() {
	for {
		message := <-w.BroadcastChan

		if socket, ok := w.Clients[message.Receiver]; ok {
			err := socket.WriteJSON(message)
			if err != nil {
				fmt.Printf("error sending message to recipient: %s\n", err.Error())
			}
		} else {
			fmt.Printf("Receiver not found!\n")
			continue
		}

		if socket, ok := w.Clients[message.Sender]; ok {
			err := socket.WriteJSON(message)
			if err != nil {
				fmt.Printf("error sending message to sender: %s\n", err.Error())
			}
		} else {
			fmt.Printf("Sender not found\n")
			continue
		}
	}
}

func (w *Websocket) Upgrade(c echo.Context) (*websocket.Conn, error) {
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
