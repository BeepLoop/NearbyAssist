package request

type RejectRequestPayload struct {
	TransactionId string `json:"transactionId" validate:"required"`
	Reason        string `json:"reason" validate:"reason"`
}
