package server

import (
	"nearbyassist/internal/handler/api/application"
	"nearbyassist/internal/handler/api/complaint"
	"nearbyassist/internal/handler/api/e2ee"
	"nearbyassist/internal/handler/api/health"
	"nearbyassist/internal/handler/api/message"
	"nearbyassist/internal/handler/api/qr"
	"nearbyassist/internal/handler/api/review"
	"nearbyassist/internal/handler/api/service"
	"nearbyassist/internal/handler/api/tag"
	"nearbyassist/internal/handler/api/transaction"
	"nearbyassist/internal/handler/api/user"
	"nearbyassist/internal/handler/api/vendor"
	"nearbyassist/internal/handler/api/verification"
	"nearbyassist/internal/middleware"
	application_repo "nearbyassist/internal/repository/application"
	complaint_repo "nearbyassist/internal/repository/complaint"
	e2ee_repo "nearbyassist/internal/repository/e2ee"
	message_repo "nearbyassist/internal/repository/message"
	review_repo "nearbyassist/internal/repository/review"
	saved_service_repo "nearbyassist/internal/repository/saved_service"
	service_repo "nearbyassist/internal/repository/service"
	tag_repo "nearbyassist/internal/repository/tag"
	transaction_repo "nearbyassist/internal/repository/transaction"
	user_repo "nearbyassist/internal/repository/user"
	vendor_repo "nearbyassist/internal/repository/vendor"
	verification_repo "nearbyassist/internal/repository/verification"
	application_service "nearbyassist/internal/service/application"
	complaint_service "nearbyassist/internal/service/complaint"
	e2ee_service "nearbyassist/internal/service/e2ee"
	health_service "nearbyassist/internal/service/health"
	message_service "nearbyassist/internal/service/message"
	qr_service "nearbyassist/internal/service/qr"
	review_service "nearbyassist/internal/service/review"
	"nearbyassist/internal/service/save_service"
	service_service "nearbyassist/internal/service/service"
	tag_service "nearbyassist/internal/service/tag"
	transaction_service "nearbyassist/internal/service/transaction"
	user_service "nearbyassist/internal/service/user"
	vendor_service "nearbyassist/internal/service/vendor"
	verification_service "nearbyassist/internal/service/verification"

	"github.com/labstack/echo/v4"
)

