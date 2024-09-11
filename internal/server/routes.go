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
		// ===== HEALTH =======
		healthRoute := v1.Group("/health")
		{
			h := service.NewHealthService()
			healthRoute.GET("", h.BaseRoute)

			protected := healthRoute.Group("/protected")
			{
				protected.Use(middleware.CheckAuth(s.JWT))

				protected.GET("", h.BaseRoute)
			}
		}

		// ===== ADMIN =======
		adminRoute := v1.Group("/admin")
		{
			adminStore := admin.NewMysqlAdminStore(s.DB)
			h := service.NewAdminService(adminStore, s.Encrypt, s.JWT)

			adminRoute.GET("", h.BaseRoute)
			adminRoute.POST("/login", h.Login)
			adminRoute.POST("/refresh", h.Refresh)

			protected := adminRoute.Group("/protected")
			{
				protected.Use(middleware.CheckAuth(s.JWT))

				protected.POST("/logout", h.Logout)
			}
		}

		// ===== USER =======
		userRoute := v1.Group("/user")
		{
			userStore := user.NewMysqlUserStore(s.DB)
			h := service.NewUserService(userStore, s.Encrypt, s.JWT)

			userRoute.POST("/login", h.Login)
			userRoute.POST("/refresh", h.Refresh)

			protected := userRoute.Group("/protected")
			{
				protected.Use(middleware.CheckAuth(s.JWT))

				protected.POST("/logout", h.Logout)
				protected.GET("/me", h.BaseRoute)
			}
		}

		// ===== RESOURCE =======
		resourceRoute := v1.Group("/resource")
		{
			resourceRoute.Use(middleware.CheckAuth(s.JWT))
			resourceRoute.Use(middleware.CheckRole(s.JWT))

			h := service.NewResourceService(s.Encrypt)

			resourceRoute.GET("/:path", h.GetFile)
		}
	}
}
