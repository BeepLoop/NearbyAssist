package handlers

import (
	"nearbyassist/internal/encryption"
	"nearbyassist/internal/hash"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
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
	req := new(request.NewAdminRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	}

	admin := models.NewAdminModel(h.server.IdGen, h.server.DB)
	if admin == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := admin.HashUsername(req.Username, h.server.Hash.Hash); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	}

	if _, err := admin.EncryptUsername(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.ENCRYPTION_ERR)
	}

	if _, err := admin.EncryptPassword(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, encryption.ENCRYPTION_ERR)
	}

	if _, err := admin.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while creating admin account")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Account registered successfully",
		"staffId": admin.Id,
	})
}
