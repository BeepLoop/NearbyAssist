package models

import "database/sql"

type ExpertiseModel struct {
	Model
	UpdateableModel
	Title string `json:"title" db:"title"`

	DateApplied        string         `db:"dateApplied"`
	DateApproved       sql.NullString `db:"dateApproved"`
	SupportingImageUrl string         `db:"supportingImageUrl"`
}
