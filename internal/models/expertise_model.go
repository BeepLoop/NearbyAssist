package models

type ExpertiseModel struct {
	Model
	UpdateableModel
	ExpertiseId string `json:"expertiseId" db:"expertiseId"`
	Title       string `json:"title" db:"title"`

	Tags []string `json:"tags" db:"tags"`
}
