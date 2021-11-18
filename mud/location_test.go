package mud

import "testing"

func TestNewLocationInDirection(t *testing.T) {
	testCases := []struct {
		loc    Location
		dir    string
		wanted Location
	}{
		{loc: NewLocation(0, 0), dir: "west", wanted: NewLocation(-1, 0)},
		{loc: NewLocation(0, 0), dir: "east", wanted: NewLocation(1, 0)},
		{loc: NewLocation(0, 0), dir: "north", wanted: NewLocation(0, -1)},
		{loc: NewLocation(0, 0), dir: "south", wanted: NewLocation(0, 1)},
	}

	for i, tc := range testCases {
		v := NewLocationInDirection(tc.loc, tc.dir)
		if v != tc.wanted {
			t.Fatalf("Testcase %d: Got %s, expected %s", i, v, tc.wanted)
		}
	}
}
