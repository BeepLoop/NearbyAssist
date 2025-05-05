package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportIDFormatter(t *testing.T) {
	t.Run("test report id formatter", func(t *testing.T) {
		tests := []struct {
			input    int
			expected string
		}{
			{input: 1, expected: "RP-00001"},
			{input: 12, expected: "RP-00012"},
			{input: 12345, expected: "RP-12345"},
			{input: 123456, expected: "RP-123456"},
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
			{input: "RP-00001", expected: 1, shouldErr: false},
			{input: "RP-01234", expected: 1234, shouldErr: false},
			{input: "01234", expected: 0, shouldErr: true},
			{input: "RP-123456", expected: 123456, shouldErr: false},
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
