package response

import "nearbyassist/internal/models"

type ServiceDetails struct {
	ServiceId   int      `json:"serviceId" db:"serviceId"`
	Description string   `json:"description" db:"description"`
	Tags        []string `json:"tags" db:"tags"`
	Rate        string   `json:"rate" db:"rate"`
	models.GeoSpatialModel
}

type ServiceVendorDetails struct {
	VendorId int    `json:"vendorId" db:"vendorId"`
	Vendor   string `json:"vendor" db:"vendor"`
	ImageUrl string `json:"imageUrl" db:"imageUrl"`
	Rating   string `json:"rating" db:"rating"`
	Job      string `json:"job" db:"job"`
}

type ServiceImages struct {
	ImageId  int    `json:"imageId" db:"imageId"`
	ImageUrl string `json:"imageUrl" db:"imageUrl"`
}

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
