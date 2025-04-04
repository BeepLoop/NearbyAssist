package request

type AcceptBookingPayload struct {
	BookingId string `json:"bookingId" validate:"required"`
	Schedule  string `json:"schedule" validate:"required"`
}
