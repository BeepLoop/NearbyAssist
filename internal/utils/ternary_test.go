package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTernary(t *testing.T) {
	t.Run("test for strings", func(t *testing.T) {
		tests := []struct {
			inputCondition bool
			inputGood      string
			inputFallback  string
			expected       string
		}{
			{
				inputCondition: true,
				inputGood:      "good",
				inputFallback:  "fallback",
				expected:       "good",
			},
			{
				inputCondition: false,
				inputGood:      "good",
				inputFallback:  "fallback",
				expected:       "fallback",
			},
		}

		for _, test := range tests {
			res := Ternary(test.inputCondition, test.inputGood, test.inputFallback)
			assert.EqualValues(t, test.expected, res)
		}
	})

	t.Run("test for objects", func(t *testing.T) {
		type Foo struct {
			bar  string
			baz  int
			bass []int
		}

		tests := []struct {
			inputCondition bool
			inputGood      Foo
			inputFallback  Foo
			expected       Foo
		}{
			{
				inputCondition: 1 == 1,
				inputGood:      Foo{bar: "bar", baz: 1, bass: []int{1, 2, 3}},
				inputFallback:  Foo{},
				expected:       Foo{bar: "bar", baz: 1, bass: []int{1, 2, 3}},
			},
			{
				inputCondition: 1 == 2,
				inputGood:      Foo{bar: "bar", baz: 1, bass: []int{1, 2, 3}},
				inputFallback:  Foo{},
				expected:       Foo{},
			},
		}

		for _, test := range tests {
			res := Ternary(test.inputCondition, test.inputGood, test.inputFallback)
			assert.EqualValues(t, test.expected, res)
		}
	})
}
