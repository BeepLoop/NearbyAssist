package types

type VendorApplication struct {
	ApplicantId int     `db:"applicantId" json:"applicantId" validate:"required"`
	Job         string  `db:"job" json:"job" validate:"required"`
	Longitude   float64 `db:"longitude" json:"longitude" validate:"required"`
	Latitude    float64 `db:"latitude" json:"latitude" validate:"required"`
}
