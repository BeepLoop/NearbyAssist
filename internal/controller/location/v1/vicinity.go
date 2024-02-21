package location

import (
	"nearbyassist/internal/db/query/location"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Point struct {
	Address   string
	Latitude  string
	Longitude string
}

func HandleVicinity(c echo.Context) error {
	latitude := c.QueryParam("lat")
	longitude := c.QueryParam("long")
	position := types.Position{Latitude: latitude, Longitude: longitude}

	result, err := query.SearchVicinity(position)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	var locations []Point

	for _, location := range result {
		latlong := utils.LatlongExtractor(location.Point)
		point := Point{Address: location.Address, Latitude: latlong[0], Longitude: latlong[1]}
		locations = append(locations, point)
	}

	return c.JSON(http.StatusOK, locations)
}
