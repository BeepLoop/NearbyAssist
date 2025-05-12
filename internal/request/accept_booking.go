package request

type AcceptBookingPayload struct {
	BookingId     string `json:"bookingId" validate:"required"`
	ScheduleStart string `json:"scheduleStart" validate:"required"`
	ScheduleEnd   string `json:"scheduleEnd" validate:"required"`
}
