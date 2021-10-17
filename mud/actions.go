package mud

import (
	"fmt"

	"github.com/google/uuid"
)

func ConnectCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		client := s.clients[clientId]
		s.world.InsertCharacterOnConnect(client.Character)
		client.reply <- "You joined\n"

		others := s.world.OtherCharactersInRoom(client.Character)
		for _, ch := range others {
			c := s.clients[ch.id]
			c.broadcast <- fmt.Sprintf("%v joined!\n", c.Character.name)
		}
		return nil
	}
}

func DisconnectCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		client := s.clients[clientId]
		client.reply <- "You are disconnecting"
		s.world.RemoveCharacterOnDisconnect(*client.Character)

		client.Disconnect()
		s.removeClientAtIndex(client.id)

		others := s.world.OtherCharactersInRoom(client.Character)
		for _, ch := range others {
			c := s.clients[ch.id]
			c.broadcast <- fmt.Sprintf("%v disconnecting...\n", c.Character.name)
		}
		return nil
	}
}

func UnknownCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		client := s.clients[clientId]
		client.reply <- fmt.Sprintf("What is %s?\n", command.contents)

		return nil
	}
}

func SayCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		client := s.clients[clientId]
		client.reply <- fmt.Sprintf("You said %s\n", command.contents)

		others := s.world.OtherCharactersInRoom(client.Character)
		for _, ch := range others {
			c := s.clients[ch.id]
			c.broadcast <- fmt.Sprintf("%s said %s\n", c.Character.name, command.contents)
		}

		return nil
	}
}

func GoCommandAction(command Command, clientId uuid.UUID) ServerAction {
	return func(s *Server) error {
		client := s.clients[clientId]

		if s.world.CanCharactorMoveInDirection(client.Character, command.contents) {
			others := s.world.OtherCharactersInRoom(client.Character)
			for _, ch := range others {
				c := s.clients[ch.id]
				c.broadcast <- fmt.Sprintf("%s moved to %s\n", c.Character.name, command.contents)
			}

			s.world.MoveCharacterInDirection(client.Character, command.contents)
			client.reply <- fmt.Sprintf("You move to %s\n", command.contents)

			others = s.world.OtherCharactersInRoom(client.Character)
			for _, ch := range others {
				c := s.clients[ch.id]
				c.broadcast <- fmt.Sprintf("%s entered from %s\n", c.Character.name, command.contents)
			}
		}

		return nil
	}
}
