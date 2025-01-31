package models

type ReviewModel struct {
	Model
	UpdateableModel
	ServiceId string `json:"serviceId" db:"serviceId" validate:"required"`
	Rating    int    `json:"rating" db:"rating" validate:"required"`
	Text      string `json:"text" db:"text" validate:"required"`

	// Additional fields for creating a review
	TransactionId string `json:"transactionId" db:"transactionId" validate:"required"`
}
