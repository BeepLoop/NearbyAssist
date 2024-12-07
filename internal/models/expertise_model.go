package models

type ExpertiseModel struct {
	Model
	UpdateableModel
	Title string `json:"title" db:"title"`

	Tags []string `json:"tags" db:"tags"`
}
