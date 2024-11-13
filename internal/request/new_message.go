package request

type NewMessagePayload struct {
	Sender   string `json:"sender" validate:"required"`
	Receiver string `json:"receiver" validate:"required"`
	Content  string `json:"content" validate:"required"`
}
