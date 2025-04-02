package models

type ReviewModel struct {
	Model
	TransactionId string `json:"transactionId" db:"transactionId" validate:"required"`
	Rating        int    `json:"rating" db:"rating" validate:"required"`
	Text          string `json:"text" db:"text" validate:"required"`
}
