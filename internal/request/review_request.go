package request

import "nearbyassist/internal/models"

type NewReview struct {
	ServiceId     string `json:"serviceId" db:"serviceId" validate:"required"`
	TransactionId string `json:"transactionId" validate:"required"`
	Rating        string `json:"rating" db:"rating" validate:"required"`
	models.Model
}
