package utils

import (
	"fmt"
	"nearbyassist/internal/models"
)

func HasScheduleOverlap(schedule string, bookings []*models.BookingModel) bool {
	sched := FormatDate(schedule)

	for _, booking := range bookings {
		fmt.Println(booking.ScheduledAt)
		if !booking.ScheduledAt.Valid {
			continue
		}

		if sched == FormatDate(booking.ScheduledAt.String) {
			return true
		}
	}

	return false
}
