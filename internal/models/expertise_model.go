package models

type ExpertiseModel struct {
	Model
	UpdateableModel
	Title string `json:"title" db:"title"`

	Tags []*TagModel `json:"tags" db:"tags"`
}
