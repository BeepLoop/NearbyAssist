package server

import (
	"nearbyassist/internal/middleware"
	"nearbyassist/internal/service"
	"nearbyassist/internal/store/admin"
	"nearbyassist/internal/store/service"
	"nearbyassist/internal/store/tag"
	"nearbyassist/internal/store/user"
	"nearbyassist/internal/store/vendor"
)

func (s *Server) routes() {
	v1 := s.Echo.Group("/api/v1")
	{
		// ===== HEALTH =======
		healthRoute := v1.Group("/health")
		{
			h := handler.NewHealthService()
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
			h := handler.NewAdminService(adminStore, s.Encrypt, s.JWT)

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
			h := handler.NewUserService(userStore, s.Encrypt, s.JWT)

			userRoute.POST("/login", h.Login)
			userRoute.POST("/refresh", h.Refresh)

			protected := userRoute.Group("/protected")
			{
				protected.Use(middleware.CheckAuth(s.JWT))

				protected.POST("/logout", h.Logout)
				protected.GET("/me", h.BaseRoute)
				protected.GET("/verified", h.Verified)
			}
		}

		// ===== RESOURCE =======
		resourceRoute := v1.Group("/resource")
		{
			resourceRoute.Use(middleware.CheckAuth(s.JWT))
			resourceRoute.Use(middleware.CheckRole(s.JWT))

			h := handler.NewResourceService(s.Encrypt)

			resourceRoute.GET("/:path", h.GetFile)
		}

		// ===== TAGS =======
		tagRoute := v1.Group("/tags")
		{
			tagStore := tag.NewMysqlTagStore(s.DB)
			h := handler.NewTagService(tagStore)

			tagRoute.GET("", h.BaseRoute)
		}

		// ===== VENDOR =======
		vendorRoute := v1.Group("/vendors")
		{
			vendorRoute.Use(middleware.CheckAuth(s.JWT))

			vendorStore := vendor.NewMysqlVendorStore(s.DB)
			h := handler.NewVendorService(vendorStore)

			vendorRoute.GET("/:vendorId", h.GetVendor)
		}

		// ===== SERVICES =======
		serviceRoute := v1.Group("/services")
		{
			serviceRoute.Use(middleware.CheckAuth(s.JWT))

			serviceStore := service.NewMysqlServiceStore(s.DB)
			h := handler.NewServiceService(serviceStore, s.Encrypt, s.JWT, s.RouteEngine, s.SuggestionEngine)

			serviceRoute.GET("", h.GetServices)
			serviceRoute.POST("", h.Create)
			serviceRoute.GET("/search", h.Search)
			serviceRoute.GET("/:serviceId", h.GetService)
			serviceRoute.PUT("/:serviceId", h.Update)
			serviceRoute.DELETE("/:serviceId", h.Delete)
			serviceRoute.GET("/vendor/:vendorId", h.GetVendorServices)
			serviceRoute.GET("/route/:serviceId", h.FindRoute)
		}
	}
}
