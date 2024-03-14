package review

import (
	"nearbyassist/internal/controller/health"

	"github.com/labstack/echo/v4"
)

func ReviewsHandler(r *echo.Group) {

    r.GET("/health", health.HealthCheck)
}
