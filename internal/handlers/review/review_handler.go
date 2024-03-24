package review

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/review/v1"

	"github.com/labstack/echo/v4"
)

func ReviewsHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "review route health check"

	r.POST("/create", review.CreateReview).Name = "post a review"
	r.GET(":vendorId", review.VendorReview).Name = "get vendor reviews"
}
