package review

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/review/v1"

	"github.com/labstack/echo/v4"
)

func ReviewsHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "review route health check"

	r.PUT("/create", review.CreateReview).Name = "post a review"
	r.GET("/service/:serviceId", review.ServiceReviews).Name = "get vendor reviews"
}
