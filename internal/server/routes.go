package server

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/handlers"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	// middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	// Custom validator
	e.Validator = &utils.Validator{Validator: validator.New()}

	// File server
	e.Static("/resource", "store/").Name = "file resource"

	// Routes
	e.GET("/", func(c echo.Context) error {
		routes := e.Routes()

		return c.JSON(http.StatusOK, routes)
	})

	e.GET("/health", health.HealthCheck).Name = "health check"
	handlers.RouteHandlerV1(e.Group("/v1"))

	return e
}
