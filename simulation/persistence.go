package simulation

import (
	"github.com/soupstore/coda/data/state"
	"github.com/soupstore/coda/simulation/model"
)

// PersistenceController is an interface over Simulation for saving and loading the simulation's state.
type PersistenceController interface {
	Save() error
}

// Save copies the game state into the persister and tells it to persist.
func (s *Simulation) Save(p state.Persister) error {
	for i := range s.characters {
		p.QueueCharacter(characterToState(s.characters[i]))
	}

	for i := range s.worlds {
		p.QueueWorld(worldToState(s.worlds[i]))
	}

	return p.Persist()
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

func worldToState(w *model.World) state.World {
	return state.World{
		ID:    string(w.WorldID),
		Rooms: mapRoomstoState(w.Rooms),
	}
}

func mapRoomstoState(r map[model.RoomID]*model.Room) []state.Room {
	var rooms []state.Room
	for _, v := range r {
		rooms = append(rooms, state.Room{
			ID:    int64(v.ID),
			Items: mapContents(v.Container),
		})
	}
	return rooms
}

func mapItem(i *model.Item) *state.Item {
	if i == nil {
		return nil
	}

	return &state.Item{
		ID:             string(i.ID),
		ItemDefinition: int64(i.Definition.ID),
		Items:          mapContents(i.Container),
	}
}

func mapContents(c model.Container) []*state.Item {
	if c == nil {
		return nil
	}
	var items []*state.Item
	for _, v := range c.Items() {
		items = append(items, mapItem(v))
	}
	return items
}
