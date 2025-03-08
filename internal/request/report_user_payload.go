package request

type ReportUserPayload struct {
	UserId string `json:"userId" validate:"required"`
	Reason string `json:"reason" validate:"required"`
	Detail string `json:"detail" validate:"required"`
}
