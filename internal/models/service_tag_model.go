package models

type ServiceTagModel struct {
	Id        int    `db:"id"`
	ServiceId int    `db:"serviceId"`
	Tag       string `db:"tag"`
}
