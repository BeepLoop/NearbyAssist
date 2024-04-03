package types

type Message struct {
	Id       int    `query:"id" db:"id"`
	Sender   int    `query:"sender" db:"sender"`
	Receiver int    `query:"receiver" db:"receiver"`
	Content  string `query:"content" db:"content"`
}
