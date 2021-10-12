package mud

import "fmt"

func ConnectCommandAction(command Command, clientId int) ServerAction {
	return func(s *Server) error {
		for _, c := range s.clients {
			if c.id == clientId {
				c.reply <- "You joined\n"
			} else {
				c.broadcast <- fmt.Sprintf("Client %d joined!\n", clientId)
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
				// update the clients list
				s.clients[i] = s.clients[len(s.clients)-1]
				s.clients = s.clients[:len(s.clients)-1]
			} else {
				c.broadcast <- fmt.Sprintf("Client %d disconnecting...\n", clientId)
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
				c.broadcast <- fmt.Sprintf("They said %s\n", command.contents)
			}
		}
		return nil
	}
}

func GoCommandAction(command Command, clientId int) ServerAction {
	return func(s *Server) error {
		for _, c := range s.clients {
			if c.id == clientId {
				c.reply <- fmt.Sprintf("You move to %s\n", command.contents)
			}
		}
		return nil
	}
}
