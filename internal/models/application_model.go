package models

type ApplicationStatusFilter string

const (
	APPLICATION_STATUS_ALL      ApplicationStatusFilter = "all"
	APPLICATION_STATUS_PENDING  ApplicationStatusFilter = "pending"
	APPLICATION_STATUS_APPROVED ApplicationStatusFilter = "approved"
	APPLICATION_STATUS_REJECTED ApplicationStatusFilter = "rejected"
)

type ApplicationModel struct {
	Model
	UpdateableModel
	GeoSpatialModel
	ApplicantId           string                  `db:"applicantId" validate:"required"`
	ExpertiseId           string                  `db:"expertiseId" validate:"required"`
	SupportingDocumentUrl string                  `db:"supportingDocumentUrl"`
	PoliceClearanceUrl    string                  `db:"policeClearanceUrl"`
	Status                ApplicationStatusFilter `db:"status"`
}
