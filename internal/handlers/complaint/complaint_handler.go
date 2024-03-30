package complaint

import (
	"nearbyassist/internal/controller/complaint/v1"
	"nearbyassist/internal/controller/health"

	"github.com/labstack/echo/v4"
)

func ComplaintsHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "complaint route health check"

	r.PUT("/create", complaint.CreateComplaint).Name = "file a complaint"
}
