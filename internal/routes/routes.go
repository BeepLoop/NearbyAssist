package routes

import (
	"nearbyassist/internal/handlers"
	"nearbyassist/internal/server"
)

func RegisterRoutes(s *server.Server) {
	// File server
	s.Echo.Static("/resource", "store/").Name = "file resource"

	// Routes
	healthHandler := handlers.NewHealthHandler(s)

	rootHandler := handlers.NewHandler(s)
	s.Echo.GET("", rootHandler.HandleBaseRoute).Name = "base route"
	s.Echo.GET("/health", healthHandler.HandleHealthCheck).Name = "base route health check"

	// V1 routes
	v1 := s.Echo.Group("/v1")
	{
		// TODO: Add base route
		v1.GET("/health", healthHandler.HandleHealthCheck).Name = "v1 route health check"
		v1.GET("", rootHandler.HandleV1BaseRoute).Name = "v1 base route"

		// Auth Routes
		auth := v1.Group("/auth")
		{
			handler := handlers.NewAuthHandler(s)
			auth.GET("/health", healthHandler.HandleHealthCheck).Name = "auth route health check"
			auth.POST("", handler.HandleBaseRoute).Name = "auth base route"

			admin := auth.Group("/admin")
			{
				admin.POST("/login", handler.HandleAdminLogin).Name = "admin login"
				admin.POST("/logout", handler.HandleLogout).Name = "admin logout"
			}

			client := auth.Group("/client")
			{
				client.POST("/register", handler.HandleRegister).Name = "client register"
				client.POST("/login", handler.HandleLogin).Name = "client login"
				client.POST("/logout", handler.HandleLogout).Name = "client logout"
			}
		}

		// User Routes
		user := v1.Group("/users")
		{
			handler := handlers.NewUserHandler(s)
			user.GET("/health", healthHandler.HandleHealthCheck).Name = "user route health check"
			user.GET("", handler.HandleBaseRoute).Name = "user base route"
			user.GET("/count", handler.HandleCount).Name = "get number of users"
			user.GET("/:userId", handler.HandleGetUser).Name = "get user details"
		}

		// Vendor Routes
		vendor := v1.Group("/vendors")
		{
			handler := handlers.NewVendorHandler(s)
			vendor.GET("/health", healthHandler.HandleHealthCheck).Name = "vendor route health check"
			vendor.GET("", handler.HandleBaseRoute).Name = "vendor base route"
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
			category.GET("", handler.HandleCategories).Name = "get all categories"
		}

		// Services Routes
		service := v1.Group("/services")
		{
			handler := handlers.NewServiceHandler(s)
			service.GET("/health", healthHandler.HandleHealthCheck).Name = "service route health check"
			service.GET("", handler.HandleGetServices).Name = "get all services"
			service.PUT("/register", handler.HandleRegisterService).Name = "register service"
			service.GET("/search", handler.HandleSearchService).Name = "search service"
			service.GET("/:serviceId", handler.HandleGetDetails).Name = "get service details"
			service.GET("/vendor/:vendorId", handler.HandleGetByVendor).Name = "get owner services"
		}

		// Complaint Routes
		complaint := v1.Group("/complaints")
		{
			handler := handlers.NewComplaintServer(s)
			complaint.GET("/health", healthHandler.HandleHealthCheck).Name = "complaint route health check"
			complaint.GET("", handler.HandleBaseRoute).Name = "complaint base route"
			complaint.GET("/count", handler.HandleCount).Name = "get number of complaints"
			complaint.PUT("/create", handler.HandleNewComplaint).Name = "file a complaint"
		}

		// Transaction Routes
		transaction := v1.Group("/transactions")
		{
			handler := handlers.NewTransactionHandler(s)
			transaction.GET("/health", healthHandler.HandleHealthCheck).Name = "health check for transactions route"
			transaction.GET("", handler.HandleBaseRoute).Name = "transaction base route"
			transaction.PUT("/create", handler.HandleNewTransaction).Name = "create new transaction"
			transaction.GET("/count", handler.HandleCount).Name = "number of transactions, takes in a filter for status"
			// TODO: maybe refactor this to be base route that takes in the following
			// userId = can be a client or vendor ID
			// Filter = view transactions as client or vendor
			// Status = transaction status (see transaction model for valid status)
			transaction.GET("/ongoing/:userId", handler.HandleOngoingTransaction).Name = "get ongoing transaction of given userId. can take a filter = client | vendor"
			transaction.GET("/history/:userId", handler.HandleHistory).Name = "get transaction history, can take in a filter = client | vendor"
		}

		// Application Routes
		application := v1.Group("/application")
		{
			handler := handlers.NewApplicationHandler(s)
			application.GET("/health", healthHandler.HandleHealthCheck).Name = "application route health check"
			application.GET("", handler.HandleGetApplications).Name = "get all vendor applications"
			application.PUT("", handler.HandleNewApplication).Name = "vendor application"
			application.GET("/count", handler.HandleCount).Name = "get number of vendor applications"
			application.PATCH("/approve/:applicationId", handler.HandleApprove).Name = "approve vendor application"
			application.PATCH("/reject/:applicationId", handler.HandleReject).Name = "reject vendor application"
		}

		// Review Routes
		review := v1.Group("/reviews")
		{
			handler := handlers.NewReviewHandler(s)
			review.GET("/health", healthHandler.HandleHealthCheck).Name = "review route health check"
			review.GET("", handler.HandleBaseRoute).Name = "review base route"
			review.GET("/:reviewId", handler.HandleGetReview).Name = "get review details"
			review.PUT("/create", handler.HandleNewReview).Name = "post a review"
			review.GET("/service/:serviceId", handler.HandleServiceReview).Name = "get reviews by service"
		}

		// Upload Routes
		upload := v1.Group("/upload")
		{
			handler := handlers.NewUploadHandler(s)
			upload.GET("/health", healthHandler.HandleHealthCheck).Name = "upload route health check"
			upload.GET("", handler.HandleBaseRoute).Name = "upload base route"
			upload.PUT("/service", handler.HandleNewServicePhoto).Name = "upload service image"
			upload.PUT("/proof", handler.HandleNewProofPhoto).Name = "upload vendor application proof"
		}

		// Chat Routes
		chat := v1.Group("/chat")
		{
			handler := handlers.NewChatHandler(s)
			chat.GET("/health", healthHandler.HandleHealthCheck).Name = "message route health check"
			chat.GET("", handler.HandleBaseRoute).Name = "chat base route"
			chat.GET("/messages", handler.HandleGetMessages).Name = "get messages between sender and receiver"
			chat.GET("/ws", handler.HandleWebsocket).Name = "websocket route for chat"
			chat.GET("/conversations", handler.HandleGetConversations).Name = "get all users you chatted with"
		}
	}
}
