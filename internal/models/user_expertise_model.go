package models

type UserExpertiseModel struct {
	Model
	UpdateableModel
	UserId             string `db:"userId"`
	ExpertiseId        string `db:"expertiseId"`
	SupportingDocument string `db:"supportingImage"`

	Expertise               string `db:"expertise"`
	DateApplied             string `db:"dateApplied"`
	DateApproved            string `db:"dateApproved"`
	SupportingDocumentImage string `db:"supportingDocumentImage"`
}
