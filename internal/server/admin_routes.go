package server

import (
	accountmanagement "nearbyassist/internal/handler/admin/account_management"
	"nearbyassist/internal/handler/admin/auth"
	"nearbyassist/internal/handler/admin/complaint"
	"nearbyassist/internal/handler/admin/dashboard"
	"nearbyassist/internal/handler/admin/expertise"
	"nearbyassist/internal/handler/admin/invitation"
	map_handler "nearbyassist/internal/handler/admin/map"
	passwordreset_handler "nearbyassist/internal/handler/admin/password_reset"
	"nearbyassist/internal/handler/admin/userManagement"
	application "nearbyassist/internal/handler/admin/vendor_application"
	"nearbyassist/internal/handler/admin/verification"
	"nearbyassist/internal/middleware"
	admin_repo "nearbyassist/internal/repository/admin"
	application_repo "nearbyassist/internal/repository/application"
	bug_report_repo "nearbyassist/internal/repository/bug_report"
	dashboard_repo "nearbyassist/internal/repository/dashboard"
	expertise_repo "nearbyassist/internal/repository/expertise"
	invitation_repo "nearbyassist/internal/repository/invitation"
	map_repo "nearbyassist/internal/repository/map"
	notification_repo "nearbyassist/internal/repository/notification"
	passwordreset_repo "nearbyassist/internal/repository/password_reset"
	policeclearance_repo "nearbyassist/internal/repository/police_clearance"
	report_user_repo "nearbyassist/internal/repository/report_user"
	service_repo "nearbyassist/internal/repository/service"
	supportingimage_repo "nearbyassist/internal/repository/supporting_image"
	tag_repo "nearbyassist/internal/repository/tag"
	user_repo "nearbyassist/internal/repository/user"
	vendor_repo "nearbyassist/internal/repository/vendor"
	verification_repo "nearbyassist/internal/repository/verification"
	adminauth_service "nearbyassist/internal/service/admin_auth"
	application_service "nearbyassist/internal/service/application"
	complaint_service "nearbyassist/internal/service/complaint"
	dashboard_service "nearbyassist/internal/service/dashboard"
	expertise_service "nearbyassist/internal/service/expertise"
	"nearbyassist/internal/service/invite_service"
	map_service "nearbyassist/internal/service/map"
	passwordreset_service "nearbyassist/internal/service/password_reset"
	resource_service "nearbyassist/internal/service/resource"
	tag_service "nearbyassist/internal/service/tag"
	user_service "nearbyassist/internal/service/user"
	"nearbyassist/internal/service/user_management_service"
	vendor_service "nearbyassist/internal/service/vendor"
	verification_service "nearbyassist/internal/service/verification"

	"github.com/labstack/echo/v4"
)

