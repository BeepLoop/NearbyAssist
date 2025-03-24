package server

import (
	"nearbyassist/internal/handler/web"
)

func (s *Server) routes() {
	s.Echo.Static("/static", "static")
	s.Echo.Static("/public", "public")

	s.Echo.GET("", web.GetIndex)
	s.Echo.GET("/privacy_policy", web.GetPrivacyPolicy)
	s.Echo.GET("/terms_and_conditions", web.GetTermsAndConditions)
	s.Echo.GET("/account_deletion", web.GetAccountDeletionInstructions)

	api := s.Echo.Group("/api")
	{
		v1 := api.Group("/v1")
		s.v1ApiRoutes(v1)
	}

	auth := s.Echo.Group("/auth")
	s.AuthRoutes(auth)

	admin := s.Echo.Group("/admin")
	s.AdminRoutes(admin)

	s.Echo.RouteNotFound("/*", web.GetNotFound)
}
