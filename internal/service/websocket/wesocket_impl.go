package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type websocketImpl struct {
	clients map[string]*websocket.Conn
	channel chan *EventModel
}

func NewWebsocket() *websocketImpl {
	return &websocketImpl{
		clients: make(map[string]*websocket.Conn),
		channel: make(chan *EventModel),
	}
}

func (w *websocketImpl) Upgrade(c echo.Context) (*websocket.Conn, error) {
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

func (w *websocketImpl) RegisterClient(userId string, conn *websocket.Conn) {
	w.clients[userId] = conn
}

func (w *websocketImpl) RemoveClient(userId string) {
	delete(w.clients, userId)
}

func (w *websocketImpl) Send(evt *EventModel) {
	w.channel <- evt
}

func (w *websocketImpl) StartListening() {
	go func() {
		for {
			select {
			case evt := <-w.channel:
				if socket, ok := w.clients[evt.ReceiverId]; ok {
					err := socket.WriteJSON(evt)
					if err != nil {
						fmt.Printf("error sending message to recipient: %s\n", err.Error())
					}
				} else {
					fmt.Printf("Receiver not online!\n")
				}
			}
		}
	}()
}

func (w *websocketImpl) Stop() {
	close(w.channel)
}
