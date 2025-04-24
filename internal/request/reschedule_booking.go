package request

type RescheduleBookingPayload struct {
	BookingId string `json:"bookingId" validte:"required"`
	Schedule  string `json:"schedule" validate:"required"`
}
