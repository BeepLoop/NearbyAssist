package models

type MessageModel struct {
	Model
	Sender   string `json:"sender" db:"sender"`
	Receiver string `json:"receiver" db:"receiver"`
	Content  string `json:"content" db:"content"`
}
