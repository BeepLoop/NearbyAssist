package dto

type UserReport struct {
	Reporter User
	Reported User
	Report   Report
}

type Report struct {
	Id               string
	ReportedUserId   string
	ReportedByUserId string
	Reason           string
	Detail           string
	Images           []string
	Status           string
	CreatedAt        string
	CompletedAt      string
}
