package simulation

import (
	"github.com/soupstore/coda-world/data"
	"github.com/soupstore/coda-world/simulation/model"
)

// Simulation is the engine of the world.
// It holds rooms, characters, items etc...
// It exposes a number of interfaces to manipulate the simulation.
type Simulation struct {
	nextCharacterID model.CharacterID
	nextContainerID model.ContainerID
	spawnRoom       *model.Room
	worlds          map[string]*model.World
	items           map[model.ItemID]*model.Item
	characters      map[model.CharacterID]*model.Character
	containers      map[model.ContainerID]*model.Container
}

// NewSimulation returns a Simulation with default params.
func NewSimulation() *Simulation {
	return &Simulation{
		nextCharacterID: 0,
		nextContainerID: 0,
		spawnRoom:       nil,
		worlds:          make(map[string]*model.World),
		items:           make(map[model.ItemID]*model.Item),
		characters:      make(map[model.CharacterID]*model.Character),
		containers:      make(map[model.ContainerID]*model.Container),
	}
}

func (s *Simulation) LoadData(d *data.Data) error {
	adminWorld := d.Worlds["admin"]
	if err := s.AddWorld("admin"); err != nil {
		return err
	}

	// load all worlds
	for roomID, room := range adminWorld {
		if err := s.loadRoom("admin", model.RoomID(roomID), room); err != nil {
			return err
		}
	}

	// load all exits
	for roomID, room := range adminWorld {
		for direction, exit := range room.Exits {
			d, err := model.StringToDirection(direction)
			if err != nil {
				return err
			}
			s.LinkRoom("admin", model.RoomID(roomID), d, "admin", model.RoomID(exit.RoomID), false)
		}
	}

	return nil
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
