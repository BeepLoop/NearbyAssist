package request

type NewTransactionPayload struct {
	VendorId    string  `json:"vendorId" validate:"required"`
	ClientId    string  `json:"clientId" validate:"required"`
	ServiceId   string  `json:"serviceId" validate:"required"`
	Cost        string  `json:"cost" validate:"required"`
	Extras      []Extra `json:"extras" validate:"required"`
	ScheduledAt string  `json:"scheduledAt" validate:"required"`
}

type Extra struct {
	Id          string  `json:"id" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
}
