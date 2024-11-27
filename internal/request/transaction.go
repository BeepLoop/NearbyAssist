package request

type NewTransactionPayload struct {
	VendorId       string `json:"vendorId" validate:"required"`
	ClientId       string `json:"clientId" validate:"required"`
	ServiceId      string `json:"serviceId" validate:"required"`
	StartDate      string `json:"startDate" validate:"required"`
	EndDate        string `json:"endDate" validate:"required"`
	Cost           string `json:"cost" validate:"required"`
	EmploymentType string `json:"employmentType" validate:"required"`
}
