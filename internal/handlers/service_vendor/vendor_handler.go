package service_vendor

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/service_vendor/v1"

	"github.com/labstack/echo/v4"
)

func VendorHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "vendor route health check"

	r.GET("/:vendorId", service_vendor.GetVendor).Name = "get vendor details"
}
