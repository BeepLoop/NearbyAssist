package admin

import (
	"nearbyassist/internal/authenticator"
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
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

	accessToken, err := h.jwt.GenerateAdminAccessToken(authenticator.AdminOptions{
		Id:       admin.Id,
		Username: req.Username,
		Role:     admin.Role,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: authenticator.ACCESS_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	refreshToken, err := h.jwt.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: authenticator.REFRESH_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	session := models.NewSessionModel(refreshToken, h.idGen)
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
