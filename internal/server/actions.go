package server

import (
	"fmt"
)

type ErrUnknownClientId struct {
	id ClientId
}

func (e ErrUnknownClientId) Error() string {
	return fmt.Sprintf("unknown client id for %s", e.id)
}

type ErrUnknownCharacter struct {
	id     ClientId
	action string
}

func (e ErrUnknownCharacter) Error() string {
	return fmt.Sprintf("unknown character. Client id for %s. Action: %s", e.id, e.action)
}

type CommandInfo struct {
	command     string
	aliases     []string
	description string
	parser      func(command, rest string) Command
	action      CommandAction
}

var loginCommandInfos = []CommandInfo{
	{
		command:     "choose",
		aliases:     []string{},
		description: "Choose a character name",
		parser: func(command, rest string) Command {
			return Command{"choose", rest}
		},
		action: NameCharacterCommandAction,
	},
}

var inGameCommandInfos = []CommandInfo{
	{
		command:     "help",
		aliases:     []string{},
		description: "List all the commands and stuff",
		parser: func(command, rest string) Command {
			return Command{"help", ""}
		},
		action: HelpCommandAction,
	},
	{
		command:     "say",
		aliases:     []string{},
		description: "Say something",
		parser: func(command, rest string) Command {
			return Command{"say", rest}
		},
		action: SayCommandAction,
	},
	{
		command:     "go",
		aliases:     []string{"n", "e", "s", "w"},
		description: "Move to east, west, north or south",
		parser: func(command, rest string) Command {
			message := rest
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
		},
		action: GoCommandAction,
	},
	{
		command:     "look",
		aliases:     []string{"l"},
		description: "Look around the room",
		parser: func(command, rest string) Command {
			return Command{"look", rest}
		},
		action: LookCommandAction,
	},
	{
		command:     "smoke",
		aliases:     []string{},
		description: "You can _start_ or _stop_ smoking",
		parser: func(command, rest string) Command {
			return Command{"smoke", rest}
		},
		action: SmokeCommandAction,
	},
}

func ConnectCommandAction(command Command, clientId ClientId) ServerAction {
	return func(s *Server) error {
		client := s.getClient(clientId)
		if client == nil {
			return ErrUnknownClientId{id: clientId}
		}
		client.reply <- "Welcome! What is your name?\n"

		return nil
	}
}

func DisconnectCommandAction(command Command, clientId ClientId) ServerAction {
	return func(s *Server) error {
		client := s.getClient(clientId)
		if client == nil {
			return ErrUnknownClientId{id: clientId}
		}

		if ch := s.world.getCharacter(clientId); ch != nil {
			s.world.RemoveCharacterOnDisconnect(ch)
			s.world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%v disconnecting...\n", ch.name),
			)
		} else {
			return ErrUnknownCharacter{id: clientId, action: command.command}
		}

		client.Disconnect()
		s.removeClient(clientId)

		return nil
	}
}

func UnknownCommandAction(command Command, clientId ClientId) ServerAction {
	return func(s *Server) error {
		client := s.getClient(clientId)
		if client == nil {
			return ErrUnknownClientId{id: clientId}
		}

		client.reply <- fmt.Sprintf("What is %s?\n", command.contents)

		return nil
	}
}

func NameCharacterCommandAction(command Command, clientId ClientId) ServerAction {
	return func(s *Server) error {
		client := s.getClient(clientId)
		if client == nil {
			return ErrUnknownClientId{id: clientId}
		}
		ch := NewCharacter(clientId, command.contents)

		// connect character with clients comms
		ch.Reply = func(message string) {
			client.reply <- message
		}
		ch.Broadcast = func(message string) {
			client.broadcast <- message
		}
		client.SetCommandRegistry(NewInGameCommandRegistry())

		s.world.InsertCharacterOnConnect(ch)

		ch.Reply(
			fmt.Sprintf("%s woke up in the world\n%s\n",
				ch.name,
				s.world.DescribeRoom(ch.Coordinate)),
		)
		ch.SetState("idle")

		s.world.BroadcastToOtherCharactersInRoom(
			ch,
			fmt.Sprintf("%v joined!\n", ch.name),
		)

		return nil
	}
}

func SayCommandAction(command Command, clientId ClientId) ServerAction {
	return func(s *Server) error {
		ch := s.world.getCharacter(clientId)
		if ch == nil {
			return ErrUnknownCharacter{id: clientId, action: command.command}
		}
		ch.Reply(fmt.Sprintf("You said %s\n", command.contents))

		s.world.BroadcastToOtherCharactersInRoom(
			ch,
			fmt.Sprintf("%s said %s\n", ch.name, command.contents),
		)

		return nil
	}
}

func LookCommandAction(command Command, clientId ClientId) ServerAction {
	return func(s *Server) error {
		ch := s.world.getCharacter(clientId)
		if ch == nil {
			return ErrUnknownCharacter{id: clientId, action: command.command}
		}

		ch.Reply(fmt.Sprintf("You look around\n%s\n", s.world.DescribeRoom(ch.Coordinate)))

		return nil
	}
}

func GoCommandAction(command Command, clientId ClientId) ServerAction {
	return func(s *Server) error {
		ch := s.world.getCharacter(clientId)
		if ch == nil {
			return ErrUnknownCharacter{id: clientId, action: command.command}
		}

		direction := DirectionFromString(command.contents)
		if command.contents == "" {
			ch.Reply("In which direction do you want to move?\n")
		} else if s.world.CanCharactorMoveInDirection(ch, direction) {
			// broadcast to old room
			s.world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s moved to %s\n", ch.name, command.contents),
			)

			s.world.MoveCharacterInDirection(ch, direction)
			ch.Reply(
				fmt.Sprintf("You move to %s\n%s\n",
					command.contents,
					s.world.DescribeRoom(ch.Coordinate)),
			)

			// broadcast to new room
			s.world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s entered from %s\n", ch.name, command.contents),
			)
		} else {
			ch.Reply("You cannot go that way!\n")
		}

		return nil
	}
}

func SmokeCommandAction(command Command, clientId ClientId) ServerAction {
	return func(s *Server) error {
		world := s.world

		ch := world.getCharacter(clientId)
		if ch == nil {
			return ErrUnknownCharacter{id: clientId, action: command.command}
		}

		switch command.contents {
		case "start":
			ch.SetState("smoking")
			ch.Reply("You started to smoke your pipe\n")

			world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s started to smoke a pipe\n", ch.name),
			)
		case "stop":
			ch.SetState("idle")
			ch.Reply("You stopped smoking your pipe\n")

			world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s stopped smoking a pipe\n", ch.name),
			)
		default:
			ch.Reply(fmt.Sprintln("You either start or stop"))
		}

		return nil
	}
}

func HelpCommandAction(command Command, clientId ClientId) ServerAction {
	return func(s *Server) error {
		client, ok := s.clients[clientId]
		if !ok {
			return nil
		}

		var output = "help:\n"
		for _, cwd := range client.registry.CommandsWithDescriptions() {
			output = fmt.Sprintf("%s\t%s\n", output, cwd)
		}

		client.reply <- output

		return nil
	}
}
