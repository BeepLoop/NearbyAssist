package application

import "nearbyassist/internal/models"

type ApplicationStore interface {
	Create(data *models.ApplicationModel) (string, error)
}
