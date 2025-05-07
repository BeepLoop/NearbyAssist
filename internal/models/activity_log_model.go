package models

import "database/sql"

type ActivityLogModel struct {
	Id         string         `db:"id"`
	AdminId    string         `db:"adminId"`
	Action     string         `db:"action"`
	TargetType string         `db:"targetType"`
	TargetId   sql.NullString `db:"targetId"`
	CreatedAt  string         `db:"createdAt"`

	AdminUsername string
	Target        string
}
