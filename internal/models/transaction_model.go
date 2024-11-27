package models

type TransactionStatusFilter string
type EmploymentType string

const (
	TRANSACTION_STATUS_PENDING   TransactionStatusFilter = "pending"
	TRANSACTION_STATUS_ONGOING   TransactionStatusFilter = "ongoing"
	TRANSACTION_STATUS_DONE      TransactionStatusFilter = "done"
	TRANSACTION_STATUS_CANCELLED TransactionStatusFilter = "cancelled"

	EMPLOYMENT_TYPE_ARAWAN EmploymentType = "arawan"
	EMPLOYMENT_TYPE_PAKYAW EmploymentType = "pakyaw"
)

type TransactionModel struct {
	Model
	UpdateableModel
	VendorId       string                  `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId       string                  `json:"clientId" db:"clientId" validate:"required"`
	ServiceId      string                  `json:"serviceId" db:"serviceId" validate:"required"`
	Cost           string                  `json:"cost" db:"cost" validate:"required"`
	StartDate      string                  `json:"startDate" db:"startDate" validate:"required"`
	EndDate        string                  `json:"endDate" db:"endDate" validate:"required"`
	EmploymentType EmploymentType          `json:"employmentType" db:"employmentType" validate:"required"`
	Status         TransactionStatusFilter `json:"status" db:"status"`
	IsReviewed     bool                    `json:"isReviewed" db:"isReviewed"`
	IsReported     bool                    `json:"isReported" db:"isReported"`

	// Additional fields for joins
	Vendor string `json:"vendor" db:"vendor"` // Vendor name
	Client string `json:"client" db:"client"` // Client name
}
