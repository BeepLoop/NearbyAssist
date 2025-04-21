package models

import "database/sql"

type IdentityVerificationModel struct {
	Id            string         `db:"id"`
	UserId        string         `db:"userId"`
	Status        string         `db:"status"`
	RejectionNote sql.NullString `db:"rejectionNote"`
	CreatedAt     string         `db:"createdAt"`
	UpdatedAt     string         `db:"updatedAt"`
	User          UserModel
}
