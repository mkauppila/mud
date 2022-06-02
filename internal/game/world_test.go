package game

import (
	"testing"
)

func TestBasicWorldFunctionality(t *testing.T) {
	w := NewWorld()
	var clientId ClientId = "clientId"
	var clientId2 ClientId = "clientId2"

	w.InsertCharacterOnConnect(NewCharacter(clientId, "abel"))
	ch := w.GetCharacter(clientId)

	if o := w.OtherCharactersInRoom(ch); len(o) > 0 {
		t.Fatal("There should not be other characters in the room")
	}

	w.InsertCharacterOnConnect(NewCharacter(clientId2, "bella"))
	if o := w.OtherCharactersInRoom(ch); len(o) != 0 && o[0].Id != clientId2 {
		t.Fatal("Bella should be in the same room")
	}

	w.MoveCharacterInDirection(w.GetCharacter(clientId2), East)
	if o := w.OtherCharactersInRoom(ch); len(o) > 0 {
		t.Fatal("There should not be other characters in the room")
	}

	if w.GetCharacter(clientId).Coordinate == w.GetCharacter(clientId2).Coordinate {
		t.Fatal("after movement characters are not in the same location")
	}

	w.RemoveCharacterOnDisconnect(w.GetCharacter(clientId))
	if (w.GetCharacter(clientId)) != nil {
		t.Fatalf("Character not removed. id: %s", clientId)
	}

	w.RemoveCharacterOnDisconnect(w.GetCharacter(clientId2))
	if (w.GetCharacter(clientId2)) != nil {
		t.Fatalf("Character not removed. id: %s", clientId2)
	}
}
