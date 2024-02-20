package auth

import (
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminLogin(c echo.Context) error {
	admin := new(types.Admin)
	if err := c.Bind(admin); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if err := c.Validate(admin); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	return c.JSON(http.StatusOK, admin)
}
