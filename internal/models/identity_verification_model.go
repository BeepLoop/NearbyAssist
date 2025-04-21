package models

import "database/sql"

type IdentityVerificationStatus string

const (
	IDENTITY_VERIF_STATUS_PENDING  = "pending"
	IDENTITY_VERIF_STATUS_APPROVED = "approved"
	IDENTITY_VERIF_STATUS_REJECTED = "rejected"
)

type IdentityVerificationModel struct {
	Id            string                     `db:"id"`
	UserId        string                     `db:"userId"`
	Status        IdentityVerificationStatus `db:"status"`
	RejectionNote sql.NullString             `db:"rejectionNote"`
	CreatedAt     string                     `db:"createdAt"`
	UpdatedAt     string                     `db:"updatedAt"`
	User          UserModel
}
