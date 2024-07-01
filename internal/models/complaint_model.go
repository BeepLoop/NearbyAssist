package models

type ComplaintModel struct {
	Model
	UpdateableModel
	Code    int    `json:"code" db:"code"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}

func NewComplaintModel() *ComplaintModel {
	return &ComplaintModel{}
}

type SystemComplaintModel struct {
	Model
	UpdateableModel
	Title  string `json:"title" db:"title"`
	Detail string `json:"detail" db:"detail"`
}
