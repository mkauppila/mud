package mud

import (
	"testing"

	"github.com/google/uuid"
)

func TestMovingCharacter(t *testing.T) {
	world := NewWorld()
	id, _ := uuid.NewRandom()
	character := NewCharacter(id)
	world.InsertCharacterOnConnect(character)

	ch := world.characters[Location{X: 0, Y: 0}][0]
	if ch.name != character.name {
		t.Fatalf("Names don't match")
	}

	world.MoveCharacterInDirection(character, "west")
	world.MoveCharacterInDirection(character, "west")
	ch = world.characters[Location{X: -2, Y: 0}][0]
	if ch.name != character.name {
		t.Fatalf("Names don't match")
	}

	world.MoveCharacterInDirection(character, "north")
	ch = world.characters[Location{X: -2, Y: -1}][0]
	if ch.name != character.name {
		t.Fatalf("Names don't match")
	}

	world.MoveCharacterInDirection(character, "east")
	ch = world.characters[Location{X: -1, Y: -1}][0]
	if ch.name != character.name {
		t.Fatalf("Names don't match")
	}

	world.MoveCharacterInDirection(character, "south")
	ch = world.characters[Location{X: -1, Y: 0}][0]
	if ch.name != character.name {
		t.Fatalf("Names don't match")
	}
}

/*
func TestMovingMultipleCharacter(t *testing.T) {
	world := NewWorld()
	ch1 := NewCharacter()
	ch2 := NewCharacter()
	world.InsertCharacterOnJoin(ch1)
	world.InsertCharacterOnJoin(ch2)

	ch := world.characters["0:0"][0]
	if ch.name != ch1.name {
		t.Fatalf("Names don't match")
	}

	world.MoveCharacterInDirection(ch1, "west")
	world.MoveCharacterInDirection(ch1, "west")
	ch = world.characters["-2:0"][0]
	if ch.name != ch1.name {
		t.Fatalf("Names don't match")
	}

	world.MoveCharacterInDirection(ch1, "north")
	ch = world.characters["-2:-1"][0]
	if ch.name != ch1.name {
		t.Fatalf("Names don't match")
	}

	world.MoveCharacterInDirection(ch1, "east")
	ch = world.characters["-1:-1"][0]
	if ch.name != ch1.name {
		t.Fatalf("Names don't match")
	}

	world.MoveCharacterInDirection(ch1, "south")
	ch = world.characters["-1:0"][0]
	if ch.name != ch1.name {
		t.Fatalf("Names don't match")
	}
}
*/
