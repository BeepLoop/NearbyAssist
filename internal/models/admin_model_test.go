package models

import (
	"nearbyassist/internal/hash"
	"testing"
)

func TestAdminHashUsername(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "admin",
			expected: "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918",
		},
		{
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			input:    "hello world",
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
	}

	h := hash.NewSha()
	for _, test := range tests {
		admin := NewAdminModelWithId("", nil)

		if _, err := admin.HashUsername(test.input, h.Hash); err != nil {
			t.Errorf("Error hashing username: %v", err)
		}

		if admin.UsernameHash != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, admin.UsernameHash)
		}
	}
}
