package mud

import "fmt"

type ServerAction func(server *Server) error
type CommandAction func(command Command, clientId ClientId) ServerAction

type CommandRegistry struct {
	commandInfos map[string]CommandInfo
	parser       CommandParser
}

func NewLoginCommandRegistry(parser CommandParser) *CommandRegistry {
	// TODO: pass the loginCommandInfos as a parameter
	commandInfos := make(map[string]CommandInfo)
	for _, i := range loginCommandInfos {
		commandInfos[i.command] = i
	}

	registry := &CommandRegistry{
		commandInfos: commandInfos,
		parser:       parser,
	}

	return registry
}

func NewInGameCommandRegistry(parser CommandParser) *CommandRegistry {
	// define the shit in an array
	// create the action registry based on it
	// key: <the whole action info>
	// return just return the function after the parsing
	commandInfos := make(map[string]CommandInfo)
	for _, i := range ingameCommandInfos {
		commandInfos[i.command] = i
	}

	registry := &CommandRegistry{
		commandInfos: commandInfos,
		parser:       parser,
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

	info, ok := c.commandInfos[command.command]
	if ok {
		return info.action(command, clientId)
	} else {
		return UnknownCommandAction(command, clientId)
	}
}

func (c *CommandRegistry) CommandsWithDescriptions() []CommandWithDescription {
	var result []CommandWithDescription
	for _, v := range c.commandInfos {
		result = append(result, CommandWithDescription{
			command:     v.command,
			description: v.description,
		})
	}
	return result
}

type CommandWithDescription struct {
	command     string
	description string
}

func (c CommandWithDescription) String() string {
	return fmt.Sprintf("%s\t%s", c.command, c.description)
}
