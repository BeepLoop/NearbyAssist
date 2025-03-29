package notification_service

import (
	"nearbyassist/internal/models"
	notification_repo "nearbyassist/internal/repository/notification"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
)

type Service struct {
	store   notification_repo.NotificationRepository
	encrypt core.Encryption
	jwt     core.Authenticator
}

func NewService(store notification_repo.NotificationRepository, encrypt core.Encryption, jwt core.Authenticator) *Service {
	return &Service{
		store:   store,
		encrypt: encrypt,
		jwt:     jwt,
	}
}

func (s *Service) GetUserNotifications(bearerToken string, status string) ([]*models.NotificationModel, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return nil, err
	}

	notifications := make([]*models.NotificationModel, 0)

	if status == "unread" {
		if notifs, err := s.store.GetAllUnreadByRecipient(userId); err != nil {
			return nil, err
		} else {
			notifications = notifs
		}
	} else {
		if notifs, err := s.store.GetAllByRecipient(userId); err != nil {
			return nil, err
		} else {
			notifications = notifs
		}
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

		notification.IsRead = notification.ReadAt.Valid
	}

	return notifications, nil
}

func (s *Service) ReadNotification(notificationId string) error {
	return s.store.UpdateRead(notificationId)
}

func (s *Service) SendNotification() {
	return
}
