package simulation

import (
	"github.com/soupstore/coda-world/simulation/model"
)

// CharacterController is an interface over the Simulation that exposes all actions a connected
// player will need to perform
type CharacterController interface {
	WakeUpCharacter(model.CharacterID) (<-chan interface{}, error)
	SleepCharacter(model.CharacterID) error
}

// WakeUpCharacter make a character wake up.
// It sends a room description to the waking character.
// It sends a character waking event to the other characters in the room.
func (s *Simulation) WakeUpCharacter(id model.CharacterID) (<-chan interface{}, error) {
	actor, ok := s.characters[id]
	if !ok {
		return nil, ErrCharacterNotFound
	}

	if actor.Awake {
		return nil, ErrCharacterAlreadyAwake
	}

	// wake character and send description
	actor.Events = make(chan interface{}, 1)
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

	s.logger.Debug("Character woke up")

	return actor.Events, nil
}

// SleepCharacter sets a character to sleeping.
// It sends a character sleeping event to all other characters in the room.
func (s *Simulation) SleepCharacter(id model.CharacterID) error {
	actor, ok := s.characters[id]
	if !ok {
		return ErrCharacterNotFound
	}

	if !actor.Awake {
		return ErrCharacterAlreadyAsleep
	}

	actor.Awake = false
	close(actor.Events)

	// send character sleeps
	sleepEvent := model.EvtCharacterFallsAsleep{Character: actor}
	for _, c := range actor.Room.GetCharacters() {
		c.Dispatch(sleepEvent)
	}

	s.logger.Debug("Character fell asleep")

	return nil
}
