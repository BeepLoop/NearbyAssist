package models

import "database/sql"

type NotificationModel struct {
	Model
	Recipient string         `json:"recipient" db:"recipient"`
	Type      string         `json:"type" db:"type"`
	Title     string         `json:"title" db:"title"`
	Content   string         `json:"content" db:"content"`
	ReadAt    sql.NullString `json:"readAt" db:"readAt"`
}
