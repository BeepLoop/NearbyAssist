package request

type CancelRequestPayload struct {
	TransactionId string `json:"transactionId"`
	Reason        string `json:"reason"`
}
