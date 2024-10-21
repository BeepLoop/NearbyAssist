package server

import (
	"nearbyassist/internal/handler/admin/auth"
	"nearbyassist/internal/handler/admin/complaint"
	"nearbyassist/internal/handler/admin/dashboard"
	"nearbyassist/internal/handler/admin/experiment"
	"nearbyassist/internal/handler/admin/management"
	map_handler "nearbyassist/internal/handler/admin/map"
	application "nearbyassist/internal/handler/admin/vendor_application"
	"nearbyassist/internal/handler/admin/verification"
	"nearbyassist/internal/middleware"
	admin_repo "nearbyassist/internal/repository/admin"
	application_repo "nearbyassist/internal/repository/application"
	complaint_repo "nearbyassist/internal/repository/complaint"
	dashboard_repo "nearbyassist/internal/repository/dashboard"
	map_repo "nearbyassist/internal/repository/map"
	verification_repo "nearbyassist/internal/repository/verification"
	admin_service "nearbyassist/internal/service/admin"
	application_service "nearbyassist/internal/service/application"
	complaint_service "nearbyassist/internal/service/complaint"
	dashboard_service "nearbyassist/internal/service/dashboard"
	management_service "nearbyassist/internal/service/management"
	map_service "nearbyassist/internal/service/map"
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
		dashboardStore := dashboard_repo.NewMysqlDashboardRepository(s.DB)
		dashboardService := dashboard_service.NewService(dashboardStore)
		dashboardHandler := dashboard.NewHandler(dashboardService)

		dashboardRoute.GET("", dashboardHandler.GetDashboard, middleware.CheckSession)
	}

	mapRoute := r.Group("/map")
	{
		mapStore := map_repo.NewMysqlMapRepository(s.DB)
		mapService := map_service.NewService(mapStore)
		mapHandler := map_handler.NewHandler(mapService)

		mapRoute.GET("", mapHandler.GetMap, middleware.CheckSession)
	}

	complaintRoute := r.Group("/complaints")
	{
		complaintStore := complaint_repo.NewMysqlComplaintRepository(s.DB)
		complaintService := complaint_service.NewService(complaintStore, s.FS, s.Encrypt)
		complaintHandler := complaint.NewHandler(complaintService)

		complaintRoute.GET("", complaintHandler.GetComplaints, middleware.CheckSession)
	}

	applicationRoute := r.Group("/vendor-applications")
	{
		applicationStore := application_repo.NewMysqlApplicationRepository(s.DB)
		applicationService := application_service.NewService(applicationStore, s.FS, s.Encrypt, s.JWT)
		applicationHandler := application.NewHandler(applicationService)

		applicationRoute.GET("", applicationHandler.GetVendorApplication, middleware.CheckSession)
	}

	verificationRoute := r.Group("/verification-requests")
	{
		requestStore := verification_repo.NewMysqlVerificationRepository(s.DB)
		requestService := verification_service.NewService(requestStore, s.FS, s.Encrypt, s.JWT)
		requestHandler := verification.NewHandler(requestService)

		verificationRoute.GET("", requestHandler.GetIdentityVerification, middleware.CheckSession)
	}

	managementRoute := r.Group("/account-management")
	{
		managementService := management_service.NewService()
		managementHandler := management.NewHandler(managementService)

		managementRoute.GET("", managementHandler.GetAccountManagement, middleware.CheckSession)
	}

	experimentRoute := r.Group("/test")
	{
		experimentHandler := experiment.NewHandler()

		experimentRoute.GET("", experimentHandler.GetExperiment, middleware.CheckSession)
	}
}
