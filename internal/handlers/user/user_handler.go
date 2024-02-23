package user

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/user/v1"

	"github.com/labstack/echo/v4"
)

func UserHandler(r *echo.Group) {

	r.GET("/health", health.HealthCheck)
	r.GET(":userId", user.GetUser)
}
