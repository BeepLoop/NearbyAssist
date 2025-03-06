package models

type ReportedUserModel struct {
	Model
	UserId string `db:"userId"`
	Title  string `db:"title"`
	Reason string `db:"reason"`

	Images []string
}
