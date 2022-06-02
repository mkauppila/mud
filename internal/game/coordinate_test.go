package game

import "testing"

func TestNewLocationInDirection(t *testing.T) {
	testCases := []struct {
		loc    Coordinate
		dir    Direction
		wanted Coordinate
	}{
		{loc: NewCoordinate(0, 0), dir: West, wanted: NewCoordinate(-1, 0)},
		{loc: NewCoordinate(0, 0), dir: East, wanted: NewCoordinate(1, 0)},
		{loc: NewCoordinate(0, 0), dir: North, wanted: NewCoordinate(0, -1)},
		{loc: NewCoordinate(0, 0), dir: South, wanted: NewCoordinate(0, 1)},
	}

	for i, tc := range testCases {
		v := CoordinateInDirection(tc.loc, tc.dir)
		if v != tc.wanted {
			t.Fatalf("Testcase %d: Got %s, expected %s", i, v, tc.wanted)
		}
	}
}
