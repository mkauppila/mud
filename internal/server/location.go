package server

import "fmt"

type Location struct {
	X, Y int
}

func NewLocation(x, y int) Location {
	return Location{X: x, Y: y}
}

func (l Location) String() string {
	return fmt.Sprintf("(%d, %d)", l.X, l.Y)
}

func NewLocationInDirection(location Location, direction string) Location {
	switch direction {
	case "west":
		location.X--
	case "east":
		location.X++
	case "north":
		location.Y--
	case "south":
		location.Y++
	}

	return location
}
