package location

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/location/v1"

	"github.com/labstack/echo/v4"
)

func LocationHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck)

	r.GET("/", location.HandleLocations)

	r.POST("/register", location.RegisterLocation)

	r.GET("/vicinity", location.HandleVicinity)
}
