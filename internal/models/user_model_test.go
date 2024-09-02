package models

import (
	"nearbyassist/internal/hash"
	"testing"
)

func TestUserEmailHash(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "user@email.com",
			expected: "0925f997eb0d742678f66d2da134d15d842d57722af5f7605c4785cb5358831b",
		},
	}

	h := hash.NewSha()
	for _, test := range tests {
		user := NewUserModelWithId("", nil)

		if _, err := user.HashEmail(test.input, h.Hash); err != nil {
			t.Errorf("Error hashing email: %v", err)
		}

		if user.Hash != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, user.Hash)
		}
	}
}
