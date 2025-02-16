package models

type TagModel struct {
	Model
	UpdateableModel
	Title string `json:"title" db:"title"`
}
