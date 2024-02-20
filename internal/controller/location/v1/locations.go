package location

import (
	"nearbyassist/internal/db"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleLocations(c echo.Context) error {
	locations, err := db.GetLocations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, locations)
}
