package handlers

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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
	req := models.NewAdminModel()
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
	}
	req.Password = string(hash)

	staffId, err := h.server.DB.NewStaff(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Register staff",
		"staffId": staffId,
	})
}
