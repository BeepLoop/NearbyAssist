package websocket

import (
	"fmt"
	"nearbyassist/internal/models"
	message_repo "nearbyassist/internal/repository/message"
	notification_service "nearbyassist/internal/service/notification"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Websocket struct {
	clients     map[string]*websocket.Conn
	messageChan chan *models.MessageModel
	senderChan  chan *models.MessageModel
	saverChan   chan *models.MessageModel
	store       message_repo.MessageRepository
}

func NewWebsocket(store message_repo.MessageRepository) *Websocket {
	return &Websocket{
		clients:     make(map[string]*websocket.Conn),
		messageChan: make(chan *models.MessageModel),
		senderChan:  make(chan *models.MessageModel),
		saverChan:   make(chan *models.MessageModel),
		store:       store,
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

func (w *Websocket) RegisterClient(userId string, conn *websocket.Conn) {
	w.clients[userId] = conn
}

func (w *Websocket) UnregisterClient(userId string) {
	delete(w.clients, userId)
}

func (w *Websocket) NewMessage(message *models.MessageModel) {
	w.messageChan <- message
}

func (w *Websocket) Start() {
	go w.listen()
	go w.processMessages()
}

func (w *Websocket) Stop() {
	close(w.messageChan)
	close(w.senderChan)
	close(w.saverChan)
}

func (w *Websocket) listen() {
	for {
		select {
		case msg := <-w.messageChan:
			w.saverChan <- msg
			w.senderChan <- msg
		}
	}
}

func (w *Websocket) processMessages() {
	for {
		select {
		case msg := <-w.saverChan:
			w.storeMessage(msg)
		case msg := <-w.senderChan:
			w.forwardMessage(msg)
		}
	}
}

func (w *Websocket) storeMessage(message *models.MessageModel) {
	if _, err := w.store.Create(message); err != nil {
		fmt.Printf("error saving message: %s\n", err.Error())
	}
}

func (w *Websocket) forwardMessage(message *models.MessageModel) {
	// NOTE: notify receiver
	oneSignal := notification_service.OneSignalInstance
	if oneSignal != nil {
		if err := oneSignal.Notify(message.Receiver, notification_service.NOTIF_TYPE_NEW_MESSAGE); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("dum dum you forgot to initialize one signal")
	}

	if socket, ok := w.clients[message.Receiver]; ok {
		err := socket.WriteJSON(message)
		if err != nil {
			fmt.Printf("error sending message to recipient: %s\n", err.Error())
		}
	} else {
		fmt.Printf("Receiver not online!\n")
	}

	if socket, ok := w.clients[message.Sender]; ok {
		err := socket.WriteJSON(message)
		if err != nil {
			fmt.Printf("error sending message to sender: %s\n", err.Error())
		}
	} else {
		fmt.Printf("Sender not online\n")
	}
}
