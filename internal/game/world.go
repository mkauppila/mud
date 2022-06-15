package game

import (
	"fmt"
	"time"
)

type Worlder interface {
	ClientJoined(
		clientId ClientId,
		directReply func(message string),
		reply func(message string),
		broadcast func(message string),
	)
	ClientDisconnected(ClientId) error
	PassMessageToClient(string, ClientId)
}

type World struct {
	accounts   []*Account
	characters map[Coordinate][]*Character
	rooms      map[Coordinate]Room
	timeStep   time.Duration
	actions    chan WorldAction
}

func (w *World) GetAccount(clientId ClientId) *Account {
	for _, a := range w.accounts {
		if a.id == clientId {
			return a
		}
	}
	return nil
}

func NewWorld() *World {
	world := &World{
		characters: make(map[Coordinate][]*Character),
		rooms:      make(map[Coordinate]Room),
		timeStep:   time.Second,
		actions:    make(chan WorldAction),
		accounts:   make([]*Account, 0),
	}

	for _, room := range BasicMap() {
		world.rooms[room.location] = room
	}

	return world
}

func (w *World) ClientJoined(
	clientId ClientId,
	directReply func(messasage string),
	reply func(message string),
	broadcast func(message string),
) {
	account := NewAccount(clientId, directReply, reply, broadcast)
	w.accounts = append(w.accounts, account)
	account.directReply("What's the character?\n > \n")
}

func (world *World) ClientDisconnected(clientId ClientId) error {
	account := world.GetAccount(clientId)
	if ch := account.loggedInCharacter; ch != nil {
		if ch := world.GetCharacter(ClientId(clientId)); ch != nil {
			world.RemoveCharacterOnDisconnect(ch)
			world.BroadcastToOtherCharactersInRoom(
				ch,
				fmt.Sprintf("%v disconnecting...\n", ch.Name),
			)
		} else {
			return ErrUnknownCharacter{id: clientId, action: "disconnecting"}
		}
	}

	for i, acc := range world.accounts {
		if acc.id == clientId {
			world.accounts[i] = world.accounts[len(world.accounts)-1]
			world.accounts = world.accounts[:len(world.accounts)-1]

			break
		}
	}

	return nil
}

func (world *World) handleAccountMessage(account *Account, msg string) {
	ch := NewCharacter(ClientId(account.id), msg)
	ch.commands = NewInGameCommandRegistry()
	ch.Reply = account.reply
	ch.Broadcast = account.broadcast
	ch.SetState("idle")
	account.loggedInCharacter = ch
	world.InsertCharacterOnConnect(ch)

	world.BroadcastToOtherCharactersInRoom(
		ch,
		fmt.Sprintf("%v joined!\n", ch.Name),
	)

	action := ch.commands.InputToAction("look", account.loggedInCharacter)
	world.actions <- action
}

func (w *World) handleCharacterMessasge(ch *Character, msg string) {
	action := ch.commands.InputToAction(msg, ch)
	w.actions <- action
}

func (world *World) PassMessageToClient(msg string, clientId ClientId) {
	if account := world.GetAccount(clientId); account != nil {
		if account.loggedInCharacter != nil {
			world.handleCharacterMessasge(account.loggedInCharacter, msg)
		} else {
			world.handleAccountMessage(account, msg)
		}
	}
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

func (w *World) InsertCharacterOnConnect(character *Character) {
	loc := character.Coordinate

	list, ok := w.characters[loc]
	if !ok {
		w.characters[loc] = []*Character{character}
	} else {
		list = append(list, character)
		w.characters[loc] = list
	}
}

func (w *World) OtherCharactersInRoom(currentCharacter *Character) []*Character {
	inRoom := w.characters[currentCharacter.Coordinate]

	var others []*Character
	for _, ch := range inRoom {
		if ch.Id != currentCharacter.Id {
			others = append(others, ch)
		}
	}
	return others
}

func (w *World) BroadcastToOtherCharactersInRoom(currentCh *Character, message string) {
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

func (w *World) RemoveCharacterOnDisconnect(ch *Character) {
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
