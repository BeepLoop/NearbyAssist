package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ModelOperation interface {
	Create() (int, error)
	Update(id int) error
	Delete(id int) error
}

type Locatable interface {
	GetGeolocation() (*GeoSpatialModel, error)
}

type Model struct {
	Id        int    `json:"id" db:"id"`
	CreatedAt string `json:"createdAt" db:"createdAt"`
}

type UpdateableModel struct {
	UpdatedAt string `json:"updatedAt" db:"updatedAt"`
}

type GeoSpatialModel struct {
	Latitude  float64 `json:"latitude" db:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" db:"longitude" validate:"required"`
	Location  string  `json:"location" db:"location"`
}

func ConstructLocationFromLatLong(model *GeoSpatialModel) {
	location := fmt.Sprintf("POINT(%f %f)", model.Latitude, model.Longitude)
	model.Location = location
}

func ExtractLatLongFromLocation(model *GeoSpatialModel) error {
	regex, err := regexp.Compile(`\((.*?)\)`)
	if err != nil {
		return err
	}

	// Extract the latitude and longitude contained within parenthesis
	// Ex. POINT(14.123456 121.123456)
	// Returns [(14.123456 121.123456), 14.123456 121.123456]
	// We only need the second element
	coordinates := regex.FindStringSubmatch(model.Location)[1]

	points := strings.Split(coordinates, " ")
	latitude := points[0]
	longitude := points[1]

	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return err
	}

	long, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return err
	}

	model.Latitude = lat
	model.Longitude = long

	return nil
}
