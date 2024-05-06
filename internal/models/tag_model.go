package models

type TagModel struct {
	Model
	Title string `json:"title" db:"title"`
}

func NewTagModel() *TagModel {
	return &TagModel{}
}
