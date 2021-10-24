package mud

import (
	"fmt"

	"github.com/google/uuid"
)

type ErrUnknownClientId struct {
	id uuid.UUID
}

func (e ErrUnknownClientId) Error() string {
	return fmt.Sprintf("unknown client id for %s", e.id)
}

func ConnectCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		client := s.getClient(clientId)
		if client == nil {
			return ErrUnknownClientId{id: clientId}
		}
		client.reply <- "Welcome! What is your name?\n"

		return nil
	}
}

func DisconnectCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		client := s.getClient(clientId)
		if client == nil {
			return ErrUnknownClientId{id: clientId}
		}

		client.reply <- "You are disconnecting"
		s.world.RemoveCharacterOnDisconnect(*client.Character)

		client.Disconnect()
		s.removeClientAtIndex(client.id)

		s.world.BroadcastToOtherCharactersInRoom(
			client.Character,
			fmt.Sprintf("%v disconnecting...\n", client.Character.name),
		)

		return nil
	}
}

func NameCharacterCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		client := s.getClient(clientId)
		if client == nil {
			return ErrUnknownClientId{id: clientId}
		}
		client.Character = NewCharacter(clientId, command.contents)

		ch := client.Character
		// connect character with clients comms
		client.Character.Reply = func(message string) {
			client.reply <- message
		}
		client.Character.Broadcast = func(message string) {
			client.broadcast <- message
		}
		client.SetCommandRegistry(NewInGameCommandRegistry(ParseInGameCommand))

		s.world.InsertCharacterOnConnect(ch)

		ch.Reply(fmt.Sprintf("%s woke up the world\n", ch.name))
		ch.SetState("idle")

		s.world.BroadcastToOtherCharactersInRoom(
			ch,
			fmt.Sprintf("%v joined!\n", ch.name),
		)

		return nil
	}
}

func UnknownCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		ch := s.world.getCharacter(clientId)
		if ch == nil {
			return ErrUnknownClientId{id: clientId}
		}

		ch.Reply(fmt.Sprintf("What is %s?\n", command.contents))

		return nil
	}
}

func SayCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		ch := s.world.getCharacter(clientId)
		if ch == nil {
			return ErrUnknownClientId{id: clientId}
		}
		ch.Reply(fmt.Sprintf("You said %s\n", command.contents))

		s.world.BroadcastToOtherCharactersInRoom(
			ch,
			fmt.Sprintf("%s said %s\n", ch.name, command.contents),
		)

		return nil
	}
}

func GoCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		ch := s.world.getCharacter(clientId)
		if ch == nil {
			return ErrUnknownClientId{id: clientId}
		}

		if s.world.CanCharactorMoveInDirection(ch, command.contents) {
			s.world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s moved to %s\n", ch.name, command.contents),
			)

			s.world.MoveCharacterInDirection(ch, command.contents)
			ch.Reply(fmt.Sprintf("You move to %s\n", command.contents))

			s.world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s moved to %s\n", ch.name, command.contents),
			)
		}

		return nil
	}
}

func StartSmokingCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		world := s.world

		ch := world.getCharacter(clientId)
		if ch == nil {
			return ErrUnknownClientId{id: clientId}
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
