package notification_repo

import "nearbyassist/internal/models"

type NotificationRepository interface {
	Create(data *models.NotificationModel) error
	FindById(id string) (*models.NotificationModel, error)
	UpdateRead(id string) error

	GetAllUnreadByRecipient(id string) ([]*models.NotificationModel, error)
}
