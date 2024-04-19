package handlers

import (
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type adminHandler struct {
	server *server.Server
}

func NewAdminHandler(server *server.Server) *adminHandler {
	return &adminHandler{
		server: server,
	}
}

func (h *adminHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Admin base route",
	})
}

func (h *adminHandler) HandleRegisterStaff(c echo.Context) error {
	// TODO: Implement registering staff accounts

	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Register staff",
	})
}
