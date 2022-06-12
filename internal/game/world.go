package game

import (
	"fmt"
	"time"
)

type Worlder interface {
	ClientJoined(clientId ClientId, reply func(message string), broadcast func(message string))
	ClientDisconnected(ClientId) error
	PassMessageToClient(string, ClientId)
}

type World struct {
	lobbyCharacters []*Character
	characters      map[Coordinate][]*Character
	rooms           map[Coordinate]Room
	timeStep        time.Duration
	actions         chan WorldAction
}

func NewWorld() *World {
	world := &World{
		characters: make(map[Coordinate][]*Character),
		rooms:      make(map[Coordinate]Room),
		timeStep:   time.Second,
		actions:    make(chan WorldAction),
	}

	for _, room := range BasicMap() {
		world.rooms[room.location] = room
	}

	return world
}
func (w *World) ClientJoined(clientId ClientId, reply func(message string), broadcast func(message string)) {
	// func ConnectCommandAction(command Command, clientId ClientId) WorldAction {
	// 	return func(w *World) error {
	// 		client := s.getClient(clientId)
	// 		if client == nil {
	// 			return ErrUnknownClientId{id: clientId}
	// 		}
	// 		client.reply <- "Welcome! What is your name?\n"
	// 		return nil
	// 	}
	// }
	ch := NewCharacter(ClientId(clientId), "connected")
	ch.commands = NewLoginCommandRegistry()
	ch.Reply = reply
	ch.Broadcast = broadcast

	// w.lobbyCharacters = append(w.lobbyCharacters, ch)
	w.InsertCharacterOnConnect(ch)
	// reply("Welcome! What is your name?\n")
}

func (world *World) ClientDisconnected(clientId ClientId) error {
	ch := world.GetCharacter(clientId)
	if ch == nil {
		panic("no client")
	}

	if ch := world.GetCharacter(ClientId(clientId)); ch != nil {
		world.RemoveCharacterOnDisconnect(ch)
		world.BroadcastToOtherCharactersInRoom(
			ch,
			fmt.Sprintf("%v disconnecting...\n", ch.Name),
		)
	} else {
		return ErrUnknownCharacter{id: clientId, action: "disconnecting"}
	}

	return nil
}

func (world *World) PassMessageToClient(msg string, clientId ClientId) {
	// here we would need to check the state of the player (it's playr id now)
	// and then handle the incoming message in a appropriate way

	ch := world.GetCharacter(clientId)
	if ch == nil {
		// for _, c := range w.lobbyCharacters {
		// 	if c.Id == clientId {
		// 		ch = c
		//
		//  }
	}

	fmt.Println("message, clientId ", msg, clientId)
	cmd := ch.commands.InputToAction(msg, clientId)
	fmt.Println("cmd: ", cmd)
	world.actions <- cmd
}

func (w *World) RunGameLoop() {
	ticker := time.NewTicker(w.timeStep)
	defer ticker.Stop()

	var actions []WorldAction
	for {
		select {
		case command, ok := <-w.actions:
			if !ok {
				panic("actions channel closed")
			}
			actions = append(actions, command)
		case _, ok := <-ticker.C:
			if !ok {
				panic("ticker closed")
			}

			for _, action := range actions {
				err := action(w)
				if err != nil {
					panic(err)
				}
			}

			w.UpdateCharacterStates(w.timeStep)
			actions = make([]WorldAction, 0)
		}
	}
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
		if ch.Id != currentCharacter.Id {
			others = append(others, ch)
		}
	}
	return others
}

func (w World) BroadcastToOtherCharactersInRoom(currentCh *Character, message string) {
	inRoom := w.characters[currentCh.Coordinate]

	for _, ch := range inRoom {
		if ch.Id != currentCh.Id {
			ch.Broadcast(message)
		}
	}
}

func (w World) GetCharacter(id ClientId) *Character {
	for _, chs := range w.characters {
		for _, ch := range chs {
			if ch.Id == id {
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
		if c.Id == ch.Id {
			chs[i] = chs[(len(chs) - 1)]
			chs = chs[:len(chs)-1]
			break
		}
	}
	w.characters[ch.Coordinate] = chs
}

func (w World) CanCharactorMoveInDirection(character *Character, direction Direction) bool {
	return w.rooms[character.Coordinate].exits&direction != 0
}

func (w World) MoveCharacterInDirection(character *Character, direction Direction) {
	old := NewCoordinate(character.Coordinate.X, character.Coordinate.Y)
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
			if character.Name == c.Name {
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
	room := w.rooms[location]
	return fmt.Sprintf("%s\n%s\n", room.description, DirectionAsStrings(room.exits))
}
