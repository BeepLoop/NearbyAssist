package suggestion_engine

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
)

type Engine interface {
	GenerateSuggestions(services []*models.GeoSpatialSearchResult) ([]*response.SearchResult, error)
}
