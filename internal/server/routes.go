package server

import (
	"nearbyassist/internal/middleware"
	service "nearbyassist/internal/service/admin"
	"nearbyassist/internal/service/user"
	store "nearbyassist/internal/store/admin"
)

func (s *Server) routes() {
	v1 := s.Echo.Group("/api/v1")
	{
		adminRoute := v1.Group("/admin")
		{
			adminStore := store.NewMysqlAdminStore(s.DB)
			h := service.NewAdminService(adminStore)

			adminRoute.GET("", h.BaseRoute)
			adminRoute.POST("/login", h.Login)
			adminRoute.POST("/refresh", h.Refresh, middleware.CheckAuth)
		}

		userRoute := v1.Group("/user")
		{
			h := user.NewHandler()
			userRoute.GET("", h.BaseRoute)
			userRoute.POST("/login", h.Login)
			userRoute.POST("/refresh", h.Refresh, middleware.CheckAuth)
		}
	}
}
