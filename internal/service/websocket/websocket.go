package websocket

import (
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/store/chat"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Websocket struct {
	Clients       map[string]*websocket.Conn
	MessageChan   chan *models.MessageModel
	BroadcastChan chan *models.MessageModel
	store         chat.ChatStore
}

func NewWebsocket(store chat.ChatStore) *Websocket {
	return &Websocket{
		Clients:       make(map[string]*websocket.Conn),
		MessageChan:   make(chan *models.MessageModel),
		BroadcastChan: make(chan *models.MessageModel),
		store:         store,
	}
}

func (w *Websocket) SaveMessages() {
	for {
		message := <-w.MessageChan

		if _, err := w.store.Create(message); err != nil {
			fmt.Printf("error saving message: %s\n", err.Error())
			continue
		}

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
			// When receiver is not online
			fmt.Printf("Receiver not online!\n")
		}

		if socket, ok := w.Clients[message.Sender]; ok {
			err := socket.WriteJSON(message)
			if err != nil {
				fmt.Printf("error sending message to sender: %s\n", err.Error())
			}
		} else {
			fmt.Printf("Sender not online\n")
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
