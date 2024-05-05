package utils

import (
	"errors"
	"nearbyassist/internal/types"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetSearchParams(c echo.Context) (*types.SearchParams, error) {
	latitude := c.QueryParam("lat")
	longitude := c.QueryParam("long")
	radius := c.QueryParam("radius")
	query := c.QueryParam("q")

	if latitude == "" || longitude == "" || query == "" {
		return nil, errors.New("missing params")
	}

	if radius == "" {
		radius = "500"
	}

	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return nil, err
	}

	long, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return nil, err
	}

	rad, err := strconv.ParseFloat(radius, 64)
	if err != nil {
		return nil, err
	}

	params := types.SearchParams{
		Latitude:  lat,
		Longitude: long,
		Radius:    rad,
		Query:     query,
	}

	return &params, nil
}
