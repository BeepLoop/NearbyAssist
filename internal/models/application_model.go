package models

type ApplicationStatus string

const (
	APPLICATION_STATUS_PENDING  ApplicationStatus = "pending"
	APPLICATION_STATUS_APPROVED ApplicationStatus = "approved"
	APPLICATION_STATUS_REJECTED ApplicationStatus = "rejected"
)

type ApplicationModel struct {
	Model
	UpdateableModel
	GeoSpatialModel
	ApplicantId int               `json:"applicantId" db:"applicantId" validate:"required"`
	Job         string            `json:"job" db:"job" validate:"required"`
	Status      ApplicationStatus `json:"status" db:"status"`
}

func NewApplicationModel() *ApplicationModel {
	return &ApplicationModel{}
}
