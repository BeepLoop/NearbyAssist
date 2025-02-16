package models

import "database/sql"

type UserAccountPageData struct {
	Id         string
	ProfileURL string
	Name       string
	Email      string
	Address    sql.NullString
	CreatedAt  string

	Expertise []string
	Services  []*ServiceModel
}
