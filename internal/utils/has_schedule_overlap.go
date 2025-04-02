package utils

import "nearbyassist/internal/models"

func HasScheduleOverlap(schedule string, transactions []*models.TransactionModel) bool {
	sched := FormatDate(schedule)

	for _, transaction := range transactions {
		if sched == FormatDate(transaction.ScheduledAt) {
			return true
		}
	}

	return false
}
