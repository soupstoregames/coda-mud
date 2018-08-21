package simulation

import (
	"github.com/soupstore/coda/simulation/model"
	state "github.com/soupstore/coda/state-data"
)

func (s *Simulation) Save() error {
	for i := range s.characters {
		s.persister.QueueCharacter(characterToState(s.characters[i]))
	}

	return s.persister.Persist()
}

func characterToState(c *model.Character) state.Character {
	return state.Character{
		ID:    string(c.ID),
		Name:  c.Name,
		Room:  int64(c.Room.ID),
		World: string(c.Room.WorldID),
		Rig:   mapRig(c.Rig),
	}
}

func mapRig(r *model.Rig) state.Rig {
	return state.Rig{
		Backpack: mapItem(r.Backpack),
	}
}

func mapItem(i *model.Item) *state.Item {
	if i == nil {
		return nil
	}

	return &state.Item{
		ID:             string(i.ID),
		ItemDefinition: int64(i.Definition.ID),
	}
}
