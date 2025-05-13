package dto

type GeospatialOperation struct {
	Id                 string
	Price              float32
	Rating             float32
	Latitude           float32
	Longitude          float32
	CompletedBookings  float32
	DistanceFromOrigin float32
}

type ServiceWithDistance struct {
	Id                 string
	Price              float64
	Rating             float64
	DistanceFromOrigin float64
	CompletedBookings  float64
}
