package models

type ReviewModel struct {
	Model
	UpdateableModel
	ServiceId int    `json:"serviceId" db:"serviceId"`
	Rating    string `json:"rating" db:"rating"`
}

func NewReviewModel() *ReviewModel {
	return &ReviewModel{}
}
