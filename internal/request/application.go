package request

type NewApplicationPayload struct {
	Job string `json:"job" validate:"required"`
}
