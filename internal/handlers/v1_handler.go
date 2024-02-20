package handlers

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/handlers/auth"
	"nearbyassist/internal/handlers/location"

	"github.com/labstack/echo/v4"
)

func RouteHandlerV1(r *echo.Group) {
	r.GET("/health", health.HealthCheck)

	authGroup := r.Group("/auth")
	auth.AuthHandler(authGroup)

	locationGroup := r.Group("/locations")
	location.LocationHandler(locationGroup)
}
