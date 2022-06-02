package game

import "fmt"

type Coordinate struct {
	X, Y int
}

func NewCoordinate(x, y int) Coordinate {
	return Coordinate{X: x, Y: y}
}

func (l Coordinate) String() string {
	return fmt.Sprintf("(%d, %d)", l.X, l.Y)
}

func CoordinateInDirection(location Coordinate, direction Direction) Coordinate {
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
