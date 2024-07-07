package routes

import (
	"nearbyassist/internal/handlers"
	"nearbyassist/internal/middleware"
	"nearbyassist/internal/server"

	"github.com/labstack/echo/v4"
)

func handleAdminRoutes(r *echo.Group, s *server.Server) {
	r.Use(middleware.CheckRole(s.Auth))

	management := r.Group("/management")
	{
		handler := handlers.NewAdminHandler(s)

		management.POST("/staff", handler.HandleRegisterStaff)
	}

	user := r.Group("/users")
	{
		handler := handlers.NewUserHandler(s)

		user.GET("", handler.HandleBaseRoute)
		user.GET("/count", handler.HandleCount)
		user.GET("/:userId", handler.HandleGetUser)
	}

	vendor := r.Group("/vendor")
	{
		handler := handlers.NewVendorHandler(s)

		vendor.GET("/count", handler.HandleCount)
		vendor.PUT("/restrict/:vendorId", handler.HandleRestrict)
		vendor.PUT("/unrestrict/:vendorId", handler.HandleUnrestrict)
	}

	application := r.Group("/application")
	{
		handler := handlers.NewApplicationHandler(s)

		application.GET("", handler.HandleGetApplications)
		application.GET("/count", handler.HandleCount)
		application.PUT("/approve/:applicationId", handler.HandleApprove)
		application.PUT("/reject/:applicationId", handler.HandleReject)
	}

	services := r.Group("/services")
	{
		handler := handlers.NewServiceHandler(s)

		services.GET("/count", handler.HandleCount)
	}

	transaction := r.Group("/transactions")
	{
		handler := handlers.NewTransactionHandler(s)

		transaction.GET("/count", handler.HandleCount)
		transaction.GET("/:transactionId", handler.HandleGetTransaction)
	}

	complaint := r.Group("/complaints")
	{
		handler := handlers.NewComplaintHandler(s)

		system := complaint.Group("/system")
		{
			system.GET("", handler.HandleGetSystemComplaint)
			system.GET("/:complaintId", handler.HandleGetSystemComplaintById)
			system.GET("/count", handler.HandleSystemComplaintCount)
		}
	}

	verification := r.Group("/verification")
	{
		handler := handlers.NewVerificationHandler(s)

		identity := verification.Group("/identity")
		{
			identity.GET("", handler.HandleGetAllIdentityVerification)
			identity.GET("/:verificationId", handler.HandleGetIdentityVerification)
		}
	}
}
