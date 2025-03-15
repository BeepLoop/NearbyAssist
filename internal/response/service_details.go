package response

import "nearbyassist/internal/models"

type CountPerRating map[string]int

func NewCountPerRating() CountPerRating {
	instance := make(CountPerRating)
	instance["five"] = 0
	instance["four"] = 0
	instance["three"] = 0
	instance["two"] = 0
	instance["one"] = 0

	return instance
}

type DetailedServiceResponse struct {
	Vendor         *models.VendorModel  `json:"vendor"`
	Service        *models.ServiceModel `json:"service"`
	CountPerRating CountPerRating       `json:"countPerRating"`
}
