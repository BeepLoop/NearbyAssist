package geodistance

import (
	"math"
)

const (
	DegreesInPiRadian = 180
	MetersPerKm       = 1000
	EarthRadiusMiles  = 3958
	EarthRadiusKm     = 6371
	EarthRadiusMeters = EarthRadiusKm * MetersPerKm

	// kilometers unit
	KM Unit = "kilometers"

	// meters unit
	M Unit = "meters"

	// miles unit
	MI Unit = "miles"
)

type Unit string
type Distance float64
type Kilometers Distance
type Meters Distance
type Miles Distance

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type Delta Coordinate

func (coordinate Coordinate) ToRadians() Coordinate {
	return Coordinate{
		Latitude:  degreesToRadians(coordinate.Latitude),
		Longitude: degreesToRadians(coordinate.Longitude),
	}
}

func (coordinate Coordinate) DistanceTo(remote Coordinate, unit Unit) Distance {
	return toUnit(haversineDistance(coordinate, remote), unit)
}

func (coordinate Coordinate) Delta(origin Coordinate) Delta {
	return Delta{
		Latitude:  coordinate.Latitude - origin.Latitude,
		Longitude: coordinate.Longitude - origin.Longitude,
	}
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / DegreesInPiRadian
}

func haversineDistance(origin, remote Coordinate) Distance {
	origin, remote = origin.ToRadians(), remote.ToRadians()
	change := origin.Delta(remote)

	angle := math.Pow(math.Sin(change.Latitude/2), 2) + math.Cos(origin.Latitude)*math.Cos(remote.Latitude)*
		math.Pow(math.Sin(change.Longitude/2), 2)

	return Distance(2 * math.Atan2(math.Sqrt(angle), math.Sqrt(1-angle)))
}

func toUnit(distance Distance, unit Unit) Distance {
	var result Distance
	switch unit {
	case KM:
		result = Distance(Kilometers(distance * EarthRadiusKm))
	case M:
		result = Distance(Meters(distance * EarthRadiusMeters))
	case MI:
		result = Distance(Miles(distance * EarthRadiusMiles))
	}
	return result
}
