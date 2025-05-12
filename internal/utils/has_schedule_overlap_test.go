package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasScheduleOverlap(t *testing.T) {
	tests := []struct {
		input     Schedule
		schedules []Schedule
		expected  bool
	}{
		{
			input: Schedule{Start: "2025-05-02", End: "2025-05-04"},
			schedules: []Schedule{
				{Start: "2025-05-05", End: "2025-05-06"},
				{Start: "2025-05-01", End: "2025-05-01"},
			},
			expected: false,
		},
		{
			input: Schedule{Start: "2025-05-02", End: "2025-05-02"},
			schedules: []Schedule{
				{Start: "2025-05-05", End: "2025-05-06"},
				{Start: "2025-05-01", End: "2025-05-01"},
			},
			expected: false,
		},
		{
			input: Schedule{Start: "2025-05-02", End: "2025-05-04"},
			schedules: []Schedule{
				{Start: "2025-05-04", End: "2025-05-06"},
				{Start: "2025-05-01", End: "2025-05-01"},
			},
			expected: true,
		},
		{
			input: Schedule{Start: "2025-05-02", End: "2025-05-04"},
			schedules: []Schedule{
				{Start: "2025-05-01", End: "2025-05-06"},
			},
			expected: true,
		},
		{
			input: Schedule{Start: "2025-05-06", End: "2025-05-07"},
			schedules: []Schedule{
				{Start: "2025-05-01", End: "2025-05-06"},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		hasOverlap, err := HasScheduleOverlap(test.input, test.schedules)

		assert.NoError(t, err)
		assert.Equal(t, test.expected, hasOverlap)
	}
}
