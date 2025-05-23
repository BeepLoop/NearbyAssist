package server

import (
	"nearbyassist/internal/handler/web"
	"net/http"

	"github.com/labstack/echo/v4"
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

	dev := s.Echo.Group("/dev")
	s.devRoutes(dev)

	auth := s.Echo.Group("/auth")
	s.authRoutes(auth)

	admin := s.Echo.Group("/admin")
	s.adminRoutes(admin)

	ws := s.Echo.Group("/ws")
	s.websocketRoute(ws)

	s.Echo.RouteNotFound("/*", func(c echo.Context) error {
		return c.String(http.StatusNotFound, "not found")
	})
}
