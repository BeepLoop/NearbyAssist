package route_engine

import "nearbyassist/internal/models"

type PolylineCode string

type Engine interface {
	GetPolyline(origin, destination *models.GeoSpatialModel) (PolylineCode, error)
	GetDistance(origin, destination *models.GeoSpatialModel) (float32, error)
}
