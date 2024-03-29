package category

import (
	"nearbyassist/internal/controller/category/v1"
	"nearbyassist/internal/controller/health"

	"github.com/labstack/echo/v4"
)

func CategoryHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "category route health check"

	r.GET("/", category.GetCategories).Name = "get all categories"
}
