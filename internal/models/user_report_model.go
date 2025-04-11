package models

import "database/sql"

type UserReportCategory string
type UserReportStatus string

var (
	category_misconduct      UserReportCategory = "misconduct"
	category_booking_related UserReportCategory = "booking_related"
)

type UserReportModel struct {
	Model
	UpdateableModel
	ReporterUserId string             `db:"reporterUserId"`
	ReportedUserId string             `db:"reportedUserId"`
	Category       UserReportCategory `db:"category"`
	BookingId      sql.NullString     `db:"bookingId"`
	Reason         string             `db:"reason"`
	Detail         string             `db:"detail"`
	Status         string             `db:"status"`

	Images         []string
	BookingIdInput string
}
