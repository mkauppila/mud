package game

import (
	"fmt"
	"strings"
)

type Command struct {
	command  string // TODO: make it enum: type disconnect, say
	contents string
}

type CommandRegistry struct {
	commandInfos map[string]CommandInfo
}

func NewCommandRegistry(givenCommandInfo []CommandInfo) *CommandRegistry {
	commandInfos := make(map[string]CommandInfo)
	for _, i := range givenCommandInfo {
		commandInfos[i.command] = i
	}

	registry := &CommandRegistry{
		commandInfos: commandInfos,
	}

	return registry
}

func NewLoginCommandRegistry() *CommandRegistry {
	return NewCommandRegistry(loginCommandInfos)
}

func NewInGameCommandRegistry() *CommandRegistry {
	return NewCommandRegistry(inGameCommandInfos)
}

func (c *CommandRegistry) InputToAction(line string, clientId ClientId) WorldAction {
	command := c.parseCommand(line)

	info, ok := c.commandInfos[command.command]
	if ok {
		return info.action(command, clientId)
	} else {
		return UnknownCommandAction(command, clientId)
	}
}

func (c *CommandRegistry) parseCommand(message string) Command {
	message = strings.ToLower(strings.TrimSpace(message))

	index := strings.IndexAny(message, " ")
	var command, rest string
	if index >= 0 {
		command = message[:index]
		rest = message[index+1:]
	} else {
		command = message
		rest = ""
	}

	for _, v := range c.commandInfos {
		if command == v.command {
			return v.parser(command, rest)
		}

		for _, alias := range v.aliases {
			if command == alias {
				return v.parser(command, rest)
			}
		}
	}

	return Command{"unknown", message}
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
