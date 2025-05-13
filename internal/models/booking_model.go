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
	VendorId           string         `db:"vendorId"`
	ClientId           string         `db:"clientId"`
	ServiceId          string         `db:"serviceId"`
	ServiceTitle       string         `db:"serviceTitle"`
	ServiceDescription string         `db:"serviceDescription"`
	Price              string         `db:"price"`
	PricingType        PricingType    `db:"pricingType"`
	Status             BookingStatus  `db:"status"`
	Quantity           int            `db:"quantity"`
	Cost               string         `db:"cost"`
	IsReviewed         bool           `db:"isReviewed"`
	ScheduleStart      sql.NullString `db:"scheduleStart"`
	ScheduleEnd        sql.NullString `db:"scheduleEnd"`
	CancelReason       sql.NullString `db:"cancelReason"`
	CancelledBy        sql.NullString `db:"cancelledBy"`

	// Additional fields for joins
	Vendor UserModel            `json:"vendor" db:"vendor"`
	Client UserModel            `json:"client" db:"client"`
	Extras []*BookingExtraModel `json:"extras" db:"extras"`
}

type BookingExtraModel struct {
	BookingId        string `db:"bookingId"`
	ExtraTitle       string `db:"extraTitle"`
	ExtraDescription string `db:"extraDescription"`
	Price            string `db:"price"`
	CreatedAt        string `db:"createdAt"`
}
