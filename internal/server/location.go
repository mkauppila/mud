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

func NewLocationInDirection(location Location, direction Direction) Location {
	switch direction {
	case West:
		location.X--
	case East:
		location.X++
	case North:
		location.Y--
	case South:
		location.Y++
	}

	return location
}
