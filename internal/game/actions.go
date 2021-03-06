package game

import (
	"fmt"
)

type WorldAction func(w *World) error
type CommandAction func(command Command, ch *Character) WorldAction

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
	aliases     []string
	description string
	parser      func(command, rest string) Command
	action      CommandAction
}

var inGameCommandInfos = []CommandInfo{
	{
		command:     "help",
		aliases:     []string{},
		description: "List all the commands and stuff",
		parser: func(command, rest string) Command {
			return Command{"help", ""}
		},
		action: HelpCommandAction,
	},
	{
		command:     "say",
		aliases:     []string{},
		description: "Say something",
		parser: func(command, rest string) Command {
			return Command{"say", rest}
		},
		action: SayCommandAction,
	},
	{
		command:     "go",
		aliases:     []string{"n", "e", "s", "w"},
		description: "Move to east, west, north or south",
		parser: func(command, rest string) Command {
			message := rest
			switch command {
			case "n":
				message = "north"
			case "e":
				message = "east"
			case "s":
				message = "south"
			case "w":
				message = "west"
			}
			return Command{"go", message}
		},
		action: GoCommandAction,
	},
	{
		command:     "look",
		aliases:     []string{"l"},
		description: "Look around the room",
		parser: func(command, rest string) Command {
			return Command{"look", rest}
		},
		action: LookCommandAction,
	},
	{
		command:     "smoke",
		aliases:     []string{},
		description: "You can _start_ or _stop_ smoking",
		parser: func(command, rest string) Command {
			return Command{"smoke", rest}
		},
		action: SmokeCommandAction,
	},
}

func UnknownCommandAction(command Command, ch *Character) WorldAction {
	return func(w *World) error {
		ch.Reply(fmt.Sprintf("What is %s?\n", command.contents))

		return nil
	}
}

func SayCommandAction(command Command, ch *Character) WorldAction {
	return func(world *World) error {
		ch.Reply(fmt.Sprintf("You said %s\n", command.contents))

		world.BroadcastToOtherCharactersInRoom(
			ch,
			fmt.Sprintf("%s said %s\n", ch.Name, command.contents),
		)

		return nil
	}
}

func LookCommandAction(command Command, ch *Character) WorldAction {
	return func(world *World) error {
		ch.Reply(fmt.Sprintf("You look around\n%s\n", world.DescribeRoom(ch.Coordinate)))

		return nil
	}
}

func GoCommandAction(command Command, ch *Character) WorldAction {
	return func(world *World) error {
		direction := DirectionFromString(command.contents)
		if command.contents == "" {
			ch.Reply("In which direction do you want to move?\n")
		} else if world.CanCharactorMoveInDirection(ch, direction) {
			// broadcast to old room
			world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s moved to %s\n", ch.Name, command.contents),
			)

			world.MoveCharacterInDirection(ch, direction)
			ch.Reply(
				fmt.Sprintf("You move to %s\n%s\n",
					command.contents,
					world.DescribeRoom(ch.Coordinate)),
			)

			// broadcast to new room
			world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s entered from %s\n", ch.Name, command.contents),
			)
		} else {
			ch.Reply("You cannot go that way!\n")
		}

		return nil
	}
}

func SmokeCommandAction(command Command, ch *Character) WorldAction {
	return func(world *World) error {
		switch command.contents {
		case "start":
			ch.SetState("smoking")
			ch.Reply("You started to smoke your pipe\n")

			world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s started to smoke a pipe\n", ch.Name),
			)
		case "stop":
			ch.SetState("idle")
			ch.Reply("You stopped smoking your pipe\n")

			world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%s stopped smoking a pipe\n", ch.Name),
			)
		default:
			ch.Reply(fmt.Sprintln("You either start or stop"))
		}

		return nil
	}
}

func HelpCommandAction(command Command, ch *Character) WorldAction {
	return func(world *World) error {
		var output = "help:\n"
		for _, cwd := range ch.commands.CommandsWithDescriptions() {
			output = fmt.Sprintf("%s\t%s\n", output, cwd)
		}

		ch.Reply(output)

		return nil
	}
}
