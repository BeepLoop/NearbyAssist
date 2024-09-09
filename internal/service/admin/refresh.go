package service

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *AdminService) Refresh(c echo.Context) error {
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

	// Check if refreshToken exists
	if err := h.store.DoesRefreshTokenExists(req.RefreshToken); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Session not found",
			Error:   err.Error(),
		})
	}

	// Check if refreshToken is blacklisted
	if err := h.store.IsRefreshTokenBlacklisted(req.RefreshToken); err == nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Token blacklisted",
			Error:   "Token blacklisted",
		})
	}

	// Generate new accessToken
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}
	adminId, ok := claims["adminId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Admin Id not found in claims",
			Error:   "Admin Id not found in claims",
		})
	}

	admin, err := h.store.FindById(adminId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Admin not found",
			Error:   "Admin not found",
		})
	}

	accessToken, err := auth.GenerateAdminAccessToken(auth.AdminOptions{
		Id:       adminId,
		Username: admin.Username,
		Role:     admin.Role,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.ACCESS_TOKEN_ERR,
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"accessToken": accessToken,
	})
}
