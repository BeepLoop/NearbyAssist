package models

type TransactionStatus string
type TransactionFilter string

const (
	TRANSACTION_STATUS_ONGOING   TransactionStatus = "ongoing"
	TRANSACTION_STATUS_DONE      TransactionStatus = "done"
	TRANSACTION_STATUS_CANCELLED TransactionStatus = "cancelled"

	FILTER_CLIENT TransactionFilter = "client"
	FILTER_VENDOR TransactionFilter = "vendor"
)

type TransactionModel struct {
	Model
	UpdateableModel
	VendorId   int               `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId   int               `json:"clientId" db:"clientId" validate:"required"`
	ServiceId  int               `json:"serviceId" db:"serviceId" validate:"required"`
	Start      string            `json:"start" db:"start" validate:"required"`
	End        string            `json:"end" db:"end" validate:"required"`
	Status     TransactionStatus `json:"status" db:"status"`
	IsReviewed bool              `json:"isReviewed" db:"isReviewed"`
}

func NewTransactionModel() *TransactionModel {
	return &TransactionModel{}
}

type DetailedTransactionModel struct {
	Model
	UpdateableModel
	Vendor    string `json:"vendor" db:"vendor"`
	Client    string `json:"client" db:"client"`
	Status    string `json:"status" db:"status"`
	CreatedAt string `json:"createdAt" db:"createdAt"`
}
