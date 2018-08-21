package simulation

import (
	"github.com/soupstore/coda/data/state"
	"github.com/soupstore/coda/simulation/model"
)

// Simulation is the engine of the world.
// It holds rooms, characters, items etc...
// It exposes a number of interfaces to manipulate the simulation.
type Simulation struct {
	persister       state.Persister
	spawnRoom       *model.Room
	worlds          map[model.WorldID]*model.World
	itemDefinitions map[model.ItemDefinitionID]*model.ItemDefinition
	items           map[model.ItemID]*model.Item
	characters      map[model.CharacterID]*model.Character
	containers      map[model.ContainerID]model.Container
}

// NewSimulation returns a Simulation with default params.
func NewSimulation(persister state.Persister) *Simulation {
	return &Simulation{
		persister:       persister,
		spawnRoom:       nil,
		worlds:          make(map[model.WorldID]*model.World),
		itemDefinitions: make(map[model.ItemDefinitionID]*model.ItemDefinition),
		items:           make(map[model.ItemID]*model.Item),
		characters:      make(map[model.CharacterID]*model.Character),
		containers:      make(map[model.ContainerID]model.Container),
	}
}
