package simulation

import (
	"github.com/soupstore/coda/simulation/model"
)

func (s *Simulation) FindItemInRoom(id model.CharacterID, alias string) (*model.Item, error) {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return nil, err
	}

	if item := actor.Room.FindItem(alias); item != nil {
		return item, nil
	}

	return nil, ErrItemNotFound
}

func (s *Simulation) FindItemInInventory(id model.CharacterID, alias string) (*model.Item, error) {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return nil, err
	}

	if item := actor.SearchInventory(alias); item != nil {
		return item, nil
	}

	return nil, ErrItemNotFound
}

func (s *Simulation) FindItemInRig(id model.CharacterID, alias string) (*model.Item, error) {
	actor, err := s.findAwakeCharacter(id)
	if err != nil {
		return nil, err
	}

	if item := actor.Rig.FindItem(alias); item != nil {
		return item, nil
	}

	return nil, ErrItemNotFound
}
