package models

type TransactionStatusFilter string

const (
	TRANSACTION_STATUS_PENDING   TransactionStatusFilter = "pending"
	TRANSACTION_STATUS_CONFIRMED TransactionStatusFilter = "confirmed"
	TRANSACTION_STATUS_REJECTED  TransactionStatusFilter = "rejected"
	TRANSACTION_STATUS_DONE      TransactionStatusFilter = "done"
	TRANSACTION_STATUS_CANCELLED TransactionStatusFilter = "cancelled"
)

type TransactionModel struct {
	Model
	UpdateableModel
	VendorId   string                  `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId   string                  `json:"clientId" db:"clientId" validate:"required"`
	ServiceId  string                  `json:"serviceId" db:"serviceId" validate:"required"`
	Cost       string                  `json:"cost" db:"cost" validate:"required"`
	Status     TransactionStatusFilter `json:"status" db:"status"`
	IsReviewed bool                    `json:"isReviewed" db:"isReviewed"`

	// Additional fields for joins
	Service *ServiceModel `json:"service,omitempty"`
	Vendor  string        `json:"vendor,omitempty" db:"vendor"` // Vendor name
	Client  string        `json:"client,omitempty" db:"client"` // Client name
	Extras  []*ExtraModel `json:"extras" db:"extras"`
}
