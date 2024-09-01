package utils

import (
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{
			input: "hello=world&value=2",
			expected: map[string]string{
				"hello": "world",
				"value": "2",
			},
		},
		{
			input: "l=1234,567",
			expected: map[string]string{
				"l": "1234,567",
			},
		},
		{
			input:    "",
			expected: map[string]string{},
		},
	}

	for _, test := range tests {
		query := ParseQuery(test.input)

		if reflect.DeepEqual(query, test.expected) == false {
			t.Errorf("ParseQuery(%s) = %v, want %v", test.input, query, test.expected)
		}
	}
}
