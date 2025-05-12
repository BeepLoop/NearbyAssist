package suggestion_engine

import (
	"nearbyassist/internal/dto"
)

type Engine interface {
	// Returns map with service ID as key and score as value
	GenerateSuggestions(services []dto.GeospatialOperation) (map[string]float32, error)
	SetWeights(Weights)
	GetValues() Weights
}
