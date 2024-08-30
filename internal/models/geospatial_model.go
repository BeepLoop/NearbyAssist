package models

import "strconv"

type GeoSpatialModel struct {
	Latitude  float64 `json:"latitude" db:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" db:"longitude" validate:"required"`
}

func NewGeoGeoSpatialModelWithData(lat, long float64) *GeoSpatialModel {
	return &GeoSpatialModel{
		Latitude:  lat,
		Longitude: long,
	}
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
