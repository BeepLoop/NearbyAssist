package dto

type GeospatialOperation struct {
	Id                 string
	Rate               float32
	Rating             float32
	Latitude           float32
	Longitude          float32
	CompletedBookings  float32
	DistanceFromOrigin float32
}

type DistanceCalculationResult struct {
	ServiceID string
	Distance  float32
}
