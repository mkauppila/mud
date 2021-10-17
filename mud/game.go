package mud

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Location struct {
	X, Y int
}

type Character struct {
	id             uuid.UUID // same as clients uuid
	health, attack int
	name           string
	Location

	currentState State
}

func NewCharacter(id uuid.UUID) *Character {
	names := []string{"Matt", "John", "Ugruk", "Sonya", "Miraboile"}
	ch := &Character{
		id:       id,
		health:   30,
		attack:   1,
		name:     names[rand.Intn(len(names))],
		Location: Location{X: 0, Y: 0},
	}

	ch.SetState("idle")

	return ch
}

func (c *Character) Tick(timeStep time.Duration) {
	c.currentState.Tick(c, timeStep)
}

func (c *Character) SetState(state CharacterState) {
	switch state {
	case idle:
		c.currentState = CreateIdleState()
	case smoking:
		c.currentState = CreateSmokingtate()
	default:
		fmt.Printf("unknown state: %s", state)
	}
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
	function    func(*Character, time.Duration)
}

func (state *State) Tick(ch *Character, timeStep time.Duration) {
	state.function(ch, timeStep)
}

func CreateIdleState() State {
	return State{
		state:       idle,
		timeLeft:    time.Second, // is not needed
		description: "X is standing idle",
		function: func(ch *Character, timeStep time.Duration) {
			fmt.Printf("%s is standing idle\n", ch.name)
		},
	}
}

func CreateSmokingtate() State {
	return State{
		state:       smoking,
		timeLeft:    time.Second,
		description: "X is smoking a pipe", // this could also be a function based on character struct
		function: func(ch *Character, timeStep time.Duration) {
			fmt.Printf("%s is smoking a pipe\n", ch.name)
			// timeLeft - time.duration
			// if timeLeft < 0:
			//   go to next state setState("idle")
			// do the state logic what ever that is!
		},
	}
}

type Room struct {
	description string
}

type World struct {
	characters map[Location][]*Character
	rooms      map[Location][]Room
}

func NewWorld() *World {
	return &World{
		characters: make(map[Location][]*Character),
		rooms:      make(map[Location][]Room),
	}
}

func (w World) InsertCharacterOnConnect(character *Character) {
	loc := character.Location

	list, ok := w.characters[loc]
	if !ok {
		w.characters[loc] = []*Character{character}
	} else {
		list = append(list, character)
		w.characters[loc] = list
	}
}

func (w World) OtherCharactersInRoom(currentCharacter *Character) []*Character {
	inRoom := w.characters[currentCharacter.Location]

	var others []*Character
	for _, ch := range inRoom {
		if ch.id != currentCharacter.id {
			others = append(others, ch)
		}
	}
	return others
}

func (w World) RemoveCharacterOnDisconnect(ch Character) {
	delete(w.characters, ch.Location)
}

func (w World) CanCharactorMoveInDirection(character *Character, direction string) bool {
	// TODO check if movement is allowed
	// if the rooms map has a room at this coord, then it's okay
	// otherwise block the movement here and modify the players respond
	// -> "Ouch, it seems the world has some boundaries"

	// There's no rooms or boundaries yet so it's always allowed
	return true
}

func (w World) MoveCharacterInDirection(character *Character, direction string) {
	old := character.Location

	switch direction {
	case "west":
		character.X--
	case "east":
		character.X++
	case "north":
		character.Y--
	case "south":
		character.Y++
	}

	new := character.Location

	// add to new location
	list, ok := w.characters[new]
	if !ok {
		w.characters[new] = []*Character{character}
	} else {
		list = append(list, character)
		w.characters[new] = list
	}

	// remove character from old
	list, ok = w.characters[old]
	if ok {
		ic := -1
		for i, c := range list {
			if character.name == c.name {
				ic = i
				break
			}
		}

		if ic > -1 {
			newLen := len(list) - 1
			list[ic] = list[newLen]
			if newLen > 0 {
				w.characters[old] = list[:newLen]
			} else {
				delete(w.characters, old)
			}
		}
	}
}

func (w World) UpdateCharacterStates(timeStep time.Duration) {
	var allChs []*Character
	for _, c := range w.characters {
		allChs = append(allChs, c...)
	}

	for _, ch := range allChs {
		ch.Tick(timeStep)
	}
}
