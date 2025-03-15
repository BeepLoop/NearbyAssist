package server

import (
	"nearbyassist/internal/handler/admin/auth"
	"nearbyassist/internal/handler/admin/complaint"
	"nearbyassist/internal/handler/admin/dashboard"
	"nearbyassist/internal/handler/admin/expertise"
	"nearbyassist/internal/handler/admin/management"
	map_handler "nearbyassist/internal/handler/admin/map"
	application "nearbyassist/internal/handler/admin/vendor_application"
	"nearbyassist/internal/handler/admin/verification"
	"nearbyassist/internal/middleware"
	admin_repo "nearbyassist/internal/repository/admin"
	application_repo "nearbyassist/internal/repository/application"
	bug_report_repo "nearbyassist/internal/repository/bug_report"
	dashboard_repo "nearbyassist/internal/repository/dashboard"
	expertise_repo "nearbyassist/internal/repository/expertise"
	map_repo "nearbyassist/internal/repository/map"
	notification_repo "nearbyassist/internal/repository/notification"
	report_user_repo "nearbyassist/internal/repository/report_user"
	tag_repo "nearbyassist/internal/repository/tag"
	user_repo "nearbyassist/internal/repository/user"
	verification_repo "nearbyassist/internal/repository/verification"
	admin_service "nearbyassist/internal/service/admin"
	application_service "nearbyassist/internal/service/application"
	complaint_service "nearbyassist/internal/service/complaint"
	dashboard_service "nearbyassist/internal/service/dashboard"
	expertise_service "nearbyassist/internal/service/expertise"
	management_service "nearbyassist/internal/service/management"
	map_service "nearbyassist/internal/service/map"
	resource_service "nearbyassist/internal/service/resource"
	tag_service "nearbyassist/internal/service/tag"
	verification_service "nearbyassist/internal/service/verification"

	"github.com/labstack/echo/v4"
)

func (s *Server) AdminRoutes(r *echo.Group) {
	authStore := admin_repo.NewMysqlAdminRepository(s.DB)
	authService := admin_service.NewService(authStore, s.Encrypt, s.Hash)
	authHandler := auth.NewHandler(authService)

	r.GET("/login", authHandler.GetLogin)
	r.POST("/login", authHandler.PostLogin)
	r.POST("/logout", authHandler.PostLogout)

	dashboardRoute := r.Group("/dashboard")
	{
		dashboardRoute.Use(middleware.CheckSession)

		dashboardStore := dashboard_repo.NewMysqlDashboardRepository(s.DB)
		dashboardService := dashboard_service.NewService(dashboardStore)
		dashboardHandler := dashboard.NewHandler(dashboardService)

		dashboardRoute.GET("", dashboardHandler.GetDashboard)
	}

	mapRoute := r.Group("/map")
	{
		mapRoute.Use(middleware.CheckSession)

		tagStore := tag_repo.NewMysqlTagRepository(s.DB)
		tagService := tag_service.NewService(tagStore)

		mapStore := map_repo.NewMysqlMapRepository(s.DB)
		mapService := map_service.NewService(mapStore)
		mapHandler := map_handler.NewHandler(mapService, tagService)

		mapRoute.GET("", mapHandler.GetMap)
	}

	complaintRoute := r.Group("/complaints")
	{
		complaintRoute.Use(middleware.CheckSession)

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
		applicationRoute.Use(middleware.CheckSession)

		notificationStore := notification_repo.NewMysqlNotificationRepository(s.DB)
		applicationStore := application_repo.NewMysqlApplicationRepository(s.DB)

		applicationService := application_service.NewService(applicationStore, notificationStore, s.FS, s.Encrypt, s.JWT)
		resourceService := resource_service.NewService(s.FS, s.Encrypt, s.Hash)

		applicationHandler := application.NewHandler(applicationService, resourceService)

		applicationRoute.GET("", applicationHandler.GetVendorApplication)
		applicationRoute.GET("/:applicationId", applicationHandler.GetVendorApplicationDetails)
		applicationRoute.POST("/accept/:applicationId", applicationHandler.AcceptRequest)
		applicationRoute.POST("/reject/:applicationId", applicationHandler.RejectRequest)
	}

	verificationRoute := r.Group("/verification-requests")
	{
		verificationRoute.Use(middleware.CheckSession)

		requestStore := verification_repo.NewMysqlVerificationRepository(s.DB)
		notificationStore := notification_repo.NewMysqlNotificationRepository(s.DB)

		requestService := verification_service.NewService(requestStore, notificationStore, s.FS, s.Encrypt, s.JWT)
		resourceService := resource_service.NewService(s.FS, s.Encrypt, s.Hash)

		requestHandler := verification.NewHandler(requestService, resourceService)

		verificationRoute.GET("", requestHandler.GetIdentityVerification)
		verificationRoute.GET("/:requestId", requestHandler.GetIdentityVerificationDetails)
		verificationRoute.POST("/accept/:requestId", requestHandler.AcceptRequest)
		verificationRoute.POST("/reject/:requestId", requestHandler.RejectRequest)
	}

	managementRoute := r.Group("/account-management")
	{
		managementRoute.Use(middleware.CheckSession)

		userStore := user_repo.NewMysqlUserRepository(s.DB)
		notifStore := notification_repo.NewMysqlNotificationRepository(s.DB)

		managementService := management_service.NewService(userStore, notifStore, s.Encrypt, s.Hash)
		managementHandler := management.NewHandler(managementService)

		managementRoute.GET("", managementHandler.GetAccountManagement)
		managementRoute.GET("/:userId", managementHandler.ViewUserAccount)
		managementRoute.POST("/ban/:userId", managementHandler.BanUser)
		managementRoute.POST("/unban/:userId", managementHandler.UnbanUser)
		managementRoute.POST("/restrict/:userId", managementHandler.RestrictUser)
		managementRoute.POST("/unrestrict/:userId", managementHandler.UnrestrictUser)
	}

	expertiseRoute := r.Group("/expertise")
	{
		expertiseRoute.Use(middleware.CheckSession)

		expertiseStore := expertise_repo.NewMysqlExpertiseRepository(s.DB)
		expertiseService := expertise_service.NewService(expertiseStore, s.Encrypt, s.Hash)
		handler := expertise.NewHandler(expertiseService)

		expertiseRoute.GET("", handler.GetAllExpertise)
		expertiseRoute.POST("", handler.CreateExpertise)
		expertiseRoute.POST("/tags", handler.AddTagToExpertise)
	}
}
