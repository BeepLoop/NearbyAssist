package suggestion_engine

import "nearbyassist/internal/models"

type Engine interface {
	GenerateSuggestability(service *models.ServiceSearchResult) (float32, error)
}
