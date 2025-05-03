package models

import "database/sql"

type BookingStatus string

const (
	BOOKING_STATUS_PENDING   BookingStatus = "pending"
	BOOKING_STATUS_CONFIRMED BookingStatus = "confirmed"
	BOOKING_STATUS_REJECTED  BookingStatus = "rejected"
	BOOKING_STATUS_DONE      BookingStatus = "done"
	BOOKING_STATUS_CANCELLED BookingStatus = "cancelled"
)

type BookingModel struct {
	Model
	UpdateableModel
	VendorId     string         `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId     string         `json:"clientId" db:"clientId" validate:"required"`
	ServiceId    string         `json:"serviceId" db:"serviceId" validate:"required"`
	Cost         string         `json:"cost" db:"cost" validate:"required"`
	Status       BookingStatus  `json:"status" db:"status"`
	IsReviewed   bool           `json:"isReviewed" db:"isReviewed"`
	ScheduledAt  sql.NullString `json:"scheduledAt" db:"scheduledAt"`
	CancelReason sql.NullString `json:"cancelReason" db:"cancelReason"`
	CancelledBy  sql.NullString `db:"cancelledBy"`

	// Additional fields for joins
	Service *ServiceModel `json:"service,omitempty"`
	Vendor  string        `json:"vendor,omitempty" db:"vendor"` // Vendor name
	Client  string        `json:"client,omitempty" db:"client"` // Client name
	Extras  []*ExtraModel `json:"extras" db:"extras"`
}
