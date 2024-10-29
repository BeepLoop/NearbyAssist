package message_service

import (
	"nearbyassist/internal/models"
	message_repo "nearbyassist/internal/repository/message"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/websocket"
	"nearbyassist/internal/utils"

	gorilla_ws "github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Service struct {
	store   message_repo.MessageRepository
	ws      *websocket.Websocket
	encrypt auth.Encryption
	jwt     auth.Authenticator
}

func NewService(store message_repo.MessageRepository, ws *websocket.Websocket, encrypt auth.Encryption, jwt auth.Authenticator) *Service {
	return &Service{store: store, ws: ws, encrypt: encrypt, jwt: jwt}
}

func (s *Service) ConnectWebsocket(c echo.Context) error {
	conn, err := s.ws.Upgrade(c)
	if err != nil {
		return err
	}
	defer conn.Close()

	token := c.QueryParam("token")
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return err
	}

	s.ws.RegisterClient(userId, conn)

	for {
		message := new(models.MessageModel)
		err := conn.ReadJSON(message)
		if err != nil {
			if gorilla_ws.IsCloseError(err, gorilla_ws.CloseNormalClosure, gorilla_ws.CloseGoingAway, gorilla_ws.CloseAbnormalClosure) {
				s.ws.UnregisterClient(userId)

				return nil
			}

			continue
		}

		s.ws.NewMessage(message)
	}
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
