package models

import "strconv"

type Location struct {
	Latitude  float64 `json:"latitude" db:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" db:"longitude" validate:"required"`
}

func NewLocationWithData(lat, long float64) *Location {
	return &Location{
		Latitude:  lat,
		Longitude: long,
	}
}

func (l *Location) String() string {
	latitude := strconv.FormatFloat(l.Latitude, 'f', -1, 64)
	longitude := strconv.FormatFloat(l.Longitude, 'f', -1, 64)

	return latitude + "," + longitude
}

func (l *Location) StringReverseOrder() string {
	latitude := strconv.FormatFloat(l.Latitude, 'f', -1, 64)
	longitude := strconv.FormatFloat(l.Longitude, 'f', -1, 64)

	return longitude + "," + latitude
}
