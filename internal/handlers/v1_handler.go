package handlers

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/handlers/auth"
	"nearbyassist/internal/handlers/category"
	"nearbyassist/internal/handlers/complaint"
	"nearbyassist/internal/handlers/message"
	"nearbyassist/internal/handlers/review"
	"nearbyassist/internal/handlers/service"
	"nearbyassist/internal/handlers/service_vendor"
	"nearbyassist/internal/handlers/upload"
	"nearbyassist/internal/handlers/user"

	"github.com/labstack/echo/v4"
)

func RouteHandlerV1(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "v1 route health check"

	authGroup := r.Group("/auth")
	auth.AuthHandler(authGroup)

	serviceGroup := r.Group("/services")
	service.ServiceHandler(serviceGroup)

	userGroup := r.Group("/users")
	user.UserHandler(userGroup)

	messageGroup := r.Group("/messages")
	message.MessageHandler(messageGroup)

	vendorGroup := r.Group("/vendors")
	service_vendor.VendorHandler(vendorGroup)

	categoryGroup := r.Group("/categories")
	category.CategoryHandler(categoryGroup)

	uploadGroup := r.Group("/uploads")
	upload.UploadHandler(uploadGroup)

	complaintGroup := r.Group("/complaints")
	complaint.ComplaintsHandler(complaintGroup)

	reviewGroup := r.Group("/reviews")
	review.ReviewsHandler(reviewGroup)
}
