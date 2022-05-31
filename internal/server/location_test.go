package server

import "testing"

func TestNewLocationInDirection(t *testing.T) {
	testCases := []struct {
		loc    Location
		dir    Direction
		wanted Location
	}{
		{loc: NewLocation(0, 0), dir: West, wanted: NewLocation(-1, 0)},
		{loc: NewLocation(0, 0), dir: East, wanted: NewLocation(1, 0)},
		{loc: NewLocation(0, 0), dir: North, wanted: NewLocation(0, -1)},
		{loc: NewLocation(0, 0), dir: South, wanted: NewLocation(0, 1)},
	}

	for i, tc := range testCases {
		v := NewLocationInDirection(tc.loc, tc.dir)
		if v != tc.wanted {
			t.Fatalf("Testcase %d: Got %s, expected %s", i, v, tc.wanted)
		}
	}
}
