package user

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	user_service "nearbyassist/internal/service/user"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userService *user_service.Service
}

func NewHandler(userService *user_service.Service) *userHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) Login(c echo.Context) error {
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

	loginResp, err := h.userService.Login(req)
	if err != nil {
		if strings.Contains(err.Error(), user_service.ERR_BANNED_USER) {
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

	return c.JSON(http.StatusCreated, loginResp)
}

func (h *userHandler) Refresh(c echo.Context) error {
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

	accessToken, err := h.userService.Refresh(bearerToken, req.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), user_service.ERR_BANNED_USER) {
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

func (h *userHandler) Logout(c echo.Context) error {
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

	if err := h.userService.Logout(req.RefreshToken); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, models.Error{
			Message: "Error logging out",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *userHandler) GetUser(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	user, err := h.userService.GetUser(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, models.Error{
			Message: "Error while retrieving user data",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"user": user,
	})
}

func (h *userHandler) GetUserVerification(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	isVerified, err := h.userService.IsVerified(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Error checking if user is verified",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"verified": isVerified,
	})
}

func (h *userHandler) AddSocial(c echo.Context) error {
	req := new(request.AddSocialPayload)
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

	if err := h.userService.AddSocial(bearerToken, req.Url); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error adding social",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, nil)
}

func (h *userHandler) DeleteSocial(c echo.Context) error {
	req := new(request.DeleteSocialPayload)
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

	if err := h.userService.DeleteSocial(bearerToken, req.Url); err != nil {
		if strings.Contains(err.Error(), "social not found") {
			return echo.NewHTTPError(http.StatusNotFound, models.Error{
				Message: "social does not exists",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error deleting social",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
