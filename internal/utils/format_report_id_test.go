package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatReportID(t *testing.T) {
	t.Run("test format report ID", func(t *testing.T) {
		tests := []struct {
			input    int
			expected string
		}{
			{input: 1, expected: "REP-00001"},
			{input: 12, expected: "REP-00012"},
			{input: 12345, expected: "REP-12345"},
			{input: 123456, expected: "REP-123456"},
		}

		for _, test := range tests {
			formatted := FormatReportID(test.input)

			assert.Equal(t, test.expected, formatted)
		}
	})

	t.Run("test report id parser", func(t *testing.T) {
		tests := []struct {
			input     string
			expected  int
			shouldErr bool
		}{
			{input: "REP-00001", expected: 1, shouldErr: false},
			{input: "REP-01234", expected: 1234, shouldErr: false},
			{input: "01234", expected: 0, shouldErr: true},
			{input: "REP-123456", expected: 123456, shouldErr: false},
		}

		for _, test := range tests {
			parsed, err := ParseReportID(test.input)

			if test.shouldErr {
				assert.Error(t, err)
			}

			if !test.shouldErr {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expected, parsed)
		}
	})
}
