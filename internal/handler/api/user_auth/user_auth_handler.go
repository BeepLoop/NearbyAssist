package userauth_handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	userauth_service "nearbyassist/internal/service/user_auth"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type userAuthHandler struct {
	authService *userauth_service.Service
}

func NewHandler(authService *userauth_service.Service) *userAuthHandler {
	return &userAuthHandler{
		authService: authService,
	}
}

func (h *userAuthHandler) ThirdPartyLogin(c echo.Context) error {
	req := new(request.UserLoginPayload)
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

	response, err := h.authService.ThirdPartyLogin(req)
	if err != nil {
		if strings.Contains(err.Error(), userauth_service.ERR_BANNED_USER) {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "User is banned",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error logging in",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *userAuthHandler) Refresh(c echo.Context) error {
	req := new(request.TokenRefreshPayload)
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

	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	accessToken, err := h.authService.Refresh(bearerToken, req.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), userauth_service.ERR_BANNED_USER) {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "User is banned",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error refresh token",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"accessToken": accessToken,
	})
}

func (h *userAuthHandler) Logout(c echo.Context) error {
	req := new(request.LogoutPayload)
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

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, models.Error{
			Message: "Error logging out",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
