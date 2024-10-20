package server

import (
	"nearbyassist/internal/handler/web"
)

func (s *Server) routes() {
	s.Echo.Static("/static", "static")
	s.Echo.Static("/public", "public")

	s.Echo.GET("", web.GetIndex)

	api := s.Echo.Group("/api")
	{
		v1 := api.Group("/v1")
		s.v1ApiRoutes(v1)
	}

	admin := s.Echo.Group("/admin")
	s.AdminRoutes(admin)

	s.Echo.RouteNotFound("/*", web.GetNotFound)
}
