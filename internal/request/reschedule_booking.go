package request

type RescheduleBookingPayload struct {
	BookingId     string `json:"bookingId" validte:"required"`
	ScheduleStart string `json:"scheduleStart" validate:"required"`
	ScheduleEnd   string `json:"scheduleEnd" validate:"required"`
}
