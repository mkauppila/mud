package mud

import "github.com/google/uuid"

type ServerAction func(server *Server) error
type CommandAction func(command Command, clientId uuid.UUID) ServerAction

type CommandRegistry struct {
	commandActions map[string]CommandAction
	parser         CommandParser
}

func NewLoginCommandRegistry(parser CommandParser) *CommandRegistry {
	registry := &CommandRegistry{
		commandActions: map[string]CommandAction{
			"choose": ChooseCharacterCommandAction,
		},
		parser: parser,
	}

	return registry
}

// Rename: New in-game registry
func NewCommandRegistry(parser CommandParser) *CommandRegistry {
	registry := &CommandRegistry{
		commandActions: map[string]CommandAction{
			"say":   SayCommandAction,
			"go":    GoCommandAction,
			"smoke": StartSmokingCommandAction,
		},
		parser: parser,
	}

	return registry
}

func (c *CommandRegistry) ConnectAction(clientId uuid.UUID) ServerAction {
	return ConnectCommandAction(Command{command: "connect", contents: ""}, clientId)
}

func (c *CommandRegistry) DisconnectAction(clientId uuid.UUID) ServerAction {
	return DisconnectCommandAction(Command{command: "disconnect", contents: ""}, clientId)
}

func (c *CommandRegistry) InputToAction(line string, clientId uuid.UUID) ServerAction {
	command := c.parser.ParseCommand(line)

	commandAction, ok := c.commandActions[command.command]
	if ok {
		return commandAction(command, clientId)
	} else {
		return UnknownCommandAction(command, clientId)
	}
}
