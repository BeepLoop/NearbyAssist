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
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	admin := models.NewAdminModel(h.server.IdGen, h.server.DB)
	if admin == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := admin.HashUsername(req.Username, h.server.Hash.Hash); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: hash.HASH_ERROR,
			Error:   err.Error(),
		})
	}

	if _, err := admin.EncryptUsername(h.server.Encrypt.EncryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: encryption.ENCRYPTION_ERR,
			Error:   err.Error(),
		})
	}

	if _, err := admin.EncryptPassword(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: encryption.ENCRYPTION_ERR,
			Error:   err.Error(),
		})
	}

	if _, err := admin.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating account",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Account registered successfully",
		"staffId": admin.Id,
	})
}
