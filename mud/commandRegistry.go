package mud

type ServerAction func(server *Server) error
type CommandAction func(command Command, clientId ClientId) ServerAction

type CommandRegistry struct {
	commandActions map[string]CommandAction
	parser         CommandParser
}

func NewLoginCommandRegistry(parser CommandParser) *CommandRegistry {
	registry := &CommandRegistry{
		commandActions: map[string]CommandAction{
			"choose": NameCharacterCommandAction,
		},
		parser: parser,
	}

	return registry
}

func NewInGameCommandRegistry(parser CommandParser) *CommandRegistry {
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

func (c *CommandRegistry) ConnectAction(clientId ClientId) ServerAction {
	return ConnectCommandAction(Command{command: "connect", contents: ""}, clientId)
}

func (c *CommandRegistry) DisconnectAction(clientId ClientId) ServerAction {
	return DisconnectCommandAction(Command{command: "disconnect", contents: ""}, clientId)
}

func (c *CommandRegistry) InputToAction(line string, clientId ClientId) ServerAction {
	command := c.parser(line)

	commandAction, ok := c.commandActions[command.command]
	if ok {
		return commandAction(command, clientId)
	} else {
		return UnknownCommandAction(command, clientId)
	}
}
