package service

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/service/v1"

	"github.com/labstack/echo/v4"
)

func ServiceHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "service route health check"

	r.GET("/", service.GetServices).Name = "get all services"
	r.POST("/register", service.RegisterService).Name = "register service"
	r.GET("/search", service.SearchService).Name = "search service"
	r.GET(":serviceId", service.GetServiceDetails).Name = "get service details"
    r.GET("/owner/:ownerId", service.GetOwnerServices).Name = "get owner services"
}
