package message_service

import (
	"fmt"
	"nearbyassist/internal/models"
	message_repo "nearbyassist/internal/repository/message"
	user_repo "nearbyassist/internal/repository/user"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"
)

type Service struct {
	userStore    user_repo.UserRepository
	messageStore message_repo.MessageRepository
	ws           websocket.Socket
	encrypt      core.Encryption
	jwt          core.Authenticator
}

func NewService(
	userStore user_repo.UserRepository,
	messageStore message_repo.MessageRepository,
	ws websocket.Socket,
	encrypt core.Encryption,
	jwt core.Authenticator,
) *Service {
	return &Service{
		userStore:    userStore,
		messageStore: messageStore,
		ws:           ws,
		encrypt:      encrypt,
		jwt:          jwt,
	}
}

func (s *Service) GetMessages(bearerToken, otherUserId string) ([]*models.MessageModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	messages, err := s.messageStore.GetMessages(userId, otherUserId)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *Service) MarkSeen(bearerToken string, messageIds []string) error {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	for _, id := range messageIds {
		message, err := s.messageStore.FindById(id)
		if err != nil {
			fmt.Println("error retrieving message with id: ", id, ", error: ", err.Error())
			continue
		}

		if message.Sender == userId {
			continue
		}

		if message.Seen {
			continue
		}

		if err := s.messageStore.MarkSeen(id); err != nil {
			fmt.Println("error marking seen. error: ", err.Error())
			continue
		}
	}

	return nil
}

func (s *Service) GetConversationList(bearerToken string) ([]*models.ConversationModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	conversations, err := s.messageStore.GetConversations(userId)
	if err != nil {
		return nil, err
	}

	for _, conversation := range conversations {
		conversation.Name = utils.Must(s.encrypt.DecryptString(conversation.Name))
	}

	return conversations, nil
}

func (s *Service) SendMessage(message *models.MessageModel) error {
	if _, err := s.messageStore.Create(message); err != nil {
		return err
	}

	sender, err := s.userStore.FindById(message.Sender)
	if err != nil {
		return err
	}

	receiver, err := s.userStore.FindById(message.Receiver)
	if err != nil {
		return err
	}

	messagePayload := response.MessageEvent{
		Id: message.Id,
		Sender: response.MessageUser{
			Id:       sender.Id,
			Name:     utils.Must(s.encrypt.DecryptString(sender.Name)),
			ImageURL: sender.ImageUrl,
		},
		Receiver: response.MessageUser{
			Id:       receiver.Id,
			Name:     utils.Must(s.encrypt.DecryptString(receiver.Name)),
			ImageURL: receiver.ImageUrl,
		},
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
		Seen:      message.Seen,
	}

	event := &websocket.EventModel{
		ReceiverId: message.Receiver,
		Type:       websocket.EVT_MSSG,
		Payload:    messagePayload,
	}

	s.ws.Send(event)

	oneSignal := notification_service.MustGetInstance()
	if err := oneSignal.NewMessageNotification(message.Receiver); err != nil {
		fmt.Println(err.Error())
	}

	return nil
}
