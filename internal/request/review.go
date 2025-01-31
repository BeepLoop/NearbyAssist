package request

type NewReviewPayload struct {
	ServiceId     string `json:"serviceId" db:"serviceId" validate:"required"`
	Rating        int    `json:"rating" db:"rating" validate:"required"`
	Text          string `json:"text" db:"text" validate:"required"`
	TransactionId string `json:"transactionId" db:"transactionId" validate:"required"`
}
