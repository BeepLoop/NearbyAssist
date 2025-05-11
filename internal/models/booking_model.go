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
	VendorId     string         `db:"vendorId"`
	ClientId     string         `db:"clientId"`
	ServiceId    string         `db:"serviceId"`
	Status       BookingStatus  `db:"status"`
	Quantity     int            `db:"quantity"`
	Cost         string         `db:"cost"`
	IsReviewed   bool           `db:"isReviewed"`
	ScheduledAt  sql.NullString `db:"scheduledAt"`
	CancelReason sql.NullString `db:"cancelReason"`
	CancelledBy  sql.NullString `db:"cancelledBy"`

	// Additional fields for joins
	Service *ServiceModel `json:"service,omitempty"`
	Vendor  UserModel     `json:"vendor" db:"vendor"`
	Client  UserModel     `json:"client" db:"client"`
	Extras  []*ExtraModel `json:"extras" db:"extras"`
}
