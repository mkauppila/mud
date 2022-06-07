package game

import (
	"fmt"
	"time"
)

type ClientId string

/*
character would have command registry
client would have a link to the character


*/
type Character struct {
	Id             ClientId
	health, attack int
	Name           string
	Coordinate     Coordinate

	Reply     func(string)
	Broadcast func(string)

	state    State
	commands *CommandRegistry
}

func NewCharacter(id ClientId, name string /*, reply func(string), broadcast func(string)*/) *Character {
	ch := &Character{
		Id:         id,
		health:     30,
		attack:     1,
		Name:       name,
		Coordinate: Coordinate{X: 0, Y: 0},
	}

	ch.SetState("idle")

	return ch
}

func (c *Character) Tick(timeStep time.Duration, world World) {
	c.state.Tick(c, world, timeStep)
}

func (c *Character) SetState(state CharacterState) {
	switch state {
	case idle:
		c.state = CreateIdleState()
	case smoking:
		c.state = CreateSmokingPipeState()
	default:
		fmt.Printf("unknown state: %s", state)
	}
}

func (c Character) String() string {
	return fmt.Sprintf("%s - %s\n", c.Id, c.Name)
}

type CharacterState string

const (
	idle    CharacterState = "idle"
	smoking CharacterState = "smoking"
)

type State struct {
	state       CharacterState
	timeLeft    time.Duration
	description string
	function    func(*Character, World, time.Duration)
}

func (state *State) Tick(ch *Character, world World, timeStep time.Duration) {
	state.function(ch, world, timeStep)
}

func CreateIdleState() State {
	return State{
		state:       idle,
		timeLeft:    time.Second,
		description: "X is standing idle",
		function: func(ch *Character, world World, timeStep time.Duration) {
		},
	}
}

func CreateSmokingPipeState() State {
	return State{
		state:       smoking,
		timeLeft:    time.Second * 5,
		description: "X is smoking a pipe",
		function: func(ch *Character, world World, timeStep time.Duration) {
			ch.state.timeLeft -= timeStep
			if ch.state.timeLeft > 0 {
				ch.Broadcast("The pipe puffs\n")

				world.BroadcastToOtherCharactersInRoom(
					ch,
					fmt.Sprintf("%s puffs the pipe\n", ch.Name),
				)
			} else {
				ch.Broadcast("You run out of tobacco and stopped smoking the pipe\n")
				ch.SetState("idle")

				world.BroadcastToOtherCharactersInRoom(
					ch,
					fmt.Sprintf("%s stopped smoking the pipe\n", ch.Name),
				)
			}
		},
	}
}
