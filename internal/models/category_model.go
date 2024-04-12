package models

type CategoryModel struct {
	Model
	UpdateableModel
	Title string `json:"title" db:"title"`
}

func NewCategoryModel() *CategoryModel {
	return &CategoryModel{}
}
