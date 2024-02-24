package types

type Message struct {
	Id       int    `query:"id" db:"id"`
	Sender   int    `query:"sender" db:"sender"`
	Reciever int    `query:"reciever" db:"reciever"`
	Content  string `query:"content" db:"content"`
}
