package auth

import (
	"nearbyassist/internal/controller/health"

	"github.com/labstack/echo/v4"
)

func AuthHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck)

	r.GET("/login", HandleLogin)
	r.POST("/register", HandleRegister)
}
