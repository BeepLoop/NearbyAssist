package models

type ReportedUserModel struct {
	Model
	UserId string `db:"userId"`
	Reason string `db:"reason"`
	Detail string `db:"detail"`

	Images []string
}
