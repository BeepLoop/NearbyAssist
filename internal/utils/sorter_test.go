package utils

import (
	"nearbyassist/internal/response"
	"reflect"
	"testing"
)

type testPair struct {
	name     string
	input    []response.SearchResult
	expected []response.SearchResult
}

func TestBubbleSort(t *testing.T) {
	tests := []testPair{
		{
			name: "test 1",
			input: []response.SearchResult{
				{Suggestability: 0.5},
				{Suggestability: 0.4},
				{Suggestability: 0.3},
				{Suggestability: 0.2},
				{Suggestability: 0.1},
			},
			expected: []response.SearchResult{
				{Suggestability: 0.1},
				{Suggestability: 0.2},
				{Suggestability: 0.3},
				{Suggestability: 0.4},
				{Suggestability: 0.5},
			},
		},
		{
			name: "test 2",
			input: []response.SearchResult{
				{Suggestability: 0.02},
				{Suggestability: 0.001},
				{Suggestability: 0.03},
			},
			expected: []response.SearchResult{
				{Suggestability: 0.001},
				{Suggestability: 0.02},
				{Suggestability: 0.03},
			},
		},
	}

	for _, test := range tests {
		result := BubbleSort(test.input)

		if reflect.DeepEqual(result, test.expected) == false {
			t.Errorf("Expected: %v \n, Got: %v\n", test.expected, result)
		}
	}
}
