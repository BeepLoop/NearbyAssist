package searchhistory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchHistoryImplementation(t *testing.T) {
	t.Run("insert with default capacity", func(t *testing.T) {
		tests := []struct {
			input    string
			expected []string
		}{
			{
				input:    "test",
				expected: []string{"test"},
			},
			{
				input:    "another",
				expected: []string{"another"},
			},
		}

		for _, test := range tests {
			hist := New()
			hist.Insert(test.input)

			assert.Equal(t, 1, hist.GetSize())
			assert.EqualValues(t, test.expected, hist.GetAll())

			Destroy()
		}
	})

	t.Run("insert when capacity is reached", func(t *testing.T) {
		tests := []struct {
			input         string
			capacity      int
			initialValues []string
			expected      []string
		}{
			{
				input:         "foo",
				capacity:      2,
				initialValues: []string{"bar"},
				expected:      []string{"foo", "bar"},
			},
			{
				input:         "baz",
				capacity:      2,
				initialValues: []string{"foo", "bar"},
				expected:      []string{"baz", "bar"},
			},
		}

		for _, test := range tests {
			hist := New()
			hist.capacity = test.capacity

			for _, input := range test.initialValues {
				hist.Insert(input)
			}

			hist.Insert(test.input)

			assert.EqualValues(t, test.expected, hist.GetAll())

			Destroy()
		}
	})

	t.Run("retrieve top k entries", func(t *testing.T) {
		tests := []struct {
			k             int
			initialValues []string
			expected      []string
		}{
			{
				k:             2,
				initialValues: []string{"foo", "bar", "baz", "bass"},
				expected:      []string{"bass", "baz"},
			},
			{
				k:             2,
				initialValues: []string{"foo"},
				expected:      []string{"foo"},
			},
		}

		for _, test := range tests {
			hist := New()

			for _, value := range test.initialValues {
				hist.Insert(value)
			}

			values := hist.GetTopKElements(test.k)

			assert.EqualValues(t, test.expected, values)

			Destroy()
		}
	})

	t.Run("skip duplicate on insert", func(t *testing.T) {
		tests := []struct {
			input         string
			initialValues []string
			expected      []string
		}{
			{
				input:         "foo",
				initialValues: []string{"foo", "bar", "baz"},
				expected:      []string{"baz", "bar", "foo"},
			},
		}

		for _, test := range tests {
			hist := New()

			for _, value := range test.initialValues {
				hist.Insert(value)
			}
			hist.Insert(test.input)

			assert.EqualValues(t, test.expected, hist.GetAll())

			Destroy()
		}
	})
}
