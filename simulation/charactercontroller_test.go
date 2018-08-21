package simulation_test

import (
	"testing"

	"github.com/soupstore/coda/simulation"
	"github.com/soupstore/coda/simulation/model"
)

// In this test, we create one room and two characters, Sleepy and Grumpy.
// Sleepy will wake up. Sleepy gets a room description. Grumpy isnt awake so doesnt get an event.
// Grumpy will wake up. Grumpy gets a room description and Sleepy gets a character wakes up.
// Grump then goes to sleep again, Sleepy gets a character goes to sleep.
func TestWakingAndSleepingCharacter(t *testing.T) {
	// set up simulation
	sim := simulation.NewSimulation(nil)
	sim.CreateWorld(model.WorldID("test"))
	sim.CreateRoom(model.WorldID("test"), model.RoomID(0), "Void", "", "Nothing", "")

	if err := sim.SetSpawnRoom("test", 0); err != nil {
		t.Error(err)
	}

	sleepyID := sim.MakeCharacter("Sleepy")
	grumpyID := sim.MakeCharacter("Grumpy")
	target := simulation.CharacterController(sim)

	// wake up sleepy
	sleepyEvents, err := target.WakeUpCharacter(sleepyID)
	if err != nil {
		t.Error(err)
	}

	// assert sleepy gets rooms description
	event := <-sleepyEvents
	roomDescriptionEvent, ok := event.(model.EvtRoomDescription)
	if !ok {
		t.Error("Event not of type EvtRoomDescription")
	}
	if roomDescriptionEvent.Room.ID != 0 {
		t.Errorf("Expected room description event to contain room %d, but got room %d", 0, roomDescriptionEvent.Room.ID)
	}

	// wake up grumpy
	grumpyEvents, err := target.WakeUpCharacter(grumpyID)
	if err != nil {
		t.Error(err)
	}

	// assert sleepy gets a wake up event
	event = <-sleepyEvents
	wakeupEvent, ok := event.(model.EvtCharacterWakesUp)
	if !ok {
		t.Error("Event not of type EvtCharacterWakesUp")
	}
	if wakeupEvent.Character.ID != grumpyID {
		t.Errorf("Expected wake up event to be about character %s, but got character %s", grumpyID, wakeupEvent.Character.ID)
	}

	// assert grumpy gets rooms description
	event = <-grumpyEvents
	roomDescriptionEvent, ok = event.(model.EvtRoomDescription)
	if !ok {
		t.Error("Event not of type EvtRoomDescription")
	}
	if roomDescriptionEvent.Room.ID != 0 {
		t.Errorf("Expected room description event to contain room %d, but got room %d", 0, roomDescriptionEvent.Room.ID)
	}

	// send grumpy to sleep
	target.SleepCharacter(grumpyID)

	// assert sleepy gets a sleep event
	event = <-sleepyEvents
	sleepEvent, ok := event.(model.EvtCharacterFallsAsleep)
	if !ok {
		t.Error("Event not of type EvtCharacterFallsAsleep")
	}
	if sleepEvent.Character.ID != grumpyID {
		t.Errorf("Expected sleep event to be about character %s, but got character %s", grumpyID, wakeupEvent.Character.ID)
	}
}

// In this test we create an empty simulation and attempt to wake up a character.
// That character doesnt exist so we get an error.
func TestWakeUpWithUnknownCharacter(t *testing.T) {
	sim := simulation.NewSimulation(nil)
	_, err := sim.WakeUpCharacter("")
	if err != simulation.ErrCharacterNotFound {
		t.Error("Did not get expected error")
	}
}

// In this test we create an empty simulation and attempt to sleep a character.
// That character doesnt exist so we get an error.
func TestSleepWithUnknownCharacter(t *testing.T) {
	sim := simulation.NewSimulation(nil)
	err := sim.SleepCharacter("")
	if err != simulation.ErrCharacterNotFound {
		t.Error("Did not get expected error")
	}
}

// Waking up an awake character implies that someone is connecting to a character
// that has already been connected to. This is an error.
func TestWakeUpWithAwakeCharacter(t *testing.T) {
	sim := simulation.NewSimulation(nil)
	sim.CreateWorld(model.WorldID("test"))
	sim.CreateRoom(model.WorldID("test"), model.RoomID(0), "Void", "", "Nothing", "")

	if err := sim.SetSpawnRoom("test", 0); err != nil {
		t.Error(err)
	}
	sleepyID := sim.MakeCharacter("Sleepy")
	sim.WakeUpCharacter(sleepyID)
	_, err := sim.WakeUpCharacter(sleepyID)
	if err != simulation.ErrCharacterAwake {
		t.Error("Did not get expected error")
	}
}

// Sleeping a character that is already asleep means that someone has disconnected
// from this character twice. This is an error.
func TestSleepWithSleepingCharacter(t *testing.T) {
	sim := simulation.NewSimulation(nil)
	sim.CreateWorld(model.WorldID("test"))
	sim.CreateRoom(model.WorldID("test"), model.RoomID(0), "Void", "", "Nothing", "")

	if err := sim.SetSpawnRoom("test", 0); err != nil {
		t.Error(err)
	}
	sleepyID := sim.MakeCharacter("Sleepy")
	err := sim.SleepCharacter(sleepyID)
	if err != simulation.ErrCharacterAsleep {
		t.Error("Did not get expected error")
	}
}
