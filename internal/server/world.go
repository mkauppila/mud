package server

import "time"

type World struct {
	characters map[Coordinate][]*Character
	rooms      map[Coordinate]Room
}

func NewWorld() *World {
	world := &World{
		characters: make(map[Coordinate][]*Character),
		rooms:      make(map[Coordinate]Room),
	}

	for _, room := range BasicMap() {
		world.rooms[room.location] = room
	}

	return world
}

func (w World) InsertCharacterOnConnect(character *Character) {
	loc := character.Coordinate

	list, ok := w.characters[loc]
	if !ok {
		w.characters[loc] = []*Character{character}
	} else {
		list = append(list, character)
		w.characters[loc] = list
	}
}

func (w World) OtherCharactersInRoom(currentCharacter *Character) []*Character {
	inRoom := w.characters[currentCharacter.Coordinate]

	var others []*Character
	for _, ch := range inRoom {
		if ch.id != currentCharacter.id {
			others = append(others, ch)
		}
	}
	return others
}

func (w World) BroadcastToOtherCharactersInRoom(currentCh *Character, message string) {
	inRoom := w.characters[currentCh.Coordinate]

	for _, ch := range inRoom {
		if ch.id != currentCh.id {
			ch.Broadcast(message)
		}
	}
}

func (w World) getCharacter(id ClientId) *Character {
	for _, chs := range w.characters {
		for _, ch := range chs {
			if ch.id == id {
				return ch
			}
		}
	}
	return nil
}

func (w World) RemoveCharacterOnDisconnect(ch *Character) {
	// remove the disconnecting ch from the room
	chs := w.characters[ch.Coordinate]
	for i, c := range chs {
		if c.id == ch.id {
			chs[i] = chs[(len(chs) - 1)]
			chs = chs[:len(chs)-1]
			break
		}
	}
	w.characters[ch.Coordinate] = chs
}

func (w World) CanCharactorMoveInDirection(character *Character, direction Direction) bool {
	newLoc := CoordinateInDirection(character.Coordinate, direction)
	if _, ok := w.rooms[newLoc]; !ok {
		return false
	} else {
		return true
	}
}

func (w World) MoveCharacterInDirection(character *Character, direction Direction) {
	old := NewCoordinate(character.Coordinate.X, character.Y)
	new := CoordinateInDirection(old, direction)

	// add to new location
	list, ok := w.characters[new]
	if !ok {
		w.characters[new] = []*Character{character}
	} else {
		list = append(list, character)
		w.characters[new] = list
	}
	character.Coordinate = new

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
		ch.Tick(timeStep, w)
	}
}

func (w World) DescribeRoom(location Coordinate) string {
	return w.rooms[location].description
}
