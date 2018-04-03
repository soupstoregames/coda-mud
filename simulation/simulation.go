package simulation

import (
	"github.com/soupstore/coda-world/simulation/model"
)

type Simulation struct {
	nextRoomID      model.RoomID
	nextCharacterID model.CharacterID
	spawnRoom       *model.Room
	rooms           map[model.RoomID]*model.Room
	items           map[model.ItemID]*model.Item
	characters      map[model.CharacterID]*model.Character
}

func NewSimulation() *Simulation {
	return &Simulation{
		nextRoomID:      0,
		nextCharacterID: 0,
		spawnRoom:       nil,
		rooms:           make(map[model.RoomID]*model.Room),
		items:           make(map[model.ItemID]*model.Item),
		characters:      make(map[model.CharacterID]*model.Character),
	}
}

func (s *Simulation) getNextRoomID() model.RoomID {
	roomID := s.nextRoomID
	s.nextRoomID = s.nextRoomID + 1
	return roomID
}

func (s *Simulation) getNextCharacterID() model.CharacterID {
	id := s.nextCharacterID
	s.nextCharacterID = s.nextCharacterID + 1
	return id
}
