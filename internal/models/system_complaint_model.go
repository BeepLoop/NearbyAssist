package models

type SystemComplaintModel struct {
	Model
	UpdateableModel
	Title  string `json:"title" db:"title"`
	Detail string `json:"detail" db:"detail"`

	Images []string `json:"images"`
}
