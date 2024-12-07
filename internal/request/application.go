package request

type NewApplicationPayload struct {
	ExpertiseId string `json:"expertiseId" validate:"required"`
}
