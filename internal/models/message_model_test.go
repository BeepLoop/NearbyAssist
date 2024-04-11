package models

import (
	"reflect"
	"testing"
)

type test struct {
	query    string
	expected *MessageModel
}

func TestMessageModelFactory(t *testing.T) {

	tests := []test{
		{
			query: "sender=1&receiver=2",
			expected: &MessageModel{
				Sender:   1,
				Receiver: 2,
			},
		},
		{
			query: "sender=2&receiver=2",
			expected: &MessageModel{
				Sender:   2,
				Receiver: 2,
			},
		},
		{
			query:    "Sender=2&Receiver=2",
			expected: nil,
		},
		{
			query:    "sender=f&to=2lk",
			expected: nil,
		},
		{
			query:    "sender=12to=2",
			expected: nil,
		},
		{
			query:    "",
			expected: nil,
		},
		{
			query:    "sender=1",
			expected: nil,
		},
		{
			query:    "receiver=2",
			expected: nil,
		},
	}

	for _, test := range tests {
		values, _ := MessageValueMapFactory(test.query)

		if values != nil && test.expected == nil {
			t.Logf("Expected: %v \n", test.expected)
			t.Logf("Got: %v \n", values)
			t.Logf("Test case: %s\n", test.query)
			t.FailNow()
		}

		if reflect.DeepEqual(values, test.expected) == false {
			t.Logf("Expected: %v \n", test.expected)
			t.Logf("Got: %v \n", values)
			t.Logf("Test case: %s\n", test.query)
			t.FailNow()
		}
	}

}