func (s *Server) AdminRoutes(r *echo.Group) {
	adminStore := admin_repo.NewMysqlAdminRepository(s.DB)
	authService := adminauth_service.NewService(adminStore, s.Encrypt, s.Hash)

	authHandler := auth.NewHandler(authService)

	r.GET("/login", authHandler.GetLogin)
	r.POST("/login", authHandler.PostLogin)
	r.POST("/logout", authHandler.PostLogout)

	resetRoute := r.Group("/reset")
	{
		adminStore := admin_repo.NewMysqlAdminRepository(s.DB)
		passwordResetStore := passwordreset_repo.NewMysqlPasswordResetRepository(s.DB)

		passwordResetService := passwordreset_service.NewService(adminStore, passwordResetStore, s.Encrypt, s.Hash)
		handler := passwordreset_handler.NewHandler(passwordResetService)

		resetRoute.POST("", handler.RequestPasswordReset)
		resetRoute.GET("/cp", handler.ChangePassword, middleware.CheckIfShouldChangePass(adminStore))
		resetRoute.POST("/cp", handler.PostChangePassword)
	}

	dashboardRoute := r.Group("/dashboard")
	{
		adminStore := admin_repo.NewMysqlAdminRepository(s.DB)

		dashboardRoute.Use(middleware.CheckSession)
		dashboardRoute.Use(middleware.CheckMustChangePass(adminStore))

		dashboardStore := dashboard_repo.NewMysqlDashboardRepository(s.DB)
		dashboardService := dashboard_service.NewService(dashboardStore)
		dashboardHandler := dashboard.NewHandler(dashboardService)

		dashboardRoute.GET("", dashboardHandler.GetDashboard)
	}

	mapRoute := r.Group("/map")
	{
		adminStore := admin_repo.NewMysqlAdminRepository(s.DB)

		mapRoute.Use(middleware.CheckSession)
		mapRoute.Use(middleware.CheckMustChangePass(adminStore))

		tagStore := tag_repo.NewMysqlTagRepository(s.DB)
		serviceStore := service_repo.NewMysqlServiceRepository(s.DB)

		tagService := tag_service.NewService(tagStore)

		mapStore := map_repo.NewMysqlMapRepository(s.DB)
		mapService := map_service.NewService(mapStore, serviceStore, s.Encrypt)
		mapHandler := map_handler.NewHandler(mapService, tagService)

		mapRoute.GET("", mapHandler.GetMap)
	}

	complaintRoute := r.Group("/complaints")
	{
		adminStore := admin_repo.NewMysqlAdminRepository(s.DB)

		complaintRoute.Use(middleware.CheckSession)
		complaintRoute.Use(middleware.CheckMustChangePass(adminStore))

		reportUserStore := report_user_repo.NewMysqlReportUserRepository(s.DB)
		bugReportStore := bug_report_repo.NewMysqlBugReportRepository(s.DB)
		notifStore := notification_repo.NewMysqlNotificationRepository(s.DB)

		complaintService := complaint_service.NewService(reportUserStore, bugReportStore, notifStore, s.FS, s.Encrypt, s.JWT)
		resourceService := resource_service.NewService(s.FS, s.Encrypt, s.Hash)

		complaintHandler := complaint.NewHandler(complaintService, resourceService)

		complaintRoute.GET("/bugs", complaintHandler.GetBugReports)
		complaintRoute.POST("/bugs/complete", complaintHandler.CompleteBug)
		complaintRoute.GET("/users", complaintHandler.GetReportedUsers)
		complaintRoute.GET("/users/:reportId", complaintHandler.GetReportedUserDetail)
		complaintRoute.POST("/users/close", complaintHandler.CloseUserReport)
	}

	applicationRoute := r.Group("/vendor-applications")
	{
		adminStore := admin_repo.NewMysqlAdminRepository(s.DB)

		applicationRoute.Use(middleware.CheckSession)
		applicationRoute.Use(middleware.CheckMustChangePass(adminStore))

		applicationStore := application_repo.NewMysqlApplicationRepository(s.DB)
		supportingImageStore := supportingimage_repo.NewMysqlImplementation(s.DB)
		policeClearanceStore := policeclearance_repo.NewMysqlImplementation(s.DB)
		notificationStore := notification_repo.NewMysqlNotificationRepository(s.DB)

		applicationService := application_service.NewService(applicationStore, supportingImageStore, policeClearanceStore, notificationStore, s.FS, s.Encrypt, s.JWT)
		resourceService := resource_service.NewService(s.FS, s.Encrypt, s.Hash)

		applicationHandler := application.NewHandler(applicationService, resourceService)

		applicationRoute.GET("", applicationHandler.GetVendorApplication)
		applicationRoute.GET("/:applicationId", applicationHandler.GetVendorApplicationDetails)
		applicationRoute.POST("/accept", applicationHandler.AcceptRequest)
		applicationRoute.POST("/reject", applicationHandler.RejectRequest)
	}

	verificationRoute := r.Group("/verification-requests")
	{
		adminStore := admin_repo.NewMysqlAdminRepository(s.DB)

		verificationRoute.Use(middleware.CheckSession)
		verificationRoute.Use(middleware.CheckMustChangePass(adminStore))

		requestStore := verification_repo.NewMysqlVerificationRepository(s.DB)
		notificationStore := notification_repo.NewMysqlNotificationRepository(s.DB)

		requestService := verification_service.NewService(requestStore, notificationStore, s.FS, s.Encrypt, s.JWT)
		resourceService := resource_service.NewService(s.FS, s.Encrypt, s.Hash)

		requestHandler := verification.NewHandler(requestService, resourceService)

		verificationRoute.GET("", requestHandler.GetIdentityVerification)
		verificationRoute.GET("/:requestId", requestHandler.GetIdentityVerificationDetails)
		verificationRoute.POST("/accept", requestHandler.AcceptRequest)
		verificationRoute.POST("/reject", requestHandler.RejectRequest)
	}

	userManagementRoute := r.Group("/user-management")
	{
		adminStore := admin_repo.NewMysqlAdminRepository(s.DB)

		userManagementRoute.Use(middleware.CheckSession)
		userManagementRoute.Use(middleware.CheckMustChangePass(adminStore))

		userStore := user_repo.NewMysqlUserRepository(s.DB)
		vendorStore := vendor_repo.NewMysqlVendorRepository(s.DB)
		serviceStore := service_repo.NewMysqlServiceRepository(s.DB)
		notifStore := notification_repo.NewMysqlNotificationRepository(s.DB)

		managementService := user_management_service.NewService(userStore, notifStore, s.Encrypt, s.Hash)
		userService := user_service.NewService(userStore, s.Encrypt, s.Hash, s.JWT)
		vendorService := vendor_service.NewService(vendorStore, serviceStore, s.Encrypt, s.Hash)
		resourceService := resource_service.NewService(s.FS, s.Encrypt, s.Hash)

		managementHandler := userManagement.NewHandler(managementService, userService, vendorService, resourceService)

		userManagementRoute.GET("/users", managementHandler.GetUserList)
		userManagementRoute.GET("/users/:userId", managementHandler.ViewUserAccount)
		userManagementRoute.GET("/vendors", managementHandler.GetVendorList)
		userManagementRoute.GET("/vendors/:userId", managementHandler.ViewVendorAccount)
		userManagementRoute.POST("/ban/:userId", managementHandler.BanUser)
		userManagementRoute.POST("/unban/:userId", managementHandler.UnbanUser)
		userManagementRoute.POST("/restrict/:userId", managementHandler.RestrictUser)
		userManagementRoute.POST("/unrestrict/:userId", managementHandler.UnrestrictUser)
	}

	expertiseRoute := r.Group("/expertise")
	{
		adminStore := admin_repo.NewMysqlAdminRepository(s.DB)

		expertiseRoute.Use(middleware.CheckSession)
		expertiseRoute.Use(middleware.CheckMustChangePass(adminStore))

		expertiseStore := expertise_repo.NewMysqlExpertiseRepository(s.DB)
		userStore := user_repo.NewMysqlUserRepository(s.DB)
		vendorStore := vendor_repo.NewMysqlVendorRepository(s.DB)
		supportingImageStore := supportingimage_repo.NewMysqlImplementation(s.DB)

		expertiseService := expertise_service.NewService(userStore, vendorStore, expertiseStore, supportingImageStore, s.FS, s.Encrypt, s.Hash, s.JWT)
		handler := expertise.NewHandler(expertiseService)

		expertiseRoute.GET("", handler.GetAllExpertise)
		expertiseRoute.POST("", handler.CreateExpertise)
		expertiseRoute.POST("/tags", handler.AddTagToExpertise)
	}

	accountManagementRoute := r.Group("/account-management")
	{
		adminStore := admin_repo.NewMysqlAdminRepository(s.DB)

		accountManagementRoute.Use(middleware.CheckSession)
		accountManagementRoute.Use(middleware.CheckMustChangePass(adminStore))

		passwordResetStore := passwordreset_repo.NewMysqlPasswordResetRepository(s.DB)

		passwordResetService := passwordreset_service.NewService(adminStore, passwordResetStore, s.Encrypt, s.Hash)

		handler := accountmanagement.NewHandler(passwordResetService)

		accountManagementRoute.GET("/add", handler.AddAccount)
		accountManagementRoute.GET("/reset", handler.ResetRequests)
		accountManagementRoute.POST("/reset/fulfill", handler.FufillResetRequest)
		accountManagementRoute.POST("/reset/reject", handler.RejectResetRequest)
	}

	invitationRoute := r.Group("/invites")
	{
		inviteStore := invitation_repo.NewMysqlRepository(s.DB)
		adminStore := admin_repo.NewMysqlAdminRepository(s.DB)

		inviteService := invite_service.NewService(s.Domain, adminStore, inviteStore, s.Mailer, s.Encrypt, s.Hash)

		handler := invitation.NewHandler(inviteService)

		invitationRoute.POST("", handler.SendInvite, middleware.CheckSession)
		invitationRoute.GET("/join", handler.JoinInvite)
	}
}
