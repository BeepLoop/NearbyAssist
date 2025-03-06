package models

type BugReportData struct {
	ComplaintId string
	Url         string
}

type BugReportImageModel struct {
	Model
	UpdateableModel
	ComplaintId string `json:"complaintId" db:"complaintId"`
	Url         string `json:"url" db:"url"`
}
