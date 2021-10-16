package mud

import "fmt"

func ConnectCommandAction(command Command, clientId int) ServerAction {
	return func(s *Server) error {
		for _, c := range s.clients {
			if c.id == clientId {
				s.world.InsertCharacterOnJoin(c.Character)
				c.reply <- "You joined\n"
			} else {
				c.broadcast <- fmt.Sprintf("%v joined!\n", c.Character.name)
			}
		}
		return nil
	}
}

func DisconnectCommandAction(command Command, clientId int) ServerAction {
	return func(s *Server) error {
		for i, c := range s.clients {
			if c.id == clientId {
				c.reply <- "You are disconnecting"
				c.Disconnect()

				s.removeClientAtIndex(i)
			} else {
				c.broadcast <- fmt.Sprintf("%v disconnecting...\n", c.Character.name)
			}
		}
		return nil
	}
}

func UnknownCommandAction(command Command, clientId int) ServerAction {
	return func(s *Server) error {
		for _, c := range s.clients {
			if c.id == clientId {
				c.reply <- fmt.Sprintf("What is %s?\n", command.contents)
			}
		}
		return nil
	}
}

func SayCommandAction(command Command, clientId int) ServerAction {
	return func(s *Server) error {
		for _, c := range s.clients {
			if c.id == clientId {
				c.reply <- fmt.Sprintf("You said %s\n", command.contents)
			} else {
				c.broadcast <- fmt.Sprintf("%s said %s\n", c.Character.name, command.contents)
			}
		}
		return nil
	}
}

func GoCommandAction(command Command, clientId int) ServerAction {
	return func(s *Server) error {
		for _, c := range s.clients {
			if c.id == clientId {
				s.world.MoveCharacterInDirection(c.Character, command.contents)
				c.reply <- fmt.Sprintf("You move to %s\n", command.contents)
			} else {
				// TODO: if moving fails, is blocke etc., this will still be send to the other clients
				c.broadcast <- fmt.Sprintf("%s moved to %s\n", c.Character.name, command.contents)
			}
		}
		return nil
	}
}
