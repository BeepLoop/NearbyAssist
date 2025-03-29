package request

type NewMessagePayload struct {
	Id       string `json:"id" validate:"required"`
	Sender   string `json:"sender" validate:"required"`
	Receiver string `json:"receiver" validate:"required"`
	Content  string `json:"content" validate:"required"`
}
