package server

import (
	"nearbyassist/internal/middleware"
	"nearbyassist/internal/service"
	"nearbyassist/internal/store/admin"
	"nearbyassist/internal/store/user"
)

func (s *Server) routes() {
	v1 := s.Echo.Group("/api/v1")
	{
		adminRoute := v1.Group("/admin")
		{
			adminStore := admin.NewMysqlAdminStore(s.DB)
			h := service.NewAdminService(adminStore, s.Encrypt)

			adminRoute.GET("", h.BaseRoute)
			adminRoute.POST("/login", h.Login)
			adminRoute.POST("/refresh", h.Refresh, middleware.CheckAuth)
			adminRoute.POST("/logout", h.Logout, middleware.CheckAuth)
		}

		userRoute := v1.Group("/user")
		{
			userStore := user.NewMysqlUserStore(s.DB)
			h := service.NewUserService(userStore, s.Encrypt)

			userRoute.GET("", h.BaseRoute)
			userRoute.POST("/login", h.Login)
			userRoute.POST("/refresh", h.Refresh, middleware.CheckAuth)
			userRoute.POST("/logout", h.Logout, middleware.CheckAuth)
		}

		resourceRoute := v1.Group("/resource")
		{
			resourceRoute.Use(middleware.CheckAuth)
			resourceRoute.Use(middleware.CheckRole)

			h := service.NewResourceService(s.Encrypt)

			resourceRoute.GET("/:path", h.GetFile)
		}
	}
}
