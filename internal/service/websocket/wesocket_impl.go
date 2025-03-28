package websocket

import (
	"errors"
	"fmt"
	message_repo "nearbyassist/internal/repository/message"
	notification_service "nearbyassist/internal/service/notification"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type websocketImpl struct {
	clients map[string]*websocket.Conn
	channel chan *EventModel
}

func NewWebsocket(store message_repo.MessageRepository) *websocketImpl {
	return &websocketImpl{
		clients: make(map[string]*websocket.Conn),
		channel: make(chan *EventModel),
	}
}

func GetInstance() (*Socket, error) {
	if Instance != nil {
		return Instance, nil
	}

	return nil, errors.New(NIL_INSTANCE_ERR)
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
				notificationHeading := "Notification Heading"
				notificationContent := "Notification Content"

				oneSignal := notification_service.OneSignalInstance
				if oneSignal != nil {
					if err := oneSignal.NewUrgentNotification(evt.ReceiverId, notificationHeading, notificationContent); err != nil {
						fmt.Println(err.Error())
					}
				} else {
					fmt.Println("dum dum you forgot to initialize one signal")
				}

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
