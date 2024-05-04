package models

type ReviewModel struct {
	Model
	UpdateableModel
	ServiceId int `json:"serviceId" db:"serviceId" validate:"required"`
	Rating    int `json:"rating" db:"rating" validate:"required"`
}

func NewReviewModel() *ReviewModel {
	return &ReviewModel{}
}
