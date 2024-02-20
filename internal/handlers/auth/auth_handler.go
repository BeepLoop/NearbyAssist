package auth

import (
	"nearbyassist/internal/controller/auth/v1"
	"nearbyassist/internal/controller/health"

	"github.com/labstack/echo/v4"
)

func AuthHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck)

	r.GET("/login", auth.HandleLogin)
	r.POST("/register", auth.HandleRegister)
}
