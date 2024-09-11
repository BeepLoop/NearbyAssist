package utils

import (
	"errors"
	"strconv"
	"strings"
)

func ParseCoordinate(coordinate string) (float64, float64, error) {
	// Split the string into latitude and longitude
	coords := strings.Split(coordinate, ",")
	if len(coords) != 2 {
		return 0, 0, errors.New("Invalid coordinate")
	}

	// Parse the latitude and longitude
	lat, err := strconv.ParseFloat(coords[0], 64)
	if err != nil {
		return 0, 0, errors.New("Invalid latitude")
	}

	lng, err := strconv.ParseFloat(coords[1], 64)
	if err != nil {
		return 0, 0, errors.New("Invalid longitude")
	}

	return lat, lng, nil
}
