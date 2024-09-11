package request

type NewApplicationPayload struct {
	ApplicantId string `json:"applicantId" validate:"required"`
	Job         string `json:"job" validate:"required"`
}
