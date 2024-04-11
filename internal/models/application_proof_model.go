package models

import "path/filepath"

type ApplicationProofModel struct {
	Model
	UpdateableModel
	ApplicationId int    `json:"applicationId" db:"applicationId"`
	ApplicantId   int    `json:"applicantId" db:"applicantId"`
	Url           string `json:"url" db:"url"`
}

func NewApplicationProofModel(applicationId, applicantId int, filename string) *ApplicationProofModel {
	fileLocation := filepath.Join("/resource/proofs", filename)

	return &ApplicationProofModel{
		ApplicationId: applicationId,
		ApplicantId:   applicantId,
		Url:           fileLocation,
	}
}
