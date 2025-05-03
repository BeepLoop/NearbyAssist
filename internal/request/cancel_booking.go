package request

type CancelBookingPayload struct {
	BookingId string `json:"bookingId" validate:"required"`
	Reason    string `json:"reason" validate:"required"`
}
