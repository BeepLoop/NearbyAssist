package server

import (
	"nearbyassist/internal/middleware"
	handler "nearbyassist/internal/service"
	"nearbyassist/internal/store/application"
	"nearbyassist/internal/store/chat"
	"nearbyassist/internal/store/complaint"
	"nearbyassist/internal/store/e2ee"
	"nearbyassist/internal/store/review"
	"nearbyassist/internal/store/service"
	"nearbyassist/internal/store/tag"
	"nearbyassist/internal/store/transaction"
	"nearbyassist/internal/store/user"
	"nearbyassist/internal/store/vendor"
	"nearbyassist/internal/store/verification"

	"github.com/labstack/echo/v4"
)

func (s *Server) api() {
	api := s.Echo.Group("/api")
	{
		v1 := api.Group("/v1")
		s.v1apiHandler(v1)
	}
}

func (s *Server) v1apiHandler(v1 *echo.Group) {
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

	// ===== USER =======
	userRoute := v1.Group("/user")
	{
		userStore := user.NewMysqlUserStore(s.DB)
		h := handler.NewUserService(userStore, s.Encrypt, s.JWT, s.Hash)

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

		h := handler.NewResourceService(s.Encrypt, s.FS)

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
		h := handler.NewServiceService(serviceStore, s.Encrypt, s.JWT, s.Hash, s.RouteEngine, s.SuggestionEngine)

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
		h := handler.NewApplicationService(appStore, s.JWT, s.Encrypt, s.FS)

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
		h := handler.NewVerificationService(verficationStore, s.JWT, s.Encrypt, s.FS)

		verificationRoute.POST("/identity", h.Create)
	}

	// ===== CHAT =======
	chatRoute := v1.Group("/chat")
	{
		chatStore := chat.NewMysqlChatStore(s.DB)
		h := handler.NewChatService(chatStore, s.JWT, s.WS, s.Encrypt)

		// NOTE: this route is separated because it is not possible to pass
		// headers to connection request, thus unable to authenticate the user.
		// Instead, access token is passed as a query parameter
		chatRoute.GET("/ws", h.Websocket)

		chatRoute.GET("/messages/:otherUserId", h.GetMessages, middleware.CheckAuth(s.JWT))
		chatRoute.GET("/conversations", h.GetConversations, middleware.CheckAuth(s.JWT))
	}

	// ===== COMPLAINT =======
	complaintRoute := v1.Group("/complaints")
	{
		complaintStore := complaint.NewMysqlComplaintStore(s.DB)
		h := handler.NewComplaintService(complaintStore, s.Encrypt, s.FS)

		complaintRoute.POST("/system", h.SystemComplaint)
		complaintRoute.POST("/vendor", h.VendorComplaint)
	}

	// ===== E2EE =======
	e2eeRoute := v1.Group("/e2ee")
	{
		e2eeRoute.Use(middleware.CheckAuth(s.JWT))

		e2eeStore := e2ee.NewMysqlE2EEStore(s.DB)
		h := handler.NewE2EEService(e2eeStore, s.JWT, s.Encrypt)

		e2eeRoute.POST("", h.SaveKeys)
		e2eeRoute.GET("/keys", h.GetKeys)
		e2eeRoute.GET("/key/:userId", h.GetPublicKey)
	}
}
