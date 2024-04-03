package message

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/message/v1"

	"github.com/labstack/echo/v4"
)

func MessageHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "message route health check"

	r.GET("/conversations", message.GetMessages).Name = "get messages between sender and receiver"
	r.GET("/chat", message.HandleChat).Name = "websocket route for chat"
	r.GET("/acquaintances", message.GetAcquaintances).Name = "get all users you chatted with"
}
