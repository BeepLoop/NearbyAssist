package service

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/service/v1"

	"github.com/labstack/echo/v4"
)

func ServiceHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck)

	r.GET("/", service.GetServices)

	r.POST("/register", service.RegisterService)

	r.GET("/search", service.SearchService)

	r.GET(":serviceId", service.GetServiceDetails)
}
