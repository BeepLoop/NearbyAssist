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

	dashboardStore := dashboard_repo.NewMysqlDashboardRepository(s.DB)
	dashboardService := dashboard_service.NewService(dashboardStore)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	r.GET("/dashboard", dashboardHandler.GetDashboard, middleware.CheckSession)

	mapStore := map_repo.NewMysqlMapRepository(s.DB)
	mapService := map_service.NewService(mapStore)
	mapHandler := map_handler.NewHandler(mapService)

	r.GET("/map", mapHandler.GetMap, middleware.CheckSession)

	complaintStore := complaint_repo.NewMysqlComplaintRepository(s.DB)
	complaintService := complaint_service.NewService(complaintStore, s.FS, s.Encrypt)
	complaintHandler := complaint.NewHandler(complaintService)

	r.GET("/complaints", complaintHandler.GetComplaints, middleware.CheckSession)

	applicationStore := application_repo.NewMysqlApplicationRepository(s.DB)
	applicationService := application_service.NewService(applicationStore, s.FS, s.Encrypt, s.JWT)
	applicationHandler := application.NewHandler(applicationService)

	r.GET("/vendor-applications", applicationHandler.GetVendorApplication, middleware.CheckSession)

	requestStore := verification_repo.NewMysqlVerificationRepository(s.DB)
	requestService := verification_service.NewService(requestStore, s.FS, s.Encrypt, s.JWT)
	requestHandler := verification.NewHandler(requestService)

	r.GET("/verification-requests", requestHandler.GetIdentityVerification, middleware.CheckSession)

	managementService := management_service.NewService()
	managementHandler := management.NewHandler(managementService)

	r.GET("/account-management", managementHandler.GetAccountManagement, middleware.CheckSession)

	experimentHandler := experiment.NewHandler()

	r.GET("/test", experimentHandler.GetExperiment, middleware.CheckSession)
}
