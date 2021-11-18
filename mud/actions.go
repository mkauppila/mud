package mud

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
	description string
	action      CommandAction
}

var loginCommandInfos = []CommandInfo{
	{
		command:     "choose",
		description: "Choose a character name",
		action:      NameCharacterCommandAction,
	},
}

var ingameCommandInfos = []CommandInfo{
	{
		command:     "help",
		description: "List all the commands and stuff",
		action:      HelpCommandAction,
	},
	{
		command:     "say",
		description: "Say something",
		action:      SayCommandAction,
	},
	{
		command:     "go",
		description: "Move to east, west, north or south",
		action:      GoCommandAction,
	},
	{
		command:     "smoke",
		description: "You can _start_ or _stop_ smoking",
		action:      SmokeCommandAction,
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
		client.SetCommandRegistry(NewInGameCommandRegistry(ParseInGameCommand))

		s.world.InsertCharacterOnConnect(ch)

		ch.Reply(
			fmt.Sprintf("%s woke up the world\n%s\n",
				ch.name,
				s.world.DescribeRoom(ch.Location)),
		)
		ch.SetState("idle")

		s.world.BroadcastToOtherCharactersInRoom(
			ch,
			fmt.Sprintf("%v joined!\n", ch.name),
		)

		return nil
	}
}

func UnknownCommandAction(command Command, clientId ClientId) ServerAction {
	return func(s *Server) error {
		ch := s.world.getCharacter(clientId)
		if ch == nil {
			return ErrUnknownCharacter{id: clientId, action: command.command}
		}

		ch.Reply(fmt.Sprintf("What is %s?\n", command.contents))

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

func GoCommandAction(command Command, clientId ClientId) ServerAction {
	return func(s *Server) error {
		ch := s.world.getCharacter(clientId)
		if ch == nil {
			return ErrUnknownCharacter{id: clientId, action: command.command}
		}

		if s.world.CanCharactorMoveInDirection(ch, command.contents) {
			// broadcast to old room
			s.world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s moved to %s\n", ch.name, command.contents),
			)

			s.world.MoveCharacterInDirection(ch, command.contents)
			ch.Reply(
				fmt.Sprintf("You move to %s\n%s\n",
					command.contents,
					s.world.DescribeRoom(ch.Location)),
			)

			// broadcast to new room
			s.world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s entered from %s\n", ch.name, command.contents),
			)
		} else {
			ch.Reply("Ouch, it seems the world has some boundaries\n")
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
		c := client.registry.CommandsWithDescriptions()

		var output = "help:\n"
		for _, cc := range c {
			output = fmt.Sprintf("%s\t%s\n", output, cc)
		}

		client.reply <- output

		return nil
	}
}
