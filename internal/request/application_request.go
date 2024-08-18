package request

import "nearbyassist/internal/models"

type NewApplication struct {
	ApplicantId string `json:"applicantId" db:"applicantId" validate:"required"`
	Job         string `json:"job" db:"job" validate:"required"`
	models.GeoSpatialModel
	models.Model
}
