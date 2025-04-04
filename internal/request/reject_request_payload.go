package request

type RejectRequestPayload struct {
	BookingId string `json:"bookingId" validate:"required"`
	Reason    string `json:"reason" validate:"reason"`
}
