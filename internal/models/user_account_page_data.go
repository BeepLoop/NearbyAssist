package models

import "database/sql"

type UserAccountPageData struct {
	Id         string
	ProfileURL string
	Name       string
	Email      string
	Address    sql.NullString
	Banned     bool
	Restricted bool
	CreatedAt  string
	Verified   bool
	VerifiedAt string

	Expertise []string
	Services  []*ServiceModel

	Stat UserBookingStats
}

type UserBookingStats struct {
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
