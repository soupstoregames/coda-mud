package simulation

import (
	"github.com/soupstore/coda/simulation/model"
	"github.com/soupstore/go-core/logging"
)

// CharacterController is an interface over the Simulation that exposes all actions a connected
// player will need to perform
type CharacterController interface {
	WakeUpCharacter(model.CharacterID) (characterEvents <-chan interface{}, err error)
	SleepCharacter(model.CharacterID) error
	QueueCommand(model.CharacterID, interface{}) error
	Look(model.CharacterID) error
	Inventory(model.CharacterID) error
}

func (s *Simulation) QueueCommand(id model.CharacterID, command interface{}) error {
	char, err := s.findAwakeCharacter(id)
	if err != nil {
		return err
	}

	char.Commands <- command

	return nil
}

// WakeUpCharacter make a character wake up.
// It sends a room description to the waking character.
// It sends a character waking event to the other characters in the room.
func (s *Simulation) WakeUpCharacter(id model.CharacterID) (characterEvents <-chan interface{}, err error) {
	actor, ok := s.characters[id]
	if !ok {
		return nil, ErrCharacterNotFound
	}

	if actor.Awake {
		return nil, ErrCharacterAwake
	}

	// wake character and send description
	actor.WakeUp()
	actor.Dispatch(model.EvtRoomDescription{Room: actor.Room})

	// send character wakes up
	wakeUpEvent := model.EvtCharacterWakesUp{Character: actor}
	for _, c := range actor.Room.Characters {
		// ignore the character that woke up
		if c == actor {
			continue
		}

		c.Dispatch(wakeUpEvent)
	}

	actor.Room.OnWake(actor)

	logging.Info("Character woke up")

	return actor.Events, nil
}

// SleepCharacter sets a character to sleeping.
// It sends a character sleeping event to all other characters in the room.
func (s *Simulation) SleepCharacter(id model.CharacterID) error {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return err
	}

	actor.Sleep()

	// send character sleeps
	sleepEvent := model.EvtCharacterFallsAsleep{Character: actor}
	for _, c := range actor.Room.Characters {
		c.Dispatch(sleepEvent)
	}

	actor.Room.OnExit(actor)

	return nil
}

// Look gives the character a room description
func (s *Simulation) Look(id model.CharacterID) error {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return err
	}

	actor.Dispatch(model.EvtRoomDescription{Room: actor.Room})
	return nil
}

// Inventory lists the users inventory and items.
func (s *Simulation) Inventory(id model.CharacterID) error {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return err
	}

	actor.Dispatch(model.EvtInventoryDescription{Rig: actor.Rig})
	return nil
}

// this checks that the character exists in the simulation and that they are awake (connected to)
func (s *Simulation) findAwakeCharacter(id model.CharacterID) (*model.Character, error) {
	actor, ok := s.characters[id]
	if !ok {
		return nil, ErrCharacterNotFound
	}

	if !actor.Awake {
		return nil, ErrCharacterAsleep
	}

	return actor, nil
}
