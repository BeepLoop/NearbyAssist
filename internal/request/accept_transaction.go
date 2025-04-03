package request

type AcceptTransactionPayload struct {
	TransactionId string `json:"transactionId" validate:"required"`
	Schedule      string `json:"schedule" validate:"required"`
}
