package mud

import "github.com/google/uuid"

type ServerAction func(server *Server) error
type CommandAction func(command Command, clientId uuid.UUID) ServerAction

type CommandHandler struct {
	commandActions map[string]CommandAction
}

func NewCommandHandler() *CommandHandler {
	commander := &CommandHandler{
		commandActions: map[string]CommandAction{
			"say":   SayCommandAction,
			"go":    GoCommandAction,
			"smoke": StartSmokingCommandAction,
		},
	}

	return commander
}

func (c *CommandHandler) ConnectAction(clientId uuid.UUID) ServerAction {
	return ConnectCommandAction(Command{command: "connect", contents: ""}, clientId)
}

func (c *CommandHandler) DisconnectAction(clientId uuid.UUID) ServerAction {
	return DisconnectCommandAction(Command{command: "disconnect", contents: ""}, clientId)
}

func (c *CommandHandler) InputToAction(line string, clientId uuid.UUID) ServerAction {
	command := ParseCommand(line)

	commandAction, ok := c.commandActions[command.command]
	if ok {
		return commandAction(command, clientId)
	} else {
		return UnknownCommandAction(command, clientId)
	}
}
