package types

type ServicePhoto struct {
	VendorId  int    `json:"vendorId" db:"vendorId" validate:"required"`
	ServiceId int    `json:"serviceId" db:"serviceId" validate:"required"`
	Url       string `db:"url"`
}

type ApplicationProof struct {
	ApplicationId int    `json:"applicationId" db:"applicationId" validate:"required"`
	ApplicantId   int    `json:"applicantId" db:"applicantId" validate:"required"`
	Url           string `db:"url"`
}
