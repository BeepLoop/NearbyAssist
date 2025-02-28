package models

import "database/sql"

type UserAccountPageData struct {
	Id         string
	ProfileURL string
	Name       string
	Email      string
	Address    sql.NullString
	Banned     bool
	CreatedAt  string

	Expertise []string
	Services  []*ServiceModel

	Stat UserTransactionStats
}

type UserTransactionStats struct {
	Sent     SentStat
	Received ReceivedStat
}

type SentStat struct {
	CurrentMonth int `db:"currentMonth"`
	LastMonth    int `db:"lastMonth"`
}

type ReceivedStat struct {
	CurrentMonth int `db:"currentMonth"`
	LastMonth    int `db:"lastMonth"`
}
