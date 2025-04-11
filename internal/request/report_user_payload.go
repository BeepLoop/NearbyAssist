package request

type ReportUserPayload struct {
	UserId    string `json:"userId" validate:"required"`
	Category  string `json:"category" validate:"required"`
	BookingId string `json:"bookingId"`
	Reason    string `json:"reason" validate:"required"`
	Detail    string `json:"detail" validate:"required"`
}
