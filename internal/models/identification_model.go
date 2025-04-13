package models

type IdentificationModel struct {
	Type       string `db:"type"`
	IdNumber   string `db:"idNumber"`
	FrontImage string `db:"frontImage"`
	BackImage  string `db:"backImage"`
}
