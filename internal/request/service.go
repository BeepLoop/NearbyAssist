package request

type NewServicePayload struct {
	VendorId    string   `json:"vendorId" validate:"required"`
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Rate        string   `json:"rate" validate:"required"`
	Tags        []string `json:"tags" validate:"required"`
	Latitude    float64  `json:"latitude" validate:"required"`
	Longitude   float64  `json:"longitude" validate:"required"`
	Extras      []struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	} `json:"extras"`
}

type UpdateServicePayload struct {
	Id          string   `json:"id" validate:"required"`
	VendorId    string   `json:"vendorId" validate:"required"`
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Rate        string   `json:"rate" validate:"required"`
	Tags        []string `json:"tags" validate:"required"`
	Latitude    float64  `json:"latitude" validate:"required"`
	Longitude   float64  `json:"longitude" validate:"required"`
}
