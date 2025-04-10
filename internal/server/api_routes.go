package server

import (
	"nearbyassist/internal/handler/api/application"
	"nearbyassist/internal/handler/api/booking"
	"nearbyassist/internal/handler/api/complaint"
	"nearbyassist/internal/handler/api/e2ee"
	"nearbyassist/internal/handler/api/expertise"
	"nearbyassist/internal/handler/api/health"
	"nearbyassist/internal/handler/api/message"
	"nearbyassist/internal/handler/api/notification"
	"nearbyassist/internal/handler/api/qr"
	recommendation_handler "nearbyassist/internal/handler/api/recommendation"
	"nearbyassist/internal/handler/api/resource"
	"nearbyassist/internal/handler/api/review"
	"nearbyassist/internal/handler/api/service"
	"nearbyassist/internal/handler/api/tag"
	"nearbyassist/internal/handler/api/user"
	userauth_handler "nearbyassist/internal/handler/api/user_auth"
	"nearbyassist/internal/handler/api/vendor"
	websocket_handler "nearbyassist/internal/handler/api/websocket"
	"nearbyassist/internal/middleware"
	application_repo "nearbyassist/internal/repository/application"
	booking_repo "nearbyassist/internal/repository/booking"
	bug_report_repo "nearbyassist/internal/repository/bug_report"
	e2ee_repo "nearbyassist/internal/repository/e2ee"
	expertise_repo "nearbyassist/internal/repository/expertise"
	message_repo "nearbyassist/internal/repository/message"
	notification_repo "nearbyassist/internal/repository/notification"
	policeclearance_repo "nearbyassist/internal/repository/police_clearance"
	report_user_repo "nearbyassist/internal/repository/report_user"
	review_repo "nearbyassist/internal/repository/review"
	saved_service_repo "nearbyassist/internal/repository/saved_service"
	service_repo "nearbyassist/internal/repository/service"
	supportingimage_repo "nearbyassist/internal/repository/supporting_image"
	tag_repo "nearbyassist/internal/repository/tag"
	user_repo "nearbyassist/internal/repository/user"
	vendor_repo "nearbyassist/internal/repository/vendor"
	verification_repo "nearbyassist/internal/repository/verification"
	application_service "nearbyassist/internal/service/application"
	booking_service "nearbyassist/internal/service/booking"
	complaint_service "nearbyassist/internal/service/complaint"
	e2ee_service "nearbyassist/internal/service/e2ee"
	expertise_service "nearbyassist/internal/service/expertise"
	health_service "nearbyassist/internal/service/health"
	message_service "nearbyassist/internal/service/message"
	notification_service "nearbyassist/internal/service/notification"
	qr_service "nearbyassist/internal/service/qr"
	recommendation_service "nearbyassist/internal/service/recommendation"
	resource_service "nearbyassist/internal/service/resource"
	review_service "nearbyassist/internal/service/review"
	"nearbyassist/internal/service/save_service"
	service_service "nearbyassist/internal/service/service"
	tag_service "nearbyassist/internal/service/tag"
	user_service "nearbyassist/internal/service/user"
	userauth_service "nearbyassist/internal/service/user_auth"
	vendor_service "nearbyassist/internal/service/vendor"
	verification_service "nearbyassist/internal/service/verification"
	"nearbyassist/internal/utils"

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

	// ===== HEALTH =======
	authRoute := v1.Group("/auth")
	{
		userStore := user_repo.NewMysqlUserRepository(s.DB)
		authService := userauth_service.NewService(userStore, s.Encrypt, s.Hash, s.JWT)
		handler := userauth_handler.NewHandler(authService)

		authRoute.POST("/thirdPartyLogin", handler.ThirdPartyLogin)
		authRoute.POST("/refresh", handler.Refresh)
		authRoute.POST("/logout", handler.Logout, middleware.CheckAuth(s.JWT))
	}

	// ===== USER =======
	userRoute := v1.Group("/user")
	{
		userRoute.Use(middleware.CheckAuth(s.JWT))

		userStore := user_repo.NewMysqlUserRepository(s.DB)
		notificationStore := notification_repo.NewMysqlNotificationRepository(s.DB)
		verificationStore := verification_repo.NewMysqlVerificationRepository(s.DB)

		userService := user_service.NewService(userStore, s.Encrypt, s.Hash, s.JWT)
		userVerificationService := verification_service.NewService(
			userStore,
			verificationStore,
			notificationStore,
			s.WS,
			s.FS,
			s.Encrypt,
			s.JWT,
		)

		handler := user.NewHandler(userService, userVerificationService)

		userRoute.GET("", handler.GetUser)
		userRoute.GET("/verify", handler.GetUserVerification)
		userRoute.POST("/verify", handler.RequestIdentityVerification)
		userRoute.POST("/socials", handler.AddSocial)
		userRoute.DELETE("/socials", handler.DeleteSocial)
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

	// ===== Expertise =======
	expertiseRoute := v1.Group("/expertise")
	{
		expertiseRoute.Use(middleware.CheckAuth(s.JWT))

		userStore := user_repo.NewMysqlUserRepository(s.DB)
		vendorStore := vendor_repo.NewMysqlVendorRepository(s.DB)
		expertiseStore := expertise_repo.NewMysqlExpertiseRepository(s.DB)
		supportingImageStore := supportingimage_repo.NewMysqlImplementation(s.DB)

		expertiseService := expertise_service.NewService(userStore, vendorStore, expertiseStore, supportingImageStore, s.FS, s.Encrypt, s.Hash, s.JWT)

		handler := expertise.NewHandler(expertiseService)

		expertiseRoute.POST("/add", handler.AddUserExpertise)
	}

	// ===== VENDOR =======
	vendorRoute := v1.Group("/vendors")
	{
		vendorRoute.Use(middleware.CheckAuth(s.JWT))

		vendorStore := vendor_repo.NewMysqlVendorRepository(s.DB)
		serviceStore := service_repo.NewMysqlServiceRepository(s.DB)

		vendorService := vendor_service.NewService(vendorStore, serviceStore, s.Encrypt, s.Hash)
		resourceService := resource_service.NewService(s.FS, s.Encrypt, s.Hash)

		handler := vendor.NewHandler(vendorService, resourceService)

		vendorRoute.GET("/:vendorId", handler.GetVendor)
		vendorRoute.GET("/services/:vendorId", handler.GetVendorServiceList)
	}

	// ===== SERVICES =======
	serviceRoute := v1.Group("/services")
	{
		serviceRoute.Use(middleware.CheckAuth(s.JWT))

		savedServiceStore := saved_service_repo.NewMysqlSavedServiceRepository(s.DB)
		serviceStore := service_repo.NewMysqlServiceRepository(s.DB)
		vendorStore := vendor_repo.NewMysqlVendorRepository(s.DB)

		serviceManager := service_service.NewService(
			serviceStore,
			vendorStore,
			s.Encrypt,
			s.Hash,
			s.JWT,
			s.SuggestionEngine,
			s.RouteEngine,
			s.FS,
		)
		savedServiceService := save_service.NewService(savedServiceStore, serviceStore, vendorStore, s.JWT, s.Encrypt)
		resourceService := resource_service.NewService(s.FS, s.Encrypt, s.Hash)

		handler := service.NewHandler(serviceManager, savedServiceService, resourceService)

		serviceRoute.POST("", handler.CreateService)
		serviceRoute.GET("/search", handler.SearchService)
		serviceRoute.GET("/:serviceId", handler.GetService)
		serviceRoute.PUT("", handler.UpdateService)
		serviceRoute.DELETE("/deleteImage/:imageId", handler.DeleteImage)
		serviceRoute.POST("/addImage/:serviceId", handler.AddImage)
		serviceRoute.POST("/addExtra", handler.AddExtra)
		serviceRoute.PUT("/editExtra", handler.EditExtra)
		serviceRoute.DELETE("/deleteExtra/:extraId", handler.DeleteExtra)
		serviceRoute.GET("/get-saved", handler.GetSavedServices)
		serviceRoute.POST("/save", handler.SaveService)
		serviceRoute.POST("/unsave", handler.UnsaveService)
		serviceRoute.GET("/route/:serviceId", handler.FindServiceRoute)
	}

	// ===== BOOKINGS =======
	bookingRoute := v1.Group("/bookings")
	{
		bookingRoute.Use(middleware.CheckAuth(s.JWT))

		userStore := user_repo.NewMysqlUserRepository(s.DB)
		bookingStore := booking_repo.NewMysqlBookingRepository(s.DB)
		notifStore := notification_repo.NewMysqlNotificationRepository(s.DB)
		serviceStore := service_repo.NewMysqlServiceRepository(s.DB)
		vendorStore := vendor_repo.NewMysqlVendorRepository(s.DB)

		userService := user_service.NewService(userStore, s.Encrypt, s.Hash, s.JWT)
		serviceService := service_service.NewService(
			serviceStore,
			vendorStore,
			s.Encrypt,
			s.Hash,
			s.JWT,
			s.SuggestionEngine,
			s.RouteEngine,
			s.FS,
		)
		bookingService := booking_service.NewService(
			notifStore,
			bookingStore,
			s.WS,
			s.Encrypt,
			s.JWT,
		)
		qrService := qr_service.NewService(utils.Must((s.Encrypt.GetKey())))

		handler := booking.NewHandler(bookingService, serviceService, userService, qrService)

		bookingRoute.POST("", handler.CreateBooking)
		bookingRoute.GET("/:bookingId", handler.GetBooking)
		bookingRoute.PUT("/cancel", handler.Cancel)
		bookingRoute.PUT("/accept", handler.Accept)
		bookingRoute.PUT("/reject", handler.Reject)
		bookingRoute.GET("/mine", handler.GetUserBookingList)
		bookingRoute.GET("/recent", handler.GetRecentBookings)
		bookingRoute.GET("/confirmed", handler.GetConfirmedBookings)
		bookingRoute.GET("/toReview", handler.GetReviewableBookings)
		bookingRoute.GET("/history", handler.GetBookingHistory)
		bookingRoute.POST("/complete/:bookingId", handler.CompleteBooking)
	}

	// ===== APPLICATION =======
	applicationRoute := v1.Group("/applications")
	{
		applicationRoute.Use(middleware.CheckAuth(s.JWT))

		userStore := user_repo.NewMysqlUserRepository(s.DB)
		userService := user_service.NewService(userStore, s.Encrypt, s.Hash, s.JWT)

		applicationStore := application_repo.NewMysqlApplicationRepository(s.DB)
		supportingImageStore := supportingimage_repo.NewMysqlImplementation(s.DB)
		policeClearanceStore := policeclearance_repo.NewMysqlImplementation(s.DB)
		notificationStore := notification_repo.NewMysqlNotificationRepository(s.DB)

		applicationService := application_service.NewService(
			applicationStore,
			supportingImageStore,
			policeClearanceStore,
			notificationStore,
			s.WS,
			s.FS,
			s.Encrypt,
			s.JWT,
		)

		handler := application.NewHandler(applicationService, userService)

		applicationRoute.POST("", handler.CreateApplication)
	}

	// ===== REVIEW =======
	reviewRoute := v1.Group("/reviews")
	{
		reviewRoute.Use(middleware.CheckAuth(s.JWT))

		bookingStore := booking_repo.NewMysqlBookingRepository(s.DB)
		reviewStore := review_repo.NewMysqlReviewRepository(s.DB)

		reviewService := review_service.NewService(bookingStore, reviewStore, s.Encrypt, s.JWT)

		handler := review.NewHandler(reviewService)

		reviewRoute.POST("", handler.CreateReview)
		reviewRoute.GET("/:reviewId", handler.GetReview)
	}

	// ===== CHAT =======
	chatRoute := v1.Group("/chat")
	{
		chatRoute.Use(middleware.CheckAuth(s.JWT))

		messageStore := message_repo.NewMysqlChatRepository(s.DB)
		messageService := message_service.NewService(messageStore, s.WS, s.Encrypt, s.JWT)
		handler := message.NewHandler(messageService)

		chatRoute.POST("/send", handler.SendMessage)
		chatRoute.GET("/messages/:otherUserId", handler.GetMessages)
		chatRoute.GET("/conversations", handler.GetConversationList)
	}

	// ===== COMPLAINT =======
	complaintRoute := v1.Group("/complaints")
	{

		reportUserStore := report_user_repo.NewMysqlReportUserRepository(s.DB)
		userStore := user_repo.NewMysqlUserRepository(s.DB)
		bugReportStore := bug_report_repo.NewMysqlBugReportRepository(s.DB)
		notifStore := notification_repo.NewMysqlNotificationRepository(s.DB)

		complaintService := complaint_service.NewService(reportUserStore, userStore, bugReportStore, notifStore, s.WS, s.FS, s.Encrypt, s.JWT)
		handler := complaint.NewHandler(complaintService)

		complaintRoute.POST("/system", handler.CreateBugReport)
		complaintRoute.POST("/user", handler.ReportUser)
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

		qrService := qr_service.NewService(utils.Must(s.Encrypt.GetKey()))
		handler := qr.NewHandler(qrService)

		qrRoute.POST("/generateSignature", handler.SignBooking)
		qrRoute.POST("/verifySignature", handler.VerifySignature)
	}

	notificationRoute := v1.Group("/notifications")
	{
		notificationRoute.Use(middleware.CheckAuth(s.JWT))

		notificationStore := notification_repo.NewMysqlNotificationRepository(s.DB)
		notificationService := notification_service.NewService(notificationStore, s.Encrypt, s.JWT)
		handler := notification.NewHandler(notificationService)

		notificationRoute.GET("", handler.GetNotifications)
		notificationRoute.POST("/:notificationId", handler.ReadNotification)
	}

	resourceRoute := v1.Group("/resource")
	{
		resourceService := resource_service.NewService(s.FS, s.Encrypt, s.Hash)
		handler := resource.NewHandler(resourceService)

		resourceRoute.GET("/:path", handler.GetPrivateFile)
		resourceRoute.GET("/public/:path", handler.GetPublicFile)
	}

	websocketRoute := v1.Group("/ws")
	{
		handler := websocket_handler.NewHandler(s.WS, s.JWT)

		// NOTE: this route is separate because I have problems passing JWT from
		// client Authorization header and continuously pass updated token on
		// reconnect. Hard skill issues
		websocketRoute.GET("", handler.Connect)
	}

	recommendationRoute := v1.Group("/recommendations")
	{
		recommendationRoute.Use(middleware.CheckAuth(s.JWT))

		serviceStore := service_repo.NewMysqlServiceRepository(s.DB)
		vendorStore := vendor_repo.NewMysqlVendorRepository(s.DB)

		recommendationService := recommendation_service.NewService(serviceStore, vendorStore, s.Encrypt)
		resourceService := resource_service.NewService(s.FS, s.Encrypt, s.Hash)

		handler := recommendation_handler.NewHandler(recommendationService, resourceService)

		recommendationRoute.GET("", handler.GetRecommendations)
	}
}
