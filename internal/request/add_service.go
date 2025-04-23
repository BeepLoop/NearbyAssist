package request

type AddServicePayload struct {
	VendorId    string     `json:"vendorId" validate:"required"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	Rate        string     `json:"rate" validate:"required"`
	Tags        []string   `json:"tags" validate:"required"`
	Extras      []NewExtra `json:"extras"`
}
