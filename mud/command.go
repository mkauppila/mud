package mud

import (
	"strings"
)

type Command struct {
	command  string // TODO: make it enum: type disconnect, say
	contents string
}

func ParseCommand(message string) Command {
	message = strings.TrimSpace(message)

	index := strings.IndexAny(message, " ")
	var command string
	if index >= 0 {
		command = message[:index]
	} else {
		command = message
	}

	switch command {
	case "say":
		return Command{"say", message[index+1:]}
	case "go", "n", "e", "s", "w":
		// TODO: should parse the next part is north, west, east, south
		message := message[index+1:]
		switch command {
		case "n":
			message = "north"
		case "e":
			message = "east"
		case "s":
			message = "south"
		case "w":
			message = "west"
		}
		return Command{"go", message}
	}

	return Command{"unknown", message}
}