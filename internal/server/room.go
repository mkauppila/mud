package server

type Room struct {
	description string
	location    Location
}

func NewRoom(description string, location Location) Room {
	return Room{
		description: description,
		location:    location,
	}
}

func BasicMap() []Room {
	return []Room{
		NewRoom("This is the room", Location{X: 0, Y: 0}),
		NewRoom("This another room", Location{X: 1, Y: 0}),
	}
}
