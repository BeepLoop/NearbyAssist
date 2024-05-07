package utils

import "testing"

func TestStringSliceContains(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if !StringSliceContains(slice, "a") {
		t.Error("Expected slice to contain 'a'")
	}

	if !StringSliceContains(slice, "b") {
		t.Error("Expected slice to contain 'b'")
	}

	if !StringSliceContains(slice, "c") {
		t.Error("Expected slice to contain 'c'")
	}

	if StringSliceContains(slice, "d") {
		t.Error("Expected slice to not contain 'd'")
	}
}
