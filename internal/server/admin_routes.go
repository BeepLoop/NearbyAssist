package server

import (
	"nearbyassist/internal/handler/admin"
	"nearbyassist/internal/middleware"

	"github.com/labstack/echo/v4"
)

func (s *Server) AdminRoutes(r *echo.Group) {
	adminHandler := admin.NewHandler(s.DB, s.Hash, s.Encrypt)

	r.GET("/login", adminHandler.GetLogin)
	r.POST("/login", adminHandler.PostLogin)

	r.POST("/logout", adminHandler.PostLogout)

	r.GET("/dashboard", adminHandler.GetDashboard, middleware.CheckSession)

	r.GET("/map", adminHandler.GetMap)

	r.GET("/complaints", adminHandler.GetComplaints)

	r.GET("/vendor-applications", adminHandler.GetVendorApplication)

	r.GET("/verification-requests", adminHandler.GetIdentityVerification)

	r.GET("/account-management", adminHandler.GetAccountManagement)

	r.GET("/test", adminHandler.GetExperiment)
}
