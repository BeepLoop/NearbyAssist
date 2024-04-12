package handlers

import (
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
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
	return c.JSON(http.StatusOK, utils.Mapper{"routes": routes})
}

func (h *handler) HandleV1BaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "V1 base route",
	})
}
