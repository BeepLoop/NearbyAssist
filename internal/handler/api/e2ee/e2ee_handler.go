package e2ee

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	e2ee_service "nearbyassist/internal/service/e2ee"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type e2eeHandler struct {
	e2eeService *e2ee_service.Service
}

func NewHandler(e2eeService *e2ee_service.Service) *e2eeHandler {
	return &e2eeHandler{e2eeService: e2eeService}
}

func (h *e2eeHandler) SaveKeys(c echo.Context) error {
	req := new(request.PEM)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
	}

	bearerToken := utils.BearerTokenFromHeader(c)

	if err := h.e2eeService.SaveKeys(bearerToken, req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error saving keys",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *e2eeHandler) GetKeys(c echo.Context) error {
	bearerToken := utils.BearerTokenFromHeader(c)

	keys, err := h.e2eeService.GetKeys(bearerToken)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return echo.NewHTTPError(http.StatusNotFound, models.Error{
				Message: "Keys not found",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving keys",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"publicKey":  keys["publicKey"],
		"privateKey": keys["privateKey"],
	})
}

func (h *e2eeHandler) GetPublicKey(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "userId must be an integer",
			Error:   "userId must be an integer",
		})
	}

	publicKey, err := h.e2eeService.GetPublicKey(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Public key not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"key": publicKey,
	})
}
