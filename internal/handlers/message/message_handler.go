package message

import (
	"nearbyassist/internal/controller/health"
	"nearbyassist/internal/controller/message/v1"

	"github.com/labstack/echo/v4"
)

func MessageHandler(r *echo.Group) {
	r.GET("/health", health.HealthCheck).Name = "message route health check"

	r.GET("/messages", message.GetMessages).Name = "get messages between sender and receiver"
	r.GET("/ws", message.HandleChat).Name = "websocket route for chat"
	r.GET("/conversations", message.GetConversations).Name = "get all users you chatted with"
}
