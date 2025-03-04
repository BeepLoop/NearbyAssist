package models

type RestrictionModel struct {
	UserId    string `db:"userId"`
	Reason    string `db:"reason"`
	StartTime string `db:"startTime"`
	EndTime   string `db:"endTime"`
}
