package models

import (
	"errors"
	"strconv"
	"strings"
)

type GeoSpatialModel struct {
	Latitude  float64 `json:"latitude" db:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" db:"longitude" validate:"required"`
}

// String coordinate {latitude},{longitude} ex: 7.544645340539252,126.14141292293003
func (l *GeoSpatialModel) FromString(coordinate string) error {
	latitude, longitude, err := ParseCoordinate(coordinate)
	if err != nil {
		return err
	}

	l.Latitude = latitude
	l.Longitude = longitude

	return nil
}

func (l *GeoSpatialModel) String() string {
	latitude := strconv.FormatFloat(l.Latitude, 'f', -1, 64)
	longitude := strconv.FormatFloat(l.Longitude, 'f', -1, 64)

	return latitude + "," + longitude
}

func (l *GeoSpatialModel) StringReverseOrder() string {
	latitude := strconv.FormatFloat(l.Latitude, 'f', -1, 64)
	longitude := strconv.FormatFloat(l.Longitude, 'f', -1, 64)

	return longitude + "," + latitude
}

type GeoSpatialSearchResult struct {
	Id                string  `db:"id"`
	VendorId          string  `db:"vendorId"`
	VendorName        string  `db:"vendorName"`
	Rate              float32 `db:"rate"`   // For price
	Rating            float32 `db:"rating"` // For rating
	Latitude          float64 `db:"latitude"`
	Longitude         float64 `db:"longitude"`
	CompletedBookings float32 `db:"bookings"` // Number of bookings completed

	Distance float32
}

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
