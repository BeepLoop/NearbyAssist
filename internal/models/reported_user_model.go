package models

type ReportedUserModel struct {
	Model
	UserId string `json:"userId" db:"userId"`
	Reason string `json:"reason" db:"reason"`
	Detail string `json:"detail" db:"detail"`

	Images []string `json:"images"`
	Name   string   `json:"name" db:"name"`
}
