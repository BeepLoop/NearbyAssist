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
		v1.Use(middleware.CheckAuth(s.Auth))

		v1.GET("/health", healthHandler.HandleHealthCheck)
		v1.GET("", rootHandler.HandleV1BaseRoute)

		// Admin only routes
		admin := v1.Group("/admin")
		handleAdminRoutes(admin, s)

		// Public routes
		public := v1.Group("/public")
		{
			vendor := public.Group("/vendors")
			{
				handler := handlers.NewVendorHandler(s)
				vendor.GET("", handler.HandleBaseRoute)
				vendor.GET("/:vendorId", handler.HandleGetVendor)
			}

			tags := public.Group("/tags")
			{
				handler := handlers.NewTagHandler(s)
				tags.GET("", handler.HandleGetTags)
			}

			service := public.Group("/services")
			{
				handler := handlers.NewServiceHandler(s)
				service.GET("", handler.HandleGetServices)
				service.POST("", handler.HandleRegisterService)
				service.GET("/search", handler.HandleSearchService)
				service.GET("/:serviceId", handler.HandleGetDetails)
				service.PUT("/:serviceId", handler.HandleUpdateService)
				service.DELETE("/:serviceId", handler.HandleDeleteService)
				service.GET("/vendor/:vendorId", handler.HandleGetByVendor)
				service.GET("/route/:serviceId", handler.HandleFindRoute)
			}

			complaint := public.Group("/complaints")
			{
				handler := handlers.NewComplaintHandler(s)
				// complaint.GET("", handler.HandleBaseRoute)
				// complaint.POST("", handler.HandleNewComplaint)
				complaint.POST("/system", handler.HandleSystemComplaint)
				complaint.POST("/vendor", handler.HandleVendorComplaint)
			}

			transaction := public.Group("/transactions")
			{
				handler := handlers.NewTransactionHandler(s)
				transaction.GET("", handler.HandleBaseRoute)
				transaction.POST("", handler.HandleNewTransaction)
				transaction.POST("/complete/:transactionId", handler.HandleCompleteTransaction)
				// TODO: maybepublic.factor this to be basev1.ute that takes in the following
				// userId = can be a client or vendor ID
				// Filter = view transactions as client or vendor
				// Status = transaction status (see transaction model for valid status)
				transaction.GET("/ongoing", handler.HandleOngoingTransaction)
				transaction.GET("/history", handler.HandleHistory)
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
				chat.GET("/messages/:otherUserId", handler.HandleGetMessages)
				chat.GET("/ws", handler.HandleWebsocket)
				chat.GET("/conversations", handler.HandleGetConversations)
			}

			verification := public.Group("/verification")
			{
				handler := handlers.NewVerificationHandler(s)
				verification.POST("/identity", handler.HandleVerifyIdentity)
			}
		}
	}

	// websocket route
	// NOTE: this route is separated because it is not possible to pass
	// headers to connection request, thus unable to authenticate the user.
	// Instead, access token is passed as a query parameter
	ws := s.Echo.Group("/chat")
	{
		handler := handlers.NewChatHandler(s)

		ws.GET("/ws", handler.HandleWebsocket)
	}
}
