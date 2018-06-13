package simulation

import (
	"github.com/soupstore/coda/simulation/model"
)

// Simulation is the engine of the world.
// It holds rooms, characters, items etc...
// It exposes a number of interfaces to manipulate the simulation.
type Simulation struct {
	nextItemID      model.ItemID
	nextCharacterID model.CharacterID
	nextContainerID model.ContainerID
	spawnRoom       *model.Room
	worlds          map[model.WorldID]*model.World
	itemDefinitions map[model.ItemDefinitionID]*model.ItemDefinition
	items           map[model.ItemID]*model.Item
	characters      map[model.CharacterID]*model.Character
	containers      map[model.ContainerID]*model.Container
}

// NewSimulation returns a Simulation with default params.
func NewSimulation() *Simulation {
	return &Simulation{
		nextItemID:      0,
		nextCharacterID: 0,
		nextContainerID: 0,
		spawnRoom:       nil,
		worlds:          make(map[model.WorldID]*model.World),
		itemDefinitions: make(map[model.ItemDefinitionID]*model.ItemDefinition),
		items:           make(map[model.ItemID]*model.Item),
		characters:      make(map[model.CharacterID]*model.Character),
		containers:      make(map[model.ContainerID]*model.Container),
	}
}

func (s *Simulation) getNextItemID() model.ItemID {
	id := s.nextItemID
	s.nextItemID = s.nextItemID + 1
	return id
}

func (s *Simulation) getNextCharacterID() model.CharacterID {
	id := s.nextCharacterID
	s.nextCharacterID = s.nextCharacterID + 1
	return id
}

func (s *Simulation) getNextContainerID() model.ContainerID {
	id := s.nextContainerID
	s.nextContainerID = s.nextContainerID + 1
	return id
}
