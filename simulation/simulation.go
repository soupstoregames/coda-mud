package simulation

import (
	"github.com/soupstore/coda-world/simulation/model"
	"go.uber.org/zap"
)

// Simulation is the engine of the world.
// It holds rooms, characters, items etc...
// It exposes a number of interfaces to manipulate the simulation.
type Simulation struct {
	logger          *zap.Logger
	nextRoomID      model.RoomID
	nextCharacterID model.CharacterID
	spawnRoom       *model.Room
	rooms           map[model.RoomID]*model.Room
	items           map[model.ItemID]*model.Item
	characters      map[model.CharacterID]*model.Character
}

// NewSimulation returns a Simulation with default params.
func NewSimulation(logger *zap.Logger) *Simulation {
	return &Simulation{
		logger:          logger,
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
	s.nextRoomID = roomID + 1
	return roomID
}

func (s *Simulation) getNextCharacterID() model.CharacterID {
	id := s.nextCharacterID
	s.nextCharacterID = s.nextCharacterID + 1
	return id
}
