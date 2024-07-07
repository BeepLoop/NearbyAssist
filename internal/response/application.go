package response

import "nearbyassist/internal/models"

type Application struct {
	Id          int                      `json:"id" db:"id"`
	ApplicantId int                      `json:"applicantId" db:"applicantId"`
	Status      models.ApplicationStatus `json:"status" db:"status"`
	CreatedAt   string                   `json:"createdAt" db:"createdAt"`
}
