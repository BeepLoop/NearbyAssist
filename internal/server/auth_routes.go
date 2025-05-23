package server

import (
	"nearbyassist/internal/handler/admin/auth"
	"nearbyassist/internal/handler/web"
	admin_repo "nearbyassist/internal/repository/admin"
	adminauth_service "nearbyassist/internal/service/admin_auth"

	"github.com/labstack/echo/v4"
)

func (s *Server) authRoutes(r *echo.Group) {
	adminStore := admin_repo.NewMysqlAdminRepository(s.DB)
	authService := adminauth_service.NewService(adminStore, s.Encrypt, s.Hash)

	authHandler := auth.NewHandler(authService)

	r.GET("/login", authHandler.GetLogin)
	r.POST("/login", authHandler.PostLogin)
	r.POST("/logout", authHandler.PostLogout)

	r.RouteNotFound("/*", web.GetNotFound)
}
