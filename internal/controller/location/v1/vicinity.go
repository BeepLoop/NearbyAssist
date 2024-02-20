package location

import (
	"nearbyassist/internal/db/query/location"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleVicinity(c echo.Context) error {
	latitude := c.QueryParam("lat")
	longitude := c.QueryParam("long")
	position := types.Position{Latitude: latitude, Longitude: longitude}

	locations, err := query.SearchVicinity(position)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, locations)
}
