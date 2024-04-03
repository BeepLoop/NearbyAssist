package service_vendor

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/service_vendor/v1"

	"github.com/labstack/echo/v4"
)

func VendorHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "vendor route health check"

	r.GET("/count", service_vendor.CountVendors).Name = "get number of vendors"
	r.GET("/vendor/:vendorId", service_vendor.GetVendor).Name = "get vendor details"
	r.PATCH("/restrict/:vendorId", service_vendor.RestrictVendor).Name = "restrict vendor"
	r.PATCH("/unrestrict/:vendorId", service_vendor.UnrestrictVendor).Name = "unrestrict vendor"
	r.PUT("/application", service_vendor.VendorApplication).Name = "vendor application"
	r.GET("/application/count", service_vendor.CountPendingApplications).Name = "get number of vendor applications"
}
