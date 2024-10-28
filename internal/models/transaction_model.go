package models

type TransactionStatusFilter string

const (
	TRANSACTION_STATUS_PENDING   TransactionStatusFilter = "pending"
	TRANSACTION_STATUS_ONGOING   TransactionStatusFilter = "ongoing"
	TRANSACTION_STATUS_DONE      TransactionStatusFilter = "done"
	TRANSACTION_STATUS_CANCELLED TransactionStatusFilter = "cancelled"
)

type TransactionModel struct {
	Model
	UpdateableModel
	VendorId    string                  `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId    string                  `json:"clientId" db:"clientId" validate:"required"`
	ServiceId   string                  `json:"serviceId" db:"serviceId" validate:"required"`
	Price       string                  `json:"price" db:"price" validate:"required"`
	StartDate   string                  `json:"startDate" db:"startDate" validate:"required"`
	EndDate     string                  `json:"endDate" db:"endDate" validate:"required"`
	Status      TransactionStatusFilter `json:"status" db:"status"`
	IsReviewed  bool                    `json:"isReviewed" db:"isReviewed"`
	IsReported  bool                    `json:"isReported" db:"isReported"`
	ConfirmCode string                  `json:"confirmCode" db:"confirmCode"`

	// Additional fields for joins
	Vendor string `json:"vendor" db:"vendor"` // Vendor name
	Client string `json:"client" db:"client"` // Client name
}
