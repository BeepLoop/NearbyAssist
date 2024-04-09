package handlers

import (
	"nearbyassist/internal/server"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type healthHandler struct {
	server *server.Server
}

func NewHealthHandler(server *server.Server) *healthHandler {
	return &healthHandler{
		server: server,
	}
}

func (h *healthHandler) HandleHealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"health": "ok",
		"time":   time.Now(),
	})
}
