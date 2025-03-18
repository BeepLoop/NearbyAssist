package models

type IdentityVerificationModel struct {
	Model
	UpdateableModel
	GeoSpatialModel
	UserId          string `json:"userId" db:"userId" validate:"required"`
	Name            string `json:"name" db:"name" validate:"required"`
	Address         string `json:"address" db:"address" validate:"required"`
	Phone           string `json:"phone" db:"phone" validate:"required"`
	IdType          string `json:"idType" db:"idType" validate:"required"`
	IdNumber        string `json:"idNumber" db:"idNumber" validate:"required"`
	FrontIdImageUrl string `json:"frontIdImageUrl" db:"frontIdImageUrl" validate:"required"`
	BackIdImageUrl  string `json:"backIdImageUrl" db:"backIdImageUrl" validate:"required"`
	FaceImageUrl    string `json:"faceImageUrl" db:"faceImageUrl" validate:"required"`
	Status          string `json:"status" db:"status"`
	RejectionReason string `db:"rejectionReason"`
}
