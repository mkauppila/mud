package server

type Room struct {
	description string
	location    Coordinate
	exits       Direction
}

func (r Room) HasExitInDirection(dir Direction) bool {
	return r.exits&dir != 0
}

func NewRoom(description string, location Coordinate) Room {
	return Room{
		description: description,
		location:    location,
	}
}

func BasicMap() []Room {
	return []Room{
		NewRoom("This is the room", Coordinate{X: 0, Y: 0}),
		NewRoom("This another room", Coordinate{X: 1, Y: 0}),
	}
}
