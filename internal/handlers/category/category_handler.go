package category

import (
	"github.com/labstack/echo/v4"
	"nearbyassist/internal/controller/health"
)

func CategoryHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck)
}
