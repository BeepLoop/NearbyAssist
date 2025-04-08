package notification_repo

import "nearbyassist/internal/models"

type NotificationRepository interface {
	Create(data *models.NotificationModel) (string, error)
	FindById(id string) (*models.NotificationModel, error)
	UpdateRead(id string) error

	GetAllUnreadByRecipient(id string) ([]*models.NotificationModel, error)
	GetAllByRecipient(id string) ([]*models.NotificationModel, error)
}
