package simulation

import (
	"github.com/soupstoregames/coda-mud/simulation/model"
	"sync"
)

// Simulation is the engine of the world.
// It holds rooms, characters, items etc...
// It exposes a number of interfaces to manipulate the simulation.
type Simulation struct {
	spawnRoom       *model.Room
	worlds          map[model.WorldID]*model.World
	itemDefinitions map[model.ItemDefinitionID]*model.ItemDefinition
	items           map[model.ItemID]*model.Item
	characters      map[model.CharacterID]*model.Character
	containers      map[model.ContainerID]model.Container

	characterLock *sync.Mutex
}

// NewSimulation returns a Simulation with default params.
func NewSimulation() *Simulation {
	return &Simulation{
		spawnRoom:       nil,
		worlds:          make(map[model.WorldID]*model.World),
		itemDefinitions: make(map[model.ItemDefinitionID]*model.ItemDefinition),
		items:           make(map[model.ItemID]*model.Item),
		characters:      make(map[model.CharacterID]*model.Character),
		containers:      make(map[model.ContainerID]model.Container),

		characterLock: &sync.Mutex{},
	}
}

func (s *Simulation) Start() {
	go func() {
		for {
			s.processPlayerCommands()
		}
	}()
}

func (s *Simulation) processPlayerCommands() {
	s.characterLock.Lock()
	defer s.characterLock.Unlock()

	// iterate through all connected characters to read commands
	for _, c := range s.characters {
		if !c.Awake {
			continue
		}

		select {
		case cmd := <-c.Commands:
			switch v := cmd.(type) {
			case model.CommandSay:
				s.say(c, v)
			case model.CommandMove:
				s.move(c, v)
			case model.CommandTake:
				s.takeItem(c, v)
			case model.CommandDrop:
				s.dropItem(c, v)
			case model.CommandEquip:
				s.equipItem(c, v)
			case model.CommandUnequip:
				s.unequipItem(c, v)
			}
		default:
			continue
		}
	}
}
