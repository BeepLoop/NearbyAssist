package utils

import (
	"nearbyassist/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasScheduleOverlap(t *testing.T) {
	tests := []struct {
		input     string
		schedules []*models.TransactionModel
		expected  bool
	}{
		{
			input: "2025-04-3",
			schedules: []*models.TransactionModel{
				{ScheduledAt: "2025-04-2"},
				{ScheduledAt: "2025-04-3"},
				{ScheduledAt: "2025-04-4"},
				{ScheduledAt: "2025-04-5"},
				{ScheduledAt: "2025-04-6"},
			},
			expected: true,
		},
		{
			input: "2025-04-3",
			schedules: []*models.TransactionModel{
				{ScheduledAt: "2025-04-2"},
				{ScheduledAt: "2025-04-4"},
				{ScheduledAt: "2025-04-6"},
				{ScheduledAt: "2025-04-8"},
			},
			expected: false,
		},
	}

	for _, test := range tests {
		hasOverlap := HasScheduleOverlap(test.input, test.schedules)

		assert.Equal(t, test.expected, hasOverlap)
	}
}
