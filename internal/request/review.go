package request

type NewReviewPayload struct {
	ServiceId     string `json:"serviceId" db:"serviceId" validate:"required"`
	Rating        int    `json:"rating" db:"rating" validate:"required"`
	TransactionId string `json:"transactionId" db:"transactionId" validate:"required"`
}