func (s *Server) v1ApiRoutes(v1 *echo.Group) {
	// ===== HEALTH =======
	healthRoute := v1.Group("/health")
	{
		healthService := health_service.NewService(s.DB)
		handler := health.NewHandler(healthService)

		healthRoute.GET("", handler.CheckHealth)

		protected := healthRoute.Group("/protected")
		{
			protected.Use(middleware.CheckAuth(s.JWT))

			protected.GET("", handler.CheckHealth)
		}
	}

	// ===== USER =======
	userRoute := v1.Group("/user")
	{
		userStore := user_repo.NewMysqlUserRepository(s.DB)
		userService := user_service.NewService(userStore, s.Encrypt, s.Hash, s.JWT)
		handler := user.NewHandler(userService)

		userRoute.POST("/login", handler.Login)
		userRoute.POST("/refresh", handler.Refresh)
		userRoute.POST("/logout", handler.Logout, middleware.CheckAuth(s.JWT))

		protected := userRoute.Group("/protected")
		{
			protected.Use(middleware.CheckAuth(s.JWT))

			protected.GET("/me", handler.GetUser)
			protected.GET("/verified", handler.GetUserVerification)
		}
	}

	// ===== TAGS =======
	tagRoute := v1.Group("/tags")
	{
		tagStore := tag_repo.NewMysqlTagRepository(s.DB)
		tagService := tag_service.NewService(tagStore)
		handler := tag.NewHandler(tagService)

		tagRoute.GET("", handler.GetTags)
		tagRoute.GET("/expertise", handler.GetExpertise)
	}

	// ===== VENDOR =======
	vendorRoute := v1.Group("/vendors")
	{
		vendorRoute.Use(middleware.CheckAuth(s.JWT))

		vendorStore := vendor_repo.NewMysqlVendorRepository(s.DB)
		vendorService := vendor_service.NewService(vendorStore, s.Encrypt)
		handler := vendor.NewHandler(vendorService)

		vendorRoute.GET("/:vendorId", handler.GetVendor)
		vendorRoute.GET("/services/:vendorId", handler.GetVendorServiceList)
	}

	// ===== SERVICES =======
	serviceRoute := v1.Group("/services")
	{
		serviceRoute.Use(middleware.CheckAuth(s.JWT))

		serviceStore := service_repo.NewMysqlServiceRepository(s.DB)
		serviceManager := service_service.NewService(
			serviceStore,
			s.Encrypt,
			s.Hash,
			s.JWT,
			s.SuggestionEngine,
			s.RouteEngine,
			s.FS,
		)

		savedServiceStore := saved_service_repo.NewMysqlSavedServiceRepository(s.DB)
		savedServiceService := save_service.NewService(savedServiceStore, serviceStore, s.JWT, s.Encrypt)

		handler := service.NewHandler(serviceManager, savedServiceService)

		serviceRoute.POST("", handler.CreateService)
		serviceRoute.GET("/search", handler.SearchService)
		serviceRoute.GET("/:serviceId", handler.GetService)
		serviceRoute.PUT("/:serviceId", handler.UpdateService)
		serviceRoute.GET("/get-saved", handler.GetSavedServices)
		serviceRoute.POST("/save", handler.SaveService)
		serviceRoute.POST("/unsave", handler.UnsaveService)
		serviceRoute.GET("/route/:serviceId", handler.FindServiceRoute)
	}

	// ===== TRANSACTIONS =======
	transactionRoute := v1.Group("/transactions")
	{
		transactionRoute.Use(middleware.CheckAuth(s.JWT))

		userStore := user_repo.NewMysqlUserRepository(s.DB)
		userService := user_service.NewService(userStore, s.Encrypt, s.Hash, s.JWT)

		transactionStore := transaction_repo.NewMysqlTransactionRepository(s.DB)
		transactionService := transaction_service.NewService(
			transactionStore,
			s.Encrypt,
			s.JWT,
		)
		handler := transaction.NewHandler(transactionService, userService)

		transactionRoute.POST("", handler.CreateTransaction)
		transactionRoute.GET("/:transactionId", handler.GetTransaction)
		transactionRoute.PUT("/cancel/:transactionId", handler.Cancel)
		transactionRoute.PUT("/accept/:transactionId", handler.Accept)
		transactionRoute.PUT("/reject/:transactionId", handler.Reject)
		transactionRoute.GET("/mine", handler.GetUserTransactionList)
		transactionRoute.GET("/recent", handler.GetRecentTransactions)
		transactionRoute.GET("/confirmed", handler.GetConfirmedTransactions)
		transactionRoute.GET("/history", handler.GetTransactionHistory)
		transactionRoute.POST("/complete/:transactionId", handler.CompleteTransaction)
	}

	// ===== APPLICATION =======
	applicationRoute := v1.Group("/applications")
	{
		applicationRoute.Use(middleware.CheckAuth(s.JWT))

		userStore := user_repo.NewMysqlUserRepository(s.DB)
		userService := user_service.NewService(userStore, s.Encrypt, s.Hash, s.JWT)

		applicationStore := application_repo.NewMysqlApplicationRepository(s.DB)
		applicationService := application_service.NewService(applicationStore, s.FS, s.Encrypt, s.JWT)
		handler := application.NewHandler(applicationService, userService)

		applicationRoute.POST("", handler.CreateApplication)
	}

	// ===== REVIEW =======
	reviewRoute := v1.Group("/reviews")
	{
		reviewRoute.Use(middleware.CheckAuth(s.JWT))

		reviewStore := review_repo.NewMysqlReviewRepository(s.DB)
		reviewService := review_service.NewService(reviewStore, s.Encrypt, s.JWT)
		handler := review.NewHandler(reviewService)

		reviewRoute.POST("", handler.CreateReview)
		reviewRoute.GET("/:reviewId", handler.GetReview)
		reviewRoute.GET("/service/:serviceId", handler.GetServiceReviews)
	}

	// ===== VERIFICATION =======
	verificationRoute := v1.Group("/verification")
	{
		verificationRoute.Use(middleware.CheckAuth(s.JWT))

		useStore := user_repo.NewMysqlUserRepository(s.DB)
		userService := user_service.NewService(useStore, s.Encrypt, s.Hash, s.JWT)

		verificationStore := verification_repo.NewMysqlVerificationRepository(s.DB)
		verificationService := verification_service.NewService(
			verificationStore,
			s.FS,
			s.Encrypt,
			s.JWT,
		)
		handler := verification.NewHandler(verificationService, userService)

		verificationRoute.POST("/identity", handler.CreateIdentityVerification)
	}

	// ===== CHAT =======
	chatRoute := v1.Group("/chat")
	{
		messageStore := message_repo.NewMysqlChatRepository(s.DB)
		messageService := message_service.NewService(messageStore, s.WS, s.Encrypt, s.JWT)
		handler := message.NewHandler(messageService)

		// NOTE: this route is separated because it is not possible to pass
		// headers to connection request, thus unable to authenticate the user.
		// Instead, access token is passed as a query parameter
		chatRoute.GET("/ws", handler.ConnectWebsocket)

		chatRoute.GET("/messages/:otherUserId", handler.GetMessages, middleware.CheckAuth(s.JWT))
		chatRoute.GET("/conversations", handler.GetConversationList, middleware.CheckAuth(s.JWT))
		chatRoute.POST("/send", handler.SendMessage, middleware.CheckAuth(s.JWT))
	}

	// ===== COMPLAINT =======
	complaintRoute := v1.Group("/complaints")
	{

		complaintStore := complaint_repo.NewMysqlComplaintRepository(s.DB)
		complaintService := complaint_service.NewService(complaintStore, s.FS, s.Encrypt)
		handler := complaint.NewHandler(complaintService)

		complaintRoute.POST("/system", handler.CreateSystemComplaint)
		complaintRoute.POST("/vendor", handler.ReportVendor)
	}

	// ===== E2E Encryption =======
	e2eeRoute := v1.Group("/e2ee")
	{
		e2eeRoute.Use(middleware.CheckAuth(s.JWT))

		e2eeStore := e2ee_repo.NewMysqlE2EERepository(s.DB)
		e2eeService := e2ee_service.NewService(e2eeStore, s.Encrypt, s.JWT)
		handler := e2ee.NewHandler(e2eeService)

		e2eeRoute.POST("", handler.SaveKeys)
		e2eeRoute.GET("/keys", handler.GetKeys)
		e2eeRoute.GET("/key/:userId", handler.GetPublicKey)
	}

	qrRoute := v1.Group("/qr")
	{
		qrRoute.Use(middleware.CheckAuth(s.JWT))

		key, err := s.Encrypt.GetKey()
		if err != nil {
			// NOTE: This should not happen unless encryption key is not properly set
			panic(err.Error())
		}

		qrService := qr_service.NewService(key)
		handler := qr.NewHandler(qrService)

		qrRoute.POST("/generateSignature", handler.SignTransaction)
		qrRoute.POST("/verifySignature", handler.VerifySignature)
	}
}
