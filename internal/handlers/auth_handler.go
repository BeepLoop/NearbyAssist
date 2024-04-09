package handlers

import (
	"nearbyassist/internal/server"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	server *server.Server
}

func NewAuthHandler(server *server.Server) *AuthHandler {
	return &AuthHandler{server}
}

func (h *AuthHandler) HandleAdminLogin(c echo.Context) error {
	return c.JSON(200, "admin login")
}

func (h *AuthHandler) HandleLogin(c echo.Context) error {
	return c.JSON(200, "login")
}

func (h *AuthHandler) HandleLogout(c echo.Context) error {
	return c.JSON(200, "logout")
}
