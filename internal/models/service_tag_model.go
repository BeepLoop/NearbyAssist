package models

type ServiceTagModel struct {
	Model
	UpdateableModel
	ServiceId string `json:"serviceId" db:"serviceId"`
	TagId     string `json:"tagId" db:"tagId"`

	// Additional fields for joining tag title
	Title string `json:"title" db:"title"`
}
