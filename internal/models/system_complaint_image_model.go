package models

type SystemComplaintImageModel struct {
	ComplaintId int    `json:"complaintId" db:"complaintId"`
	Url         string `json:"url" db:"url"`
	Model
	UpdateableModel
}
