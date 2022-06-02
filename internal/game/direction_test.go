package game

import "testing"

func TestDirectionParsing(t *testing.T) {
	testCases := []struct {
		input    string
		expected Direction
	}{
		{input: "north", expected: North},
		{input: "east", expected: East},
		{input: "south", expected: South},
		{input: "west", expected: West},
		{input: "error", expected: None},
	}

	for _, tc := range testCases {
		if DirectionFromString(tc.input) != tc.expected {
			t.Fatalf("%d did not match %s", tc.expected, tc.input)
		}
	}
}
