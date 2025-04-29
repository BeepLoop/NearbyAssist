package models

import "database/sql"

type ExpertiseModel struct {
	Model
	UpdateableModel
	Title string `json:"title" db:"title"`

	Tags               []*TagModel    `json:"tags" db:"tags"`
	DateApplied        string         `db:"dateApplied"`
	DateApproved       sql.NullString `db:"dateApproved"`
	SupportingImageUrl string         `db:"supportingImageUrl"`
}
