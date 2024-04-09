package routes

import (
	"nearbyassist/internal/handlers"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRoutes(s *server.Server) {
	// Middlewares
	s.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.RemoveTrailingSlash())
	s.Echo.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Custom validator
	s.Echo.Validator = &utils.Validator{Validator: validator.New()}

	// File server
	s.Echo.Static("/resource", "store/").Name = "file resource"

	// Routes
	healthHandler := handlers.NewHealthHandler(s)

	handler := handlers.NewHandler(s)
	s.Echo.GET("/", handler.HandleRoot).Name = "root route"
	s.Echo.GET("/health", healthHandler.HandleHealthCheck).Name = "health check"

	// V1 routes
	v1 := s.Echo.Group("/v1")
	{
		v1.GET("/health", healthHandler.HandleHealthCheck)

		// Auth Routes
		auth := v1.Group("/auth")
		{
			handler := handlers.NewAuthHandler(s)
			auth.GET("/health", healthHandler.HandleHealthCheck).Name = "auth route health check"
			auth.POST("/", handler.HandleAdminLogin).Name = "admin login"
			auth.POST("/login", handler.HandleLogin).Name = "client login"
			auth.POST("/logout", handler.HandleLogout).Name = "client logout"
		}

		// User Routes
		user := v1.Group("/users")
		{
			handler := handlers.NewUserHandler(s)
			user.GET("/health", healthHandler.HandleHealthCheck).Name = "user route health check"
			user.GET("/count", handler.HandleCount).Name = "get number of users"
			user.GET("/:userId", handler.HandleGetUser).Name = "get user details"
		}

		// Vendor Routes
		vendor := v1.Group("/vendors")
		{
			handler := handlers.NewVendorHandler(s)
			vendor.GET("/health", healthHandler.HandleHealthCheck).Name = "vendor route health check"
			vendor.GET("/count", handler.HandleCount).Name = "get number of vendors"
			vendor.GET("/:vendorId", handler.HandleGetVendor).Name = "get vendor details"
			vendor.PATCH("/restrict/:vendorId", handler.HandleRestrict).Name = "restrict vendor"
			vendor.PATCH("/unrestrict/:vendorId", handler.HandleUnrestrict).Name = "unrestrict vendor"
		}

		// Category Routes
		category := v1.Group("/category")
		{
			handler := handlers.NewCategoryHandler(s)
			category.GET("/health", healthHandler.HandleHealthCheck).Name = "category route health check"
			category.GET("/", handler.HandleCategories).Name = "get all categories"
			category.GET("", handler.HandleCategories).Name = "get all categories"
		}

		// Services Routes
		service := v1.Group("/services")
		{
			handler := handlers.NewServiceHandler(s)
			service.GET("/health", healthHandler.HandleHealthCheck).Name = "service route health check"
			service.GET("/", handler.HandleGetServices).Name = "get all services"
			service.PUT("/register", handler.HandleRegisterService).Name = "register service"
			service.GET("/search", handler.HandleSearchService).Name = "search service"
			service.GET("/:serviceId", handler.HandleGetDetails).Name = "get service details"
			service.GET("/owner/:ownerId", handler.HandleGetByOwner).Name = "get owner services"
		}

		// Complaint Routes
		complaint := v1.Group("/complaints")
		{
			handler := handlers.NewComplaintServer(s)
			complaint.GET("/health", healthHandler.HandleHealthCheck).Name = "complaint route health check"
			complaint.GET("/count", handler.HandleCount).Name = "get number of complaints"
			complaint.PUT("/create", handler.HandleNewComplaint).Name = "file a complaint"
		}

		// Transaction Routes
		transaction := v1.Group("/transactions")
		{
			handler := handlers.NewTransactionHandler(s)
			transaction.GET("/health", healthHandler.HandleHealthCheck).Name = "health check for transactions route"
			transaction.PUT("/create", handler.HandleNewTransaction).Name = "create new transaction"
			transaction.GET("/count", handler.HandleCount).Name = "number of transactions, takes in a filter for status"

			ongoing := transaction.Group("/ongoing")
			{
				ongoing.GET("/client/:userId", handler.HandleClientOngoingTransaction).Name = "get client ongoing transactions"
				ongoing.GET("/vendor/:userId", handler.HandleVendorOngoingTransaction).Name = "get vendor ongoing transactions"
			}

			history := transaction.Group("/history")
			{
				history.GET("/client/:userId", handler.HandleClientHistory).Name = "get client transaction history"
				history.GET("/vendor/:userId", handler.HandleVendorHistory).Name = "get vendor transaction history"
			}
		}

		// Application Routes
		application := v1.Group("/application")
		{
			handler := handlers.NewApplicationHandler(s)
			application.GET("/health", healthHandler.HandleHealthCheck)
			application.PUT("", handler.HandleNewApplication).Name = "vendor application"
			application.GET("", handler.HandleGetApplicants).Name = "get all vendor applications"
			application.GET("/count", handler.HandleCount).Name = "get number of vendor applications"
			application.PATCH("/approve/:applicationId", handler.HandleApprove).Name = "approve vendor application"
			application.PATCH("/reject/:applicationId", handler.HandleReject).Name = "reject vendor application"
		}

		// Review Routes
		review := v1.Group("/reviews")
		{
			handler := handlers.NewReviewHandler(s)
			review.GET("/health", healthHandler.HandleHealthCheck).Name = "review route health check"
			review.PUT("/create", handler.HandleNewReview).Name = "post a review"
			review.GET("/service/:serviceId", handler.HandleServiceReview).Name = "get reviews by service"
		}

		// Upload Routes
		upload := v1.Group("/upload")
		{
			handler := handlers.NewUploadHandler(s)
			upload.GET("/health", healthHandler.HandleHealthCheck)
			upload.PUT("/service", handler.HandleNewServicePhoto).Name = "upload service image"
			upload.PUT("/proof", handler.HandleNewProofPhoto).Name = "upload vendor application proof"
		}

		// Chat Routes
		chat := v1.Group("/chat")
		{
			handler := handlers.NewChatHandler(s)
			chat.GET("/health", healthHandler.HandleHealthCheck).Name = "message route health check"
			chat.GET("/messages", handler.HandleGetMessages).Name = "get messages between sender and receiver"
			chat.GET("/ws", handler.HandleWebsocket).Name = "websocket route for chat"
			chat.GET("/conversations", handler.HandleWebsocket).Name = "get all users you chatted with"
		}
	}
}
