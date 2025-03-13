package models

import "database/sql"

type ReportedUserModel struct {
	Model
	ReportedBy  string         `json:"reportedBy" db:"reportedBy"`
	UserId      string         `json:"userId" db:"userId"`
	Reason      string         `json:"reason" db:"reason"`
	Detail      string         `json:"detail" db:"detail"`
	CompletedAt sql.NullString `json:"completedAt" db:"completedAt"`

	Images []string `json:"images,omitempty"`
	Name   string   `json:"name" db:"name"`
}
