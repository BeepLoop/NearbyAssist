package models

import "database/sql"

type ApplicationModel struct {
	Model
	GeoSpatialModel
	ApplicantId        string         `db:"applicantId" validate:"required"`
	ExpertiseId        string         `db:"expertiseId" validate:"required"`
	SupportingDocument string         `db:"supportingDocument"`
	PoliceClearance    string         `db:"policeClearance"`
	Status             string         `db:"status"`
	RejectionReason    string         `db:"rejectionReason"`
	UpdatedAt          sql.NullString `db:"updatedAt"`

	// Join table fields
	ApplicantName         string `db:"applicantName"`
	Expertise             string `db:"expertise"`
	SupportingDocumentUrl string `db:"supportingDocumentUrl"`
	PoliceClearanceUrl    string `db:"policeClearanceUrl"`
}
