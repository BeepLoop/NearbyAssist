package server

import (
	"nearbyassist/internal/middleware"
	"nearbyassist/internal/service"
	"nearbyassist/internal/store/admin"
	"nearbyassist/internal/store/analytics"
	"nearbyassist/internal/store/application"
	"nearbyassist/internal/store/chat"
	"nearbyassist/internal/store/complaint"
	"nearbyassist/internal/store/review"
	"nearbyassist/internal/store/service"
	"nearbyassist/internal/store/tag"
	"nearbyassist/internal/store/transaction"
	"nearbyassist/internal/store/user"
	"nearbyassist/internal/store/vendor"
	"nearbyassist/internal/store/verification"
)

func (s *Server) routes() {
	v1 := s.Echo.Group("/api/v1")
	{
		// ===== HEALTH =======
		healthRoute := v1.Group("/health")
		{
			h := handler.NewHealthService()
			healthRoute.GET("", h.BaseRoute)

			protected := healthRoute.Group("/protected")
			{
				protected.Use(middleware.CheckAuth(s.JWT))

				protected.GET("", h.BaseRoute)
			}
		}

		// ===== ADMIN =======
		adminRoute := v1.Group("/admin")
		{
			adminStore := admin.NewMysqlAdminStore(s.DB)
			h := handler.NewAdminService(adminStore, s.Encrypt, s.JWT)

			adminRoute.GET("", h.BaseRoute)
			adminRoute.POST("/login", h.Login)
			adminRoute.POST("/refresh", h.Refresh)
			adminRoute.POST("/logout", h.Logout, middleware.CheckAuth(s.JWT))

			management := adminRoute.Group("/management")
			{
				management.Use(middleware.CheckAuth(s.JWT))

				analyticStore := analytics.NewMysqlAnalyticsStore(s.DB)
				h := handler.NewAnalyticsService(analyticStore)

				management.GET("/analytics", h.Analytics)
			}
		}

		// ===== USER =======
		userRoute := v1.Group("/user")
		{
			userStore := user.NewMysqlUserStore(s.DB)
			h := handler.NewUserService(userStore, s.Encrypt, s.JWT)

			userRoute.POST("/login", h.Login)
			userRoute.POST("/refresh", h.Refresh)
			userRoute.POST("/logout", h.Logout, middleware.CheckAuth(s.JWT))

			protected := userRoute.Group("/protected")
			{
				protected.Use(middleware.CheckAuth(s.JWT))

				protected.GET("/me", h.BaseRoute)
				protected.GET("/verified", h.Verified)
			}
		}

		// ===== RESOURCE =======
		resourceRoute := v1.Group("/resource")
		{
			resourceRoute.Use(middleware.CheckAuth(s.JWT))
			resourceRoute.Use(middleware.CheckRole(s.JWT))

			h := handler.NewResourceService(s.Encrypt)

			resourceRoute.GET("/:path", h.GetFile)
		}

		// ===== TAGS =======
		tagRoute := v1.Group("/tags")
		{
			tagStore := tag.NewMysqlTagStore(s.DB)
			h := handler.NewTagService(tagStore)

			tagRoute.GET("", h.BaseRoute)
		}

		// ===== VENDOR =======
		vendorRoute := v1.Group("/vendors")
		{
			vendorRoute.Use(middleware.CheckAuth(s.JWT))

			vendorStore := vendor.NewMysqlVendorStore(s.DB)
			h := handler.NewVendorService(vendorStore)

			vendorRoute.GET("/:vendorId", h.GetVendor)
		}

		// ===== SERVICES =======
		serviceRoute := v1.Group("/services")
		{
			serviceRoute.Use(middleware.CheckAuth(s.JWT))

			serviceStore := service.NewMysqlServiceStore(s.DB)
			h := handler.NewServiceService(serviceStore, s.Encrypt, s.JWT, s.RouteEngine, s.SuggestionEngine)

			serviceRoute.GET("", h.GetServices)
			serviceRoute.POST("", h.Create)
			serviceRoute.GET("/search", h.Search)
			serviceRoute.GET("/:serviceId", h.GetService)
			serviceRoute.PUT("/:serviceId", h.Update)
			serviceRoute.DELETE("/:serviceId", h.Delete)
			serviceRoute.GET("/vendor/:vendorId", h.GetVendorServices)
			serviceRoute.GET("/route/:serviceId", h.FindRoute)
		}

		// ===== TRANSACTIONS =======
		transactionRoute := v1.Group("/transactions")
		{
			transactionRoute.Use(middleware.CheckAuth(s.JWT))

			transactionStore := transaction.NewMysqlTransactionStore(s.DB)
			h := handler.NewTransactionService(transactionStore, s.JWT, s.Encrypt)

			transactionRoute.POST("", h.Create)
			transactionRoute.GET("", h.GetAll)
			transactionRoute.GET("/mine", h.GetMyTransactions)
			transactionRoute.GET("/ongoing", h.GetOngoing)
			transactionRoute.GET("/history", h.GetHistory)
			transactionRoute.POST("/complete/:transactionId", h.Complete)
		}

		// ===== APPLICATION =======
		applicationRoute := v1.Group("/applications")
		{
			applicationRoute.Use(middleware.CheckAuth(s.JWT))

			appStore := application.NewMysqlApplicationStore(s.DB)
			h := handler.NewApplicationService(appStore)

			applicationRoute.POST("", h.Create)
		}

		// ===== REVIEW =======
		reviewRoute := v1.Group("/reviews")
		{
			reviewRoute.Use(middleware.CheckAuth(s.JWT))

			reviewStore := review.NewMysqlReviewStore(s.DB)
			h := handler.NewReviewService(reviewStore, s.JWT)

			reviewRoute.POST("", h.Create)
			reviewRoute.GET("/:reviewId", h.GetById)
			reviewRoute.GET("/service/:serviceId", h.GetByService)
		}

		// ===== VERIFICATION =======
		verificationRoute := v1.Group("/verification")
		{
			verificationRoute.Use(middleware.CheckAuth(s.JWT))

			verficationStore := verification.NewMysqlVerificationStore(s.DB)
			h := handler.NewVerificationService(verficationStore, s.Encrypt, s.Storage)

			verificationRoute.POST("/identity", h.Create)
		}

		// ===== CHAT =======
		chatRoute := v1.Group("/chat")
		{
			chatStore := chat.NewMysqlChatStore(s.DB)
			h := handler.NewChatService(chatStore, s.JWT, *s.Websocket, s.Encrypt)

			// NOTE: this route is separated because it is not possible to pass
			// headers to connection request, thus unable to authenticate the user.
			// Instead, access token is passed as a query parameter
			chatRoute.GET("/ws", h.Websocket)

			chatRoute.GET("/messages/:otherUserId", h.GetMessages, middleware.CheckAuth(s.JWT))
			chatRoute.GET("/conversations", h.GetConversations, middleware.CheckAuth(s.JWT))
		}

		// ===== CHAT =======
		complaintRoute := v1.Group("/complaints")
		{
			complaintStore := complaint.NewMysqlComplaintStore(s.DB)
			h := handler.NewComplaintService(complaintStore, s.Encrypt)

			complaintRoute.POST("/system", h.SystemComplaint)
			complaintRoute.POST("/vendor", h.VendorComplaint)
		}
	}
}
