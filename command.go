package main

import (
	"strings"
)

type Command struct {
	command  string // type disconnect, say
	contents string
	clientId int
}

func ParseCommand(message string) Command {
	message = strings.TrimSpace(message)

	index := strings.IndexAny(message, " ")
	command := message[:index]

	switch command {
	case "say":
		return Command{"say", message[index+1:], 0}
	case "go":
		// TODO: should parse the next part is north, west, east, south
		return Command{"go", message[index+1:], 0}
	case "connect":
		return Command{command: "connect", contents: "", clientId: 0}
	case "disconnect":
		return Command{command: "disconnect", contents: message, clientId: 0}

	}

	return Command{"unknown", message, 0}
}