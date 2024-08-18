package models

type SystemComplaintImageModel struct {
	ComplaintId string `json:"complaintId" db:"complaintId"`
	Url         string `json:"url" db:"url"`
	Model
	UpdateableModel
}
