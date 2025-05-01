package models

import "database/sql"

type MessageModel struct {
	Model
	Sender   string         `db:"sender"`
	Receiver string         `db:"receiver"`
	Content  string         `db:"content"`
	Seen     bool           `db:"seen"`
	SeenAt   sql.NullString `db:"seenAt"`
}
