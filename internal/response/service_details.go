package response

import "nearbyassist/internal/models"

type ServiceDetails struct {
	ServiceId   int    `json:"serviceId" db:"serviceId"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Category    string `json:"category" db:"category"`
	Rate        string `json:"rate" db:"rate"`
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
