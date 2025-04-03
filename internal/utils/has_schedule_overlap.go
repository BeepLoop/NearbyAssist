package utils

import (
	"fmt"
	"nearbyassist/internal/models"
)

func HasScheduleOverlap(schedule string, transactions []*models.TransactionModel) bool {
	sched := FormatDate(schedule)

	for _, transaction := range transactions {
		fmt.Println(transaction.ScheduledAt)
		if !transaction.ScheduledAt.Valid {
			continue
		}

		if sched == FormatDate(transaction.ScheduledAt.String) {
			return true
		}
	}

	return false
}
