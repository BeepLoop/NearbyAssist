package models

import (
	"nearbyassist/internal/utils"
	"strconv"
)

type GeoSpatialModel struct {
	Latitude  float64 `json:"latitude" db:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" db:"longitude" validate:"required"`
}

// String coordinate {latitude},{longitude} ex: 7.544645340539252,126.14141292293003
func (l *GeoSpatialModel) FromString(coordinate string) error {
	latitude, longitude, err := utils.ParseCoordinate(coordinate)
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
	Id                    string  `db:"id"`
	VendorId              string  `db:"vendorId"`
	VendorName            string  `db:"vendorName"`
	Rate                  float32 `db:"rate"`   // for price
	Rating                float32 `db:"rating"` // for rating
	Latitude              float64 `db:"latitude"`
	Longitude             float64 `db:"longitude"`
	CompletedTransactions float32 `db:"transactions"` // number of transactions completed

	Distance float32
}
