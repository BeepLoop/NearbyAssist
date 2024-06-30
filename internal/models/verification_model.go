package models

type IdentityVerificationModel struct {
	Model
	UpdateableModel
	Name     string `json:"name" db:"name" validate:"required"`
	Address  string `json:"address" db:"address" validate:"required"`
	IdType   string `json:"idType" db:"idType" validate:"required"`
	IdNumber string `json:"idNumber" db:"idNumber" validate:"required"`
	FrontId  int    `json:"frontId" db:"frontId" validate:"required"`
	BackId   int    `json:"backId" db:"backId" validate:"required"`
	Face     int    `json:"face" db:"face" validate:"required"`
}

type FrontIdModel struct {
	Model
	UpdateableModel
	Url string `json:"url" db:"url"`
}

type BackIdModel struct {
	Model
	UpdateableModel
	Url string `json:"url" db:"url"`
}

type FaceModel struct {
	Model
	UpdateableModel
	Url string `json:"url" db:"url"`
}
