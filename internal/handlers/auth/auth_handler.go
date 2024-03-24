package auth

import (
	"nearbyassist/internal/controller/auth/v1"
	"nearbyassist/internal/controller/health"

	"github.com/labstack/echo/v4"
)

func AuthHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "auth route health check"

	r.POST("/", auth.AdminLogin).Name = "admin login"
	r.POST("/login", auth.HandleLogin).Name = "client login"
	r.POST("/logout", auth.HandlLogout).Name = "client logout"
}
