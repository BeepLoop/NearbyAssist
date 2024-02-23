package location

import (
	"nearbyassist/internal/db/query/location"
	"nearbyassist/internal/types"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleVicinity(c echo.Context) error {
	params, err := utils.GetSearchParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	result, err := query.SearchVicinity(params)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	var locations []types.TransformedServiceData

	for _, location := range result {
		lat, long := utils.LatlongExtractor(location.Location)
		point := types.TransformedServiceData{
			Vendor:      location.Vendor,
			Title:       location.Title,
			Description: location.Description,
			Rate:        location.Rate,
			Latitude:    lat,
			Longitude:   long,
		}
		locations = append(locations, point)
	}

	return c.JSON(http.StatusOK, locations)
}
