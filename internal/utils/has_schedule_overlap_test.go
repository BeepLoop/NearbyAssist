package utils

import (
	"database/sql"
	"nearbyassist/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasScheduleOverlap(t *testing.T) {
	tests := []struct {
		input     string
		schedules []*models.BookingModel
		expected  bool
	}{
		{
			input: "2025-04-3",
			schedules: []*models.BookingModel{
				{ScheduledAt: sql.NullString{String: "2025-04-2", Valid: true}},
				{ScheduledAt: sql.NullString{String: "2025-04-3", Valid: true}},
				{ScheduledAt: sql.NullString{String: "2025-04-4", Valid: true}},
				{ScheduledAt: sql.NullString{String: "2025-04-5", Valid: true}},
				{ScheduledAt: sql.NullString{String: "2025-04-6", Valid: true}},
			},
			expected: true,
		},
		{
			input: "2025-04-3",
			schedules: []*models.BookingModel{
				{ScheduledAt: sql.NullString{String: "2025-04-2", Valid: true}},
				{ScheduledAt: sql.NullString{String: "2025-04-4", Valid: true}},
				{ScheduledAt: sql.NullString{String: "2025-04-5", Valid: true}},
				{ScheduledAt: sql.NullString{String: "2025-04-6", Valid: true}},
			},
			expected: false,
		},
	}

	for _, test := range tests {
		hasOverlap := HasScheduleOverlap(test.input, test.schedules)

		assert.Equal(t, test.expected, hasOverlap)
	}
}
