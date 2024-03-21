package message

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/message/v1"

	"github.com/labstack/echo/v4"
)

func MessageHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck)

	r.GET("/conversations", message.GetMessages)
	r.GET("/chat", message.HandleChat)
	r.GET("/acquaintances", message.GetAcquaintances)
}
