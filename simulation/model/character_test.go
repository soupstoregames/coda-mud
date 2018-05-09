package model_test

import (
	"sync"
	"testing"

	"github.com/soupstore/coda-world/simulation/model"
)

func TestNewCharacter(t *testing.T) {
	spawnRoom := model.NewRoom(model.RoomID(4433), model.WorldID("test"), 0, "Spawn Room", "")
	character := model.NewCharacter(model.CharacterID(3453), "Test Name", spawnRoom)

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

	expected := model.EvtItemNotHere{}
	var result interface{}

	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		result = <-character.Events
		wait.Done()
	}()

	if expected != result {
		t.Error("Correct event did not come through character event channel")
	}
}
