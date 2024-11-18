package models

type SavedServiceModel struct {
	Model
	UpdateableModel
	UserId    string `db:"userId"`
	ServiceId string `db:"serviceId"`
}
