package models

import "database/sql"

type UserReportCategory string
type UserReportStatus string

var (
	CATEGORY_MISCONDUCT      UserReportCategory = "misconduct"
	CATEGORY_BOOKING_RELATED UserReportCategory = "booking_related"

	REPORT_STATUS_PENDING   UserReportStatus = "pending"
	REPORT_STATUS_RESOLVED  UserReportStatus = "resolved"
	REPORT_STATUS_DISMISSED UserReportStatus = "dismissed"
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
	Status         UserReportStatus   `db:"status"`

	Images         []string
	BookingIdInput string
}
