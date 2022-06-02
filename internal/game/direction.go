package game

type Direction uint8

const (
	None  Direction = 0x0
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

	return None
}

func DirectionAsStrings(dir Direction) []string {
	dirs := make([]string, 2)
	if dir&West != 0 {
		dirs = append(dirs, "west")
	}
	if dir&East != 0 {
		dirs = append(dirs, "east")
	}
	if dir&North != 0 {
		dirs = append(dirs, "north")
	}
	if dir&South != 0 {
		dirs = append(dirs, "south")
	}

	return dirs
}
