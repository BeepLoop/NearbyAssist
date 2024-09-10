package hash

// test for sha256.go
import (
	"testing"
)

type test struct {
	input    string
	expected string
}

func TestNewSha(t *testing.T) {
	tests := []test{
		{
			input:    "test",
			expected: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		},
		{
			input:    "test 2",
			expected: "dec2e4bc4992314a9c9a51bbd859e1b081b74178818c53c19d18d6f761f5d804",
		},
		{
			input:    "foo@email.com",
			expected: "86b4ae9987424467ed614e59a98a1c4952e7f69cc619a9e49761c9ea37953fbb",
		},
	}

	sha := NewSha()
	for i, test := range tests {
		i = i + 1
		result, err := sha.Hash([]byte(test.input))
		if err != nil {
			t.Errorf("Error testing input: %d, %v", i, err)
		}

		if result != test.expected {
			t.Logf("test input: %s", test.input)
			t.Errorf("Test %d failed: \nExpected %s\n but got %s\n", i, test.expected, result)
		}
	}
}
