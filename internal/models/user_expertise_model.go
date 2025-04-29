package models

import "database/sql"

type UserExpertiseModel struct {
	Model
	UpdateableModel
	UserId             string `db:"userId"`
	ExpertiseId        string `db:"expertiseId"`
	SupportingDocument string `db:"supportingImage"`

	Expertise               string         `db:"expertise"`
	DateApplied             string         `db:"dateApplied"`
	DateApproved            sql.NullString `db:"dateApproved"`
	SupportingDocumentImage string         `db:"supportingDocumentImage"`
}
