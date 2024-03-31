package utils

import (
	"errors"
	"testing"
)

type test struct {
	err    error
	expect bool
}

func TestDetermineNoRowsError(t *testing.T) {

	tests := []test{
		{err: errors.New("no rows in result set"), expect: true},
		{err: errors.New("sql: no rows in result set"), expect: true},
		{err: errors.New("some other error"), expect: false},
		{err: errors.New("rows in result"), expect: false},
	}

	for _, test := range tests {
		result := DetermineNoRowsError(test.err)
		if result != test.expect {
			t.Errorf("Expected %v, got %v", test.expect, result)
		}
	}
}
