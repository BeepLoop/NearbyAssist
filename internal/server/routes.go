package server

import (
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	e.Validator = &utils.Validator{Validator: validator.New()}

	e.GET("/health", s.HealthHandler)

	s.HandleVersionOneRoutes(e.Group("/v1"))

	return e
}
