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
}
