package models

type PoliceClearanceModel struct {
	Model
	UpdateableModel
	ApplicationId string `db:"applicationId"`
	ApplicantId   string `db:"applicantId"`
	Url           string `db:"url"`
}
