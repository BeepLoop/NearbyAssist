package models

import "database/sql"

type BugReportModel struct {
	Id          int            `json:"id" db:"id"`
	Title       string         `json:"title" db:"title"`
	Detail      string         `json:"detail" db:"detail"`
	CreatedAt   string         `json:"createdAt,omitempty" db:"createdAt"`
	CompletedAt sql.NullString `json:"completedAt,omitempty" db:"completedAt"`
	Images      []string       `json:"images"`
}
