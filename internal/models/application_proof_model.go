package models

import "path/filepath"

type ApplicationProofModel struct {
	Model
	UpdateableModel
	ApplicationId string `json:"applicationId" db:"applicationId"`
	ApplicantId   string `json:"applicantId" db:"applicantId"`
	Url           string `json:"url" db:"url"`
}

func NewApplicationProofModel(applicationProofId, applicationId, applicantId string, filename string) *ApplicationProofModel {
	fileLocation := filepath.Join("/resource/proofs", filename)

	return &ApplicationProofModel{
		Model:         Model{Id: applicationProofId},
		ApplicationId: applicationId,
		ApplicantId:   applicantId,
		Url:           fileLocation,
	}
}
