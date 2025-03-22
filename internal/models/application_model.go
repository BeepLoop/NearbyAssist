package models

type ApplicationModel struct {
	Model
	UpdateableModel
	GeoSpatialModel
	ApplicantId        string `db:"applicantId" validate:"required"`
	ExpertiseId        string `db:"expertiseId" validate:"required"`
	SupportingDocument string `db:"supportingDocument"`
	PoliceClearance    string `db:"policeClearance"`
	Status             string `db:"status"`
	RejectionReason    string `db:"rejectionReason"`

	// Join table fields
	ApplicantName         string `db:"applicantName"`
	Expertise             string `db:"expertise"`
	SupportingDocumentUrl string `db:"supportingDocumentUrl"`
	PoliceClearanceUrl    string `db:"policeClearanceUrl"`
}
