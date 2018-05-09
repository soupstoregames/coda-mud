package model_test

import (
	"testing"

	"github.com/soupstore/coda-world/simulation/model"
)

func TestNewCharacter(t *testing.T) {
	spawnRoom := model.NewRoom(model.RoomID(4433), model.WorldID("test"), 0, "Spawn Room", "")
	character := model.NewCharacter(model.CharacterID(3453), "Test Name", spawnRoom)

	if character.ID != model.CharacterID(3453) {
		t.Error("Incorrect ID assigned")
	}

	if character.Name != "Test Name" {
		t.Error("Incorrect name set on character")
	}

	if character.Room != spawnRoom {
		t.Error("Character thinks its in the wrong room")
	}

	if character.Rig == nil {
		t.Error("Rig not initialized")
	}

	if character.Awake {
		t.Error("Characters are asleep until a client connects to them")
	}

	if character.Backpack != nil {
		t.Error("Character should spawn with no equipment")
	}
}

func TestDispatch(t *testing.T) {
	spawnRoom := model.NewRoom(model.RoomID(4433), model.WorldID("test"), 0, "Spawn Room", "")
	character := model.NewCharacter(model.CharacterID(3453), "Test Name", spawnRoom)
	character.Events = make(chan interface{})
	character.Awake = true

	expected := model.EvtItemNotHere{}
	var result interface{}

	go func() {
		character.Dispatch(expected)
	}()

	result = <-character.Events

	if expected != result {
		t.Error("Correct event did not come through character event channel")
	}
}
