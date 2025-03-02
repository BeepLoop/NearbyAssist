package response

type Analytics struct {
	User               int           `json:"user"`
	Vendor             int           `json:"vendor"`
	VerifiedUser       int           `json:"verifiedUser"`
	PendingApplication int           `json:"pendingApplication"`
	Complaint          int           `json:"complaint"`
	BugReportData      BugReportData `json:"bugReportData"`
}

type BugReportData struct {
	Total int   `json:"total"`
	Daily []int `json:"daily"`
}
