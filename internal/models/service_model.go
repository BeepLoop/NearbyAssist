package models

import "database/sql"

type ServiceSearchResult struct {
	ServiceModel
	Vendor string `json:"vendor" db:"vendor"`
}

type PricingType string
type ServiceStatus string

const (
	FIXED_PRICING  PricingType = "fixed"
	HOURLY_PRICING PricingType = "per_hour"
	DAILY_PRICING  PricingType = "per_day"

	SERVICE_STATUS_UNDER_REVIEW ServiceStatus = "under_review"
	SERVICE_STATUS_ACCEPTED     ServiceStatus = "accepted"
	SERVICE_STATUS_REJECTED     ServiceStatus = "rejected"
)

type ServiceModel struct {
	Model
	UpdateableModel
	VendorId     string         `db:"vendorId"`
	Title        string         `db:"title"`
	Description  string         `db:"description"`
	Price        string         `db:"price"`
	PricingType  PricingType    `db:"pricingType"`
	Signature    string         `db:"signature"`
	Disabled     bool           `db:"disabled"`
	Status       ServiceStatus  `db:"status"`
	RejectReason sql.NullString `db:"rejectReason"`
	AcceptedAt   sql.NullString `db:"acceptedAt"`
	RejectedAt   sql.NullString `db:"rejectedAt"`

	// Additional fields for joins
	Vendor       VendorModel
	Tags         []*TagModel          `json:"tags" db:"tags" validate:"required"`
	TagsAsString []string             `json:"-"`
	Extras       []*ExtraModel        `json:"extras" db:"extras"`
	Images       []*ServicePhotoModel `json:"images"`
	Address      AddressModel
}
