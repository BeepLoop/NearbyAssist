package user

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
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

	return c.JSON(http.StatusOK, req)
}
