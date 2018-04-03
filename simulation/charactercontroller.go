package simulation

import "github.com/soupstore/coda-world/simulation/model"

// CharacterController is an interface over the Simulation that exposes all actions a connected
// player will need to perform
type CharacterController interface {
	WakeUpCharacter(model.CharacterID) <-chan interface{}
	SleepCharacter(model.CharacterID)
}

// WakeUpCharacter make a character wake up.
// It sends a room description to the waking character.
// It sends a character waking event to the other characters in the room.
func (s *Simulation) WakeUpCharacter(id model.CharacterID) <-chan interface{} {
	actor, ok := s.characters[id]
	if !ok {
		// error
	}

	// wake character and send description
	actor.Awake = true
	actor.Dispatch(model.EvtRoomDescription{Room: actor.Room})

	// send character wakes up
	wakeUpEvent := model.EvtCharacterWakesUp{Character: actor}
	for _, c := range actor.Room.GetCharacters() {
		// ignore the character that woke up
		if c == actor {
			continue
		}

		c.Dispatch(wakeUpEvent)
	}

	return actor.Events
}

// SleepCharacter sets a character to sleeping.
// It sends a character sleeping event to all other characters in the room.
func (s *Simulation) SleepCharacter(id model.CharacterID) {
	actor, ok := s.characters[id]
	if !ok {
		// error
	}

	actor.Awake = false

	// send character sleeps
	sleepEvent := model.EvtCharacterFallsAsleep{Character: actor}
	for _, c := range actor.Room.GetCharacters() {
		c.Dispatch(sleepEvent)
	}
}
