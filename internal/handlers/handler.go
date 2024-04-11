package handlers

import (
	"nearbyassist/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	server *server.Server
}

func NewHandler(server *server.Server) *handler {
	return &handler{
		server: server,
	}
}

func (h *handler) HandleBaseRoute(c echo.Context) error {
	routes := h.server.Echo.Routes()
	return c.JSON(http.StatusOK, routes)
}
