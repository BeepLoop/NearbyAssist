package routing_engine

import "nearbyassist/internal/models"

type PolylineCode string

type Engine interface {
	FindRoute(origin, destination *models.Location) (PolylineCode, error)
}
