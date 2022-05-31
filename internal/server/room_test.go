package server

import "testing"

func TestRoomCreation(t *testing.T) {
	room := NewRoom("desc", NewCoordinate(0, 0))

	if room.description != "desc" {
		t.Fatal("Description now match")
	}

	if room.location.X != 0 || room.location.Y != 0 {
		t.Fatal("location is incorrect")
	}
}

func TestRoomMultipleExits(t *testing.T) {
	room := Room{}
	room.exits = East | South | North

	t.Logf("exits %b", room.exits)

	if !room.HasExitInDirection(North) {
		t.Fatal("north exit should exist")
	}
	if !room.HasExitInDirection(South) {
		t.Fatal("south exit should exist")
	}
	if !room.HasExitInDirection(East) {
		t.Fatal("east exit should exist")
	}

	if room.HasExitInDirection(West) {
		t.Fatal("west exit should NOT exist")
	}
}
