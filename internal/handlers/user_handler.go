package handlers

import (
	"nearbyassist/internal/encryption"
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	server *server.Server
}

func NewUserHandler(server *server.Server) *userHandler {
	return &userHandler{
		server: server,
	}
}

func (h *userHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "User base route",
	})
}

func (h *userHandler) HandleCount(c echo.Context) error {
	param := c.QueryParam("filter")
	var filter models.UserStatusFilter
	switch param {
	case "":
		filter = models.USER_STATUS_ALL
	case "all":
		filter = models.USER_STATUS_ALL
	case "verified":
		filter = models.USER_STATUS_VERIFIED
	case "unverified":
		filter = models.USER_STATUS_UNVERIFIED
	default:
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid filter",
			Error:   "Invalid filter",
		})
	}

	user := models.NewUserModel(h.server.IdGen, h.server.DB)
	if user == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	count, err := user.Count(filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting user count",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"count": count,
	})
}

func (h *userHandler) HandleCheckVerification(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "User ID not found in JWT",
			Error:   "User ID not found in JWT",
		})
	}

	user := models.NewUserModel(h.server.IdGen, h.server.DB)
	if user == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}
	user.Id = id
	isVerified := user.IsVerified()

	return c.JSON(http.StatusOK, utils.Mapper{
		"userId":   id,
		"verified": isVerified,
	})
}

func (h *userHandler) HandleGetMyDetails(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "User ID not found in JWT",
			Error:   "User ID not found in JWT",
		})
	}

	user := models.NewUserModel(h.server.IdGen, h.server.DB)
	if user == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := user.FindById(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "User not found",
			Error:   err.Error(),
		})
	}

	if _, err := user.DecryptName(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error decrypting name",
			Error:   encryption.DECRYPTION_ERR,
		})
	}
	if _, err := user.DecryptEmail(h.server.Encrypt.DecryptString); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error decrypting email",
			Error:   encryption.DECRYPTION_ERR,
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"user": user,
	})
}
