package request

type UpdateServicePayload struct {
	Id          string   `json:"id" validate:"required"`
	VendorId    string   `json:"vendorId" validate:"required"`
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Price       string   `json:"price" validate:"required"`
	PricingType string   `json:"pricingType" validate:"required"`
	Tags        []string `json:"tags" validate:"required"`
}
