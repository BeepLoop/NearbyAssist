package request

type ReportUserPayload struct {
	UserId string `json:"userId" validate:"required"`
	Title  string `json:"title" validate:"required"`
	Reason string `json:"reason" validate:"required"`
}
