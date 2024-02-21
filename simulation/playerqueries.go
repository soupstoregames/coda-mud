package simulation

import (
	"github.com/soupstoregames/coda-mud/simulation/model"
)

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

	actor.Dispatch(model.EvtInventoryDescription{Character: actor})
	return nil
}
