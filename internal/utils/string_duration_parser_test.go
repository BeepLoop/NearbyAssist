package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseStringDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{
			input:    "10s",
			expected: time.Second * 10,
		},
		{
			input:    "05s",
			expected: time.Second * 5,
		},
		{
			input:    "1m",
			expected: time.Minute * 1,
		},
		{
			input:    "2h",
			expected: time.Hour * 2,
		},
		{
			input:    "02h",
			expected: time.Hour * 2,
		},
		{
			input:    "10d",
			expected: (time.Hour * 24) * 10,
		},
		{
			input:    "02d",
			expected: (time.Hour * 24) * 2,
		},
	}

	for _, test := range tests {
		duration, err := ParseStringDuration(test.input)

		assert.Nil(t, err)
		assert.Equal(t, test.expected, duration)
	}
}
