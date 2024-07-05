package handlers

import (
	"nearbyassist/internal/hash"
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

	if hashed, err := h.server.Hash.Hash([]byte(req.Username)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	} else {
		req.UsernameHash = hashed
	}

	if cipher, err := h.server.Encrypt.EncryptString(req.Username); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	} else {
		req.Username = cipher
	}

	if hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, hash.HASH_ERROR)
	} else {
		req.Password = string(hashed)
	}

	staffId, err := h.server.DB.NewStaff(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Register staff",
		"staffId": staffId,
	})
}
