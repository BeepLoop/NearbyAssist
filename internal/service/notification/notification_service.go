package notification_service

import (
	"nearbyassist/internal/models"
	notification_repo "nearbyassist/internal/repository/notification"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/utils"
)

type Service struct {
	store   notification_repo.NotificationRepository
	encrypt auth.Encryption
	jwt     auth.Authenticator
}

func NewService(store notification_repo.NotificationRepository, encrypt auth.Encryption, jwt auth.Authenticator) *Service {
	return &Service{
		store:   store,
		encrypt: encrypt,
		jwt:     jwt,
	}
}

func (s *Service) GetUnreadNotifications(bearerToken string) ([]*models.NotificationModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	notifications, err := s.store.GetAllUnreadByRecipient(userId)
	if err != nil {
		return nil, err
	}

	for _, notification := range notifications {
		if plainText, err := s.encrypt.DecryptString(notification.Title); err != nil {
			return nil, err
		} else {
			notification.Title = plainText
		}

		if plainText, err := s.encrypt.DecryptString(notification.Content); err != nil {
			return nil, err
		} else {
			notification.Content = plainText
		}
	}

	return notifications, nil
}

func (s *Service) ReadNotification(notificationId string) error {
	return s.store.UpdateRead(notificationId)
}
