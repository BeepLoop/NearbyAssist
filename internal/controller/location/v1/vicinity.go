package location

import (
	"nearbyassist/internal/db/query/location"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Point struct {
	OwnerId   int
	Address   string
	Latitude  float64
	Longitude float64
}

func HandleVicinity(c echo.Context) error {
	latitude := c.QueryParam("lat")
	longitude := c.QueryParam("long")
	position := types.Position{Latitude: latitude, Longitude: longitude}

	radius := c.QueryParam("radius")
	if radius == "" {
		radius = "200"
	}

	result, err := query.SearchVicinity(position, radius)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	var locations []Point

	for _, location := range result {
		lat, long := utils.LatlongExtractor(location.Point)
		point := Point{
			OwnerId:   location.OwnerId,
			Address:   location.Address,
			Latitude:  lat,
			Longitude: long,
		}
		locations = append(locations, point)
	}

	return c.JSON(http.StatusOK, locations)
}
