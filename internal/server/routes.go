package server

import (
	"nearbyassist/internal/service/admin"
	"nearbyassist/internal/service/user"
	store "nearbyassist/internal/store/admin"
)

func (s *Server) routes() {
	v1 := s.Echo.Group("/api/v1")
	{
		adminRoute := v1.Group("/admin")
		{
			adminStore := store.NewAdminStore(s.DB)
			h := admin.NewHandler(adminStore, s.Auth, s.IdGen)

			adminRoute.GET("", h.BaseRoute)
			adminRoute.POST("/login", h.Login)
			adminRoute.POST("/refresh", h.Refresh)
		}

		userRoute := v1.Group("/user")
		{
			h := user.NewHandler()
			userRoute.GET("", h.BaseRoute)
			userRoute.POST("/login", h.Login)
			userRoute.POST("/refresh", h.Refresh)
		}
	}
}
