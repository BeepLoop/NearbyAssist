package dto

type UserReport struct {
	Id               string
	ReportedUserId   string
	ReportedByUserId string
	CreatedAt        string
}

type UserReportDetail struct {
	Reporter User
	Reported User
	Report   Report
}

type Report struct {
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
