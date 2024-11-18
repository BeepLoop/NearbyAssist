package request

type SaveServicePayload struct {
	ServiceId string `json:"serviceId" validate:"required"`
}

type UnsaveServicePayload struct {
	ServiceId string `json:"serviceId" validate:"required"`
}
