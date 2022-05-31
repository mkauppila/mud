package server

type Direction uint8

const (
	Error Direction = 0
	North           = 0x1
	East            = 0x2
	South           = 0x4
	West            = 0x8
)

func DirectionFromString(dir string) Direction {
	switch dir {
	case "west":
		return West
	case "east":
		return East
	case "north":
		return North
	case "south":
		return South
	}

	return Error
}
