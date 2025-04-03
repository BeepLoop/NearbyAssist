package request

type CancelRequestPayload struct {
	TransactionId string `json:"transactionId" validate:"required"`
	Reason        string `json:"reason" validate:"required"`
}
