package server

import "github.com/labstack/echo/v4/middleware"

func (s *Server) registerMiddleware() {
	s.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.RemoveTrailingSlash())
	s.Echo.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
}
