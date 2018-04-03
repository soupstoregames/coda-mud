package simulation_test

import (
	"testing"

	"github.com/soupstore/coda-world/simulation"
	"github.com/soupstore/coda-world/simulation/model"
)

// In this test, we create one room and two characters, Sleepy and Grumpy.
// Sleepy will wake up. Sleepy gets a room description. Grumpy isnt awake so doesnt get an event.
// Grumpy will wake up. Grumpy gets a room description and Sleepy gets a character wakes up.
// Grump then goes to sleep again, Sleepy gets a character goes to sleep.
func TestWakingAndSleepingCharacter(t *testing.T) {
	// set up simulation
	sim := simulation.NewSimulation()

	// set up world
	wc := simulation.WorldController(sim)
	roomID := wc.MakeRoom("Void", "Nothing")
	wc.SetSpawnRoom(roomID)

	// create our actor
	rc := simulation.RegistrationController(sim)
	sleepyID := rc.MakeCharacter("Sleepy")
	grumpyID := rc.MakeCharacter("Grumpy")

	// wake up sleepy
	cc := simulation.CharacterController(sim)
	sleepyEvents := cc.WakeUpCharacter(sleepyID)

	// assert sleepy gets rooms description
	event := <-sleepyEvents
	roomDescriptionEvent, ok := event.(model.EvtRoomDescription)
	if !ok {
		t.Error("Event not of type EvtRoomDescription")
	}
	if roomDescriptionEvent.Room.ID != roomID {
		t.Errorf("Expected room description event to contain room %d, but got room %d", roomID, roomDescriptionEvent.Room.ID)
	}

	// wake up grumpy
	grumpyEvents := cc.WakeUpCharacter(grumpyID)

	// assert sleepy gets a wake up event
	event = <-sleepyEvents
	wakeupEvent, ok := event.(model.EvtCharacterWakesUp)
	if !ok {
		t.Error("Event not of type EvtCharacterWakesUp")
	}
	if wakeupEvent.Character.ID != grumpyID {
		t.Errorf("Expected wake up event to be about character %d, but got character %d", grumpyID, wakeupEvent.Character.ID)
	}

	// assert grumpy gets rooms description
	event = <-grumpyEvents
	roomDescriptionEvent, ok = event.(model.EvtRoomDescription)
	if !ok {
		t.Error("Event not of type EvtRoomDescription")
	}
	if roomDescriptionEvent.Room.ID != roomID {
		t.Errorf("Expected room description event to contain room %d, but got room %d", roomID, roomDescriptionEvent.Room.ID)
	}

	// send grumpy to sleep
	cc.SleepCharacter(grumpyID)

	// assert sleepy gets a sleep event
	event = <-sleepyEvents
	sleepEvent, ok := event.(model.EvtCharacterFallsAsleep)
	if !ok {
		t.Error("Event not of type EvtCharacterFallsAsleep")
	}
	if sleepEvent.Character.ID != grumpyID {
		t.Errorf("Expected sleep event to be about character %d, but got character %d", grumpyID, wakeupEvent.Character.ID)
	}
}
