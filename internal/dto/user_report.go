package dto

type ReportItem struct {
	Id        string
	Reported  ReportUser
	Reporter  ReportUser
	CreatedAt string
}

type ReportUser struct {
	Id    string
	Name  string
	Email string
}

type UserReportDetail struct {
	Reporter            User
	Reported            User
	ReporterHistory     ReporterHistory
	ReportedUserHistory ReportedUserHistory
	Report              ReportDetail
	Booking             Booking
}

type ReportedUserHistory struct {
	Bookings          int
	CompletedBookings int
	RejectedBookings  int
	ActiveBookings    int
	PreviousReports   []PreviousReport
	AccountCreatedAt  string
	JoinedVendorAt    string
	Rating            string
}

type ReporterHistory struct {
	ReportsFiled      int
	FalseReports      int
	CancelledBookings int
	AccountCreatedAt  string
}

type ReportDetail struct {
	Id               string
	ReportedUserId   string
	ReportedByUserId string
	Category         string
	BookingId        string
	Reason           string
	Detail           string
	Images           []string
	Status           string
	CreatedAt        string
	CompletedAt      string
}

type PreviousReport struct {
	Id               string
	ReportedByUserId string
	ReportedByName   string
	Category         string
	BookingId        string
	Reason           string
	Detail           string
	Images           []string
	Status           string
	AdminId          string
	AdminUsername    string
	AdminNote        string
	CreatedAt        string
	CompletedAt      string
}
