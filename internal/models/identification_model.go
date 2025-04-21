package models

type IdentificationModel struct {
	Id              string `db:"id"`
	Type            string `db:"type"`
	ReferenceNumber string `db:"referenceNumber"`
	FrontImageUrl   string `db:"frontImageUrl"`
	BackImageUrl    string `db:"backImageUrl"`
	SelfieImageUrl  string `db:"selfieImageUrl"`
	CreatedAt       string `db:"createdAt"`
}
