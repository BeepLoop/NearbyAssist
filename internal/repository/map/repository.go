package map_repo

import "nearbyassist/internal/models"

type MapRepository interface {
	GetAllByTag(tag string) ([]*models.ServiceModel, error)
}
