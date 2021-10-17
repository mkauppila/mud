package mud

import (
	"math/rand"

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
}

func NewCharacter(id uuid.UUID) *Character {
	names := []string{"Matt", "John", "Ugruk", "Sonya", "Miraboile"}
	return &Character{
		id:       id,
		health:   30,
		attack:   1,
		name:     names[rand.Intn(len(names))],
		Location: Location{X: 0, Y: 0},
	}
}

type Room struct {
	description string
}

type World struct {
	characters map[Location][]*Character
	rooms      map[Location][]Room
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

func NewWorld() *World {
	return &World{
		characters: make(map[Location][]*Character),
		rooms:      make(map[Location][]Room),
	}
}
