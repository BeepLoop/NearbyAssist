package routes

import (
	"nearbyassist/internal/handlers"
	"nearbyassist/internal/middleware"
	"nearbyassist/internal/server"
)

func RegisterRoutes(s *server.Server) {
	// File server
	s.Echo.Static("/resource", "store/")

	// Routes
	healthHandler := handlers.NewHealthHandler(s)

	rootHandler := handlers.NewHandler(s)
	s.Echo.RouteNotFound("/*", rootHandler.HandleUnknownRoute)
	s.Echo.GET("", rootHandler.HandleBaseRoute)
	s.Echo.GET("/health", healthHandler.HandleHealthCheck)

	// Auth Routes
	auth := s.Echo.Group("/auth")
	{
		handler := handlers.NewAuthHandler(s)

		auth.GET("/health", healthHandler.HandleHealthCheck)
		auth.GET("", handler.HandleBaseRoute)
		auth.POST("/refresh", handler.HandleTokenRefresh)
		auth.POST("/admin/login", handler.HandleAdminLogin)
		auth.POST("/client/login", handler.HandleLogin)
		auth.POST("/logout", handler.HandleLogout)
	}

	// V1 routes
	v1 := s.Echo.Group("/v1")
	{
		v1.GET("/health", healthHandler.HandleHealthCheck)
		v1.GET("", rootHandler.HandleV1BaseRoute)

		// Admin only routes
		admin := s.Echo.Group("/admin")
		handleAdminRoutes(admin, s)

		// Public routes
		public := v1.Group("/public")
		{
			public.Use(middleware.CheckAuth(s.Auth))

			user := public.Group("/users")
			{
				handler := handlers.NewUserHandler(s)
				user.GET("", handler.HandleBaseRoute)
				user.GET("/:userId", handler.HandleGetUser)
			}

			vendor := public.Group("/vendors")
			{
				handler := handlers.NewVendorHandler(s)
				vendor.GET("", handler.HandleBaseRoute)
				vendor.GET("/:vendorId", handler.HandleGetVendor)
			}

			category := public.Group("/category")
			{
				handler := handlers.NewCategoryHandler(s)
				category.GET("", handler.HandleCategories)
			}

			service := public.Group("/services")
			{
				handler := handlers.NewServiceHandler(s)
				service.GET("", handler.HandleGetServices)
				service.POST("", handler.HandleRegisterService)
				service.GET("/search", handler.HandleSearchService)
				service.GET("/:serviceId", handler.HandleGetDetails)
				service.GET("/vendor/:vendorId", handler.HandleGetByVendor)
				service.GET("/route/:serviceId", handler.HandleFindRoute)
			}

			complaint := public.Group("/complaints")
			{
				handler := handlers.NewComplaintServer(s)
				complaint.GET("", handler.HandleBaseRoute)
				complaint.POST("", handler.HandleNewComplaint)
			}

			transaction := public.Group("/transactions")
			{
				handler := handlers.NewTransactionHandler(s)
				transaction.GET("", handler.HandleBaseRoute)
				transaction.POST("", handler.HandleNewTransaction)
				// TODO: maybepublic.factor this to be basev1.ute that takes in the following
				// userId = can be a client or vendor ID
				// Filter = view transactions as client or vendor
				// Status = transaction status (see transaction model for valid status)
				transaction.GET("/ongoing/:userId", handler.HandleOngoingTransaction)
				transaction.GET("/history/:userId", handler.HandleHistory)
			}

			application := public.Group("/application")
			{
				handler := handlers.NewApplicationHandler(s)
				application.GET("", handler.HandleGetApplications)
				application.POST("", handler.HandleNewApplication)
			}

			review := public.Group("/reviews")
			{
				handler := handlers.NewReviewHandler(s)
				review.GET("", handler.HandleBaseRoute)
				review.POST("", handler.HandleNewReview)
				review.GET("/:reviewId", handler.HandleGetReview)
				review.GET("/service/:serviceId", handler.HandleServiceReview)
			}

			upload := public.Group("/upload")
			{
				handler := handlers.NewUploadHandler(s)
				upload.GET("", handler.HandleBaseRoute)
				upload.POST("/service", handler.HandleNewServicePhoto)
				upload.POST("/proof", handler.HandleNewProofPhoto)
			}

			chat := public.Group("/chat")
			{
				handler := handlers.NewChatHandler(s)
				chat.GET("", handler.HandleBaseRoute)
				chat.GET("/messages", handler.HandleGetMessages)
				chat.GET("/ws", handler.HandleWebsocket)
				chat.GET("/conversations", handler.HandleGetConversations)
			}
		}
	}
}
