package models

import "database/sql"

type TransactionStatus string

const (
	TRANSACTION_STATUS_PENDING   TransactionStatus = "pending"
	TRANSACTION_STATUS_CONFIRMED TransactionStatus = "confirmed"
	TRANSACTION_STATUS_REJECTED  TransactionStatus = "rejected"
	TRANSACTION_STATUS_DONE      TransactionStatus = "done"
	TRANSACTION_STATUS_CANCELLED TransactionStatus = "cancelled"
)

type TransactionModel struct {
	Model
	UpdateableModel
	VendorId     string            `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId     string            `json:"clientId" db:"clientId" validate:"required"`
	ServiceId    string            `json:"serviceId" db:"serviceId" validate:"required"`
	Cost         string            `json:"cost" db:"cost" validate:"required"`
	Status       TransactionStatus `json:"status" db:"status"`
	IsReviewed   bool              `json:"isReviewed" db:"isReviewed"`
	ScheduledAt  sql.NullString    `json:"scheduledAt" db:"scheduledAt"`
	CancelReason sql.NullString    `json:"cancelReason" db:"cancelReason"`

	// Additional fields for joins
	Service *ServiceModel `json:"service,omitempty"`
	Vendor  string        `json:"vendor,omitempty" db:"vendor"` // Vendor name
	Client  string        `json:"client,omitempty" db:"client"` // Client name
	Extras  []*ExtraModel `json:"extras" db:"extras"`
}
