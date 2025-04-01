package searchhistory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchHistoryImplementation(t *testing.T) {
	t.Run("test get all", func(t *testing.T) {
		tests := []struct {
			input         string
			initialValues []string
			expected      []string
		}{
			{
				input:         "foo",
				initialValues: []string{"bar"},
				expected:      []string{"foo", "bar"},
			},
			{
				input:         "baz",
				initialValues: []string{"foo", "bar"},
				expected:      []string{"baz", "bar", "foo"},
			},
		}

		for _, test := range tests {
			hist := New()

			for _, input := range test.initialValues {
				hist.Insert(input)
			}

			hist.Insert(test.input)

			assert.EqualValues(t, test.expected, hist.GetAll())
			Instance = nil
		}
	})

	t.Run("test get items with count", func(t *testing.T) {
		tests := []struct {
			count         int
			initialValues []string
			expected      []string
		}{
			{
				count:         2,
				initialValues: []string{"foo", "bar", "baz", "bass"},
				expected:      []string{"bass", "baz"},
			},
			{
				count:         2,
				initialValues: []string{"foo"},
				expected:      []string{"foo"},
			},
		}

		for _, test := range tests {
			hist := New()

			for _, value := range test.initialValues {
				hist.Insert(value)
			}

			values := hist.GetCount(test.count)

			assert.EqualValues(t, test.expected, values)
			Instance = nil
		}
	})
}
