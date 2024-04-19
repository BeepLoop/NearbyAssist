package server

import (
	"net/http"

	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) registerMiddleware() {
	s.Echo.Pre(middleware.RemoveTrailingSlash())
	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	s.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"127.0.0.1:5173", "localhost:5173"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	s.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
}
