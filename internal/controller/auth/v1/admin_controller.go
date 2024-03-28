package auth

import (
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminLogin(c echo.Context) error {
	admin := new(types.Admin)
	if err := c.Bind(admin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(admin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, admin)
}
