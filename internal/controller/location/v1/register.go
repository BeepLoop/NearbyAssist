package location

import (
	"nearbyassist/internal/db"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterLocation(c echo.Context) error {
	var location types.LocationRegister
	err := c.Bind(&location)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	err = c.Validate(location)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	err = db.RegisterLocation(location)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}
