package service

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *AdminService) Login(c echo.Context) error {
	req := new(request.AdminLoginPayload)
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

	usernameHash, err := auth.Sha256([]byte(req.Username))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error hashing username",
			Error:   err.Error(),
		})
	}

	admin, err := h.store.FindByUsernameHash(usernameHash)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
			Message: "Invalid credentials",
			Error:   err.Error(),
		})
	}

	if auth.IsPasswordMatch(admin.Password, req.Password) == false {
		return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
			Message: "Invalid credentials",
			Error:   "Invalid credentials",
		})
	}

	accessToken, err := auth.GenerateAdminAccessToken(auth.AdminOptions{
		Id:       admin.Id,
		Username: req.Username,
		Role:     admin.Role,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.ACCESS_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.REFRESH_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	sessionId, err := auth.GenerateNanoId()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.NANO_ID_ERR,
			Error:   err.Error(),
		})
	}
	session := models.NewSessionModel(sessionId, refreshToken)
	if session == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if err := h.store.Login(session); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating session",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"role":         admin.Role,
		"adminId":      admin.Id,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
