package message_service

import (
	"fmt"
	"nearbyassist/internal/models"
	message_repo "nearbyassist/internal/repository/message"
	"nearbyassist/internal/service/core"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
)

type Service struct {
	store   message_repo.MessageRepository
	ws      websocket.Socket
	encrypt core.Encryption
	jwt     core.Authenticator
}

func NewService(store message_repo.MessageRepository, ws websocket.Socket, encrypt core.Encryption, jwt core.Authenticator) *Service {
	return &Service{store: store, ws: ws, encrypt: encrypt, jwt: jwt}
}

func (s *Service) GetMessages(bearerToken, otherUserId string) ([]*models.MessageModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	messages, err := s.store.GetMessages(userId, otherUserId)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *Service) GetConversationList(bearerToken string) ([]*models.ConversationModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	conversations, err := s.store.GetConversations(userId)
	if err != nil {
		return nil, err
	}

	for _, conversation := range conversations {
		if plain, err := s.encrypt.DecryptString(conversation.Name); err != nil {
			return nil, err
		} else {
			conversation.Name = plain
		}
	}

	return conversations, nil
}

func (s *Service) SendMessage(message *models.MessageModel) error {
	if _, err := s.store.Create(message); err != nil {
		return err
	}

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewMessageNotification(message.Receiver); err != nil {
		fmt.Println(err.Error())
	}

	event := &websocket.EventModel{
		Type:    websocket.EVT_MSSG,
		Payload: message,
	}

	s.ws.Send(event)

	return nil
}
