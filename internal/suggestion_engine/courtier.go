package suggestion_engine

import (
	"math/rand"
	"nearbyassist/internal/models"
)

type Courtier struct{}

func NewCourtier() *Courtier {
	return &Courtier{}
}

func (c *Courtier) GenerateSuggestability(service *models.ServiceSearchResult) (float32, error) {
	rng := rand.Float32()
	return rng, nil
}
