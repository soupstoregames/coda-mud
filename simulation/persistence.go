package simulation

import (
	"fmt"
	"time"

	"github.com/soupstore/coda/simulation/data/state"
	"github.com/soupstore/coda/simulation/model"
	"github.com/soupstore/go-core/logging"
)

// StateController is an interface over Simulation for saving and loading the simulation's state.
type StateController interface {
	Save(state.Persister) error
	Load(characters []state.Character, worlds []state.World) error
}

// Save copies the game state into the persister and tells it to persist.
func (s *Simulation) Save(p state.Persister) error {
	start := time.Now()

	for i := range s.characters {
		p.QueueCharacter(characterToState(s.characters[i]))
	}

	for i := range s.worlds {
		p.QueueWorld(worldToState(s.worlds[i]))
	}

	err := p.Persist()

	logging.Info(fmt.Sprintf("Saved game in %v", time.Since(start)))

	return err
}

// Load takes in characters and world states and writes them into the simulation.
// It is naive in that it wont sync the states, simply load the state on top.
// It is only to be used once, before starting the simulation.
func (s *Simulation) Load(characters []state.Character, worlds []state.World) error {
	loaded := 0
	for _, ch := range characters {
		room, err := s.GetRoom(model.WorldID(ch.World), model.RoomID(ch.Room))
		if err != nil {
			return err
		}

		// create new character
		character := model.NewCharacter(ch.Name, room)
		character.ID = model.CharacterID(ch.ID)

		// equip character's rig
		if ch.Rig.Backpack != nil {
			itemDefinition := s.itemDefinitions[model.ItemDefinitionID(ch.Rig.Backpack.ItemDefinition)]
			item := itemDefinition.Spawn()
			item.ID = model.ItemID(ch.Rig.Backpack.ID)
			character.Rig.Backpack = item
		}

		// spawn in character's items
		for _, i := range ch.Items {
			definition, ok := s.itemDefinitions[model.ItemDefinitionID(i.ItemDefinition)]
			if !ok {
				logging.Error("failed to load item for character")
				continue
			}

			item := definition.Spawn()
			item.ID = model.ItemID(i.ID)
			character.Container.PutItem(item)
		}

		s.characters[character.ID] = character

		// add character to room
		room.AddCharacter(character)

		loaded++
	}
	logging.Info(fmt.Sprintf("Loaded %d characters", loaded))

	for _, w := range worlds {
		world, ok := s.worlds[model.WorldID(w.ID)]
		if !ok {
			logging.Warn(fmt.Sprintf("Tried to load state for non-existant world %s", w.ID))
			continue
		}

		for _, r := range w.Rooms {
			room, ok := world.Rooms[model.RoomID(r.ID)]
			if !ok {
				logging.Warn(fmt.Sprintf("Tried to load state for non-existant room %d in world %s", r.ID, w.ID))
				continue
			}

			for _, i := range r.Items {
				definition, ok := s.itemDefinitions[model.ItemDefinitionID(i.ItemDefinition)]
				if !ok {
					logging.Warn(fmt.Sprintf("Tried to load item for non-existant definition %d in room %d in world %s", i.ItemDefinition, r.ID, w.ID))
					continue
				}

				item := definition.Spawn()
				item.ID = model.ItemID(i.ID)
				room.Container.PutItem(item)
			}
		}
	}

	return nil
}

func characterToState(c *model.Character) state.Character {
	return state.Character{
		ID:    string(c.ID),
		Name:  c.Name,
		Room:  int64(c.Room.ID),
		World: string(c.Room.WorldID),
		Rig:   mapRig(c.Rig),
		Items: mapContents(c.Container),
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
