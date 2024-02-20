package location

import (
	"nearbyassist/internal/db/query/location"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleLocations(c echo.Context) error {
	locations, err := query.GetLocations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, locations)
}
