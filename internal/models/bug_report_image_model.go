package models

type BugReportImageModel struct {
	Model
	UpdateableModel
	ComplaintId string `json:"complaintId" db:"complaintId"`
	Url         string `json:"url" db:"url"`
}
