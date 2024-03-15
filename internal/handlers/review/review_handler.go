package review

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/review/v1"

	"github.com/labstack/echo/v4"
)

func ReviewsHandler(r *echo.Group) {

	r.GET("/health", health.HealthCheck)
	r.POST("/create", review.CreateReview)
	r.GET(":vendorId", review.VendorReview)
}
