package models

type ComplaintModel struct {
	Model
	UpdateableModel
	Title  string `json:"title" db:"title"`
	Detail string `json:"content" db:"detail"`
}
