package mud

import (
	"testing"
)

func TestBasicWorldFunctionality(t *testing.T) {
	w := NewWorld()
	var clientId ClientId = "clientId"
	var clientId2 ClientId = "clientId2"

	w.InsertCharacterOnConnect(NewCharacter(clientId, "abel"))
	ch := w.getCharacter(clientId)

	if o := w.OtherCharactersInRoom(ch); len(o) > 0 {
		t.Fatal("There should not be other characters in the room")
	}

	w.InsertCharacterOnConnect(NewCharacter(clientId2, "bella"))
	if o := w.OtherCharactersInRoom(ch); len(o) != 0 && o[0].id != clientId2 {
		t.Fatal("Bella should be in the same room")
	}

	w.MoveCharacterInDirection(w.getCharacter(clientId2), "east")
	if o := w.OtherCharactersInRoom(ch); len(o) > 0 {
		t.Fatal("There should not be other characters in the room")
	}

	if w.getCharacter(clientId).Location == w.getCharacter(clientId2).Location {
		t.Fatal("after movement characters are not in the same location")
	}

	w.RemoveCharacterOnDisconnect(w.getCharacter(clientId))
	if (w.getCharacter(clientId)) != nil {
		t.Fatalf("Character not removed. id: %s", clientId)
	}

	w.RemoveCharacterOnDisconnect(w.getCharacter(clientId2))
	if (w.getCharacter(clientId2)) != nil {
		t.Fatalf("Character not removed. id: %s", clientId2)
	}
}
